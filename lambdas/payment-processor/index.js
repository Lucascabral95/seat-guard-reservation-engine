const Stripe = require("stripe");

const stripeSecretKey = process.env.STRIPE_SECRET_KEY;
const stripe = stripeSecretKey ? new Stripe(stripeSecretKey) : null;

const bookingServiceBaseUrl = ("http://monorepo-prod-alb-17945421.us-east-1.elb.amazonaws.com:8080" || "http://monorepo-prod-alb-17945421.us-east-1.elb.amazonaws.com:8080" || "").replace(/\/$/, "");

const requestJson = async (method, url, body) => {
  const controller = new AbortController();
  const timeoutMs = Number(process.env.HTTP_TIMEOUT_MS || 8000);
  const t = setTimeout(() => controller.abort(), timeoutMs);

  try {
    const res = await fetch(url, {
      method,
      headers: {
        "content-type": "application/json",
      },
      body: body === undefined ? undefined : JSON.stringify(body),
      signal: controller.signal,
    });

    const text = await res.text();
    let json;
    try {
      json = text ? JSON.parse(text) : undefined;
    } catch {
      json = text;
    }

    if (!res.ok) {
      const err = new Error(`HTTP ${res.status} ${method} ${url}`);
      err.status = res.status;
      err.response = json;
      throw err;
    }

    return json;
  } finally {
    clearTimeout(t);
  }
};

const parseSeatIds = (seatIdsRaw) => {
  if (Array.isArray(seatIdsRaw)) {
    return seatIdsRaw.map(String).map((s) => s.trim()).filter(Boolean);
  }
  if (typeof seatIdsRaw === "string") {
    return seatIdsRaw
      .split(",")
      .map((s) => String(s).trim())
      .filter(Boolean);
  }
  return [];
};

const toPaymentStatus = (paymentStatusRaw) => {
  const s = String(paymentStatusRaw || "").toLowerCase();
  if (s === "paid" || s === "complete" || s === "completed" || s === "succeeded") return "COMPLETED";
  if (s === "failed" || s === "canceled" || s === "cancelled" || s === "unpaid") return "FAILED";
  return "PENDING";
};

const normalizeMessage = (body) => {
  if (body && typeof body === "object") {
    if (body.userId && body.seatIds) {
      return {
        userId: String(body.userId),
        amount: Number(body.amount || 0),
        paymentStatus: body.status,
        seatIds: parseSeatIds(body.seatIds),
        paymentProviderId: body.paymentProviderId || "",
      };
    }

    if (body.type && body.data && body.data.object) {
      const o = body.data.object;
      const meta = (o && o.metadata) || {};
      return {
        userId: String(meta.user_id || ""),
        amount: Number(o.amount_total || 0) / 100,
        paymentStatus: o.payment_status,
        seatIds: parseSeatIds(meta.seat_ids),
        paymentProviderId: o.payment_intent || o.id || "",
      };
    }
  }
  return null;
};

const ensureBookingService = () => {
  if (!bookingServiceBaseUrl) {
    throw new Error("BOOKING_SERVICE_BASE_URL/BOOKING_SERVICE_URL no configurado");
  }
};

const markSeatsSold = async (seatIds) => {
  ensureBookingService();
  for (const seatId of seatIds) {
    await requestJson("PATCH", `${bookingServiceBaseUrl}/api/v1/seats/${encodeURIComponent(seatId)}`, { status: "SOLD" });
  }
};

const bookingOrderExists = async (paymentProviderId) => {
  if (!paymentProviderId) return false;
  ensureBookingService();
  try {
    const orders = await requestJson("GET", `${bookingServiceBaseUrl}/api/v1/booking-orders`);
    if (!Array.isArray(orders)) return false;
    return orders.some((o) => o && String(o.paymentProviderId || "") === String(paymentProviderId));
  } catch (e) {
    console.error("❌ Error verificando idempotencia en booking-orders:", e && e.message ? e.message : e);
    return false;
  }
};

const createBookingOrder = async ({ userId, amount, paymentStatus, seatIds, paymentProviderId }) => {
  ensureBookingService();
  return requestJson("POST", `${bookingServiceBaseUrl}/api/v1/booking-orders`, {
    userId,
    amount,
    status: toPaymentStatus(paymentStatus),
    seatIds,
    paymentProviderId,
  });
};

const processOne = async (rawBody) => {
  const msg = normalizeMessage(rawBody);
  if (!msg) throw new Error("Mensaje no soportado (no es Stripe event ni BookingMessage)");
  if (!msg.userId) throw new Error("Falta userId");
  if (!msg.seatIds || msg.seatIds.length === 0) throw new Error("Faltan seatIds");
  if (!Number.isFinite(msg.amount)) throw new Error("Amount inválido");

  const shouldCheckIdempotency = String(process.env.ENABLE_IDEMPOTENCY_LOOKUP || "true").toLowerCase() !== "false";
  const already = shouldCheckIdempotency ? await bookingOrderExists(msg.paymentProviderId) : false;

  await markSeatsSold(msg.seatIds);
  if (!already) {
    await createBookingOrder(msg);
  }

  return { alreadyProcessed: already };
};

exports.handler = async (event) => {
  console.log("🔔 Webhook/SQS Handler iniciado");

  // Si el evento viene de SQS, itero los registros
  if (event.Records) {
    const batchItemFailures = [];

    for (const record of event.Records) {
      try {
        const body = typeof record.body === "string" ? JSON.parse(record.body) : record.body;
        console.log("📩 Mensaje SQS recibido:", JSON.stringify(body, null, 2));
        await processOne(body);
      } catch (err) {
        console.error("❌ Error procesando mensaje SQS:", err && err.message ? err.message : err);
        if (record && record.messageId) {
          batchItemFailures.push({ itemIdentifier: record.messageId });
        }
      }
    }

    if (batchItemFailures.length > 0) {
      return { batchItemFailures };
    }
    return { statusCode: 200, body: "SQS Processed" };
  }

  console.log("Evento directo (no SQS):", JSON.stringify(event));
  try {
    await processOne(event);
  } catch (err) {
    console.error("❌ Error procesando evento directo:", err && err.message ? err.message : err);
    return { statusCode: 500, body: "ERROR" };
  }
  return { statusCode: 200, body: "OK" };
};