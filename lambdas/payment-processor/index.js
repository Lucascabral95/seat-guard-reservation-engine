const Stripe = require("stripe");

const stripeSecretKey = process.env.STRIPE_SECRET_KEY;
const stripe = stripeSecretKey ? new Stripe(stripeSecretKey) : null;

const internalApiKey = process.env.INTERNAL_API_KEY || process.env.SECRET_X_INTERNAL_SECRET || "";

console.log("🔑 Internal API Key configurada:", internalApiKey ? `SÍ (Longitud: ${internalApiKey.length})` : "NO");

const defaultUrl = "http://monorepo-prod-alb-1569205323.us-east-1.elb.amazonaws.com:8080"; 
const bookingServiceBaseUrl = (process.env.BOOKING_SERVICE_URL || defaultUrl).replace(/\/$/, "");

console.log("🌐 URL base del servicio:", bookingServiceBaseUrl);

const requestJson = async (method, url, body) => {
  const controller = new AbortController();
  const timeoutMs = Number(process.env.HTTP_TIMEOUT_MS || 10000);
  const t = setTimeout(() => controller.abort(), timeoutMs);

  const headers = { "content-type": "application/json" };
  if (internalApiKey) {
    headers["X-Internal-Secret"] = internalApiKey;
  }

  console.log(`📡 Intentando ${method} a: ${url}`);

  try {
    const res = await fetch(url, {
      method,
      headers,
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
      console.error(`❌ SERVER ERROR ${res.status} ${url}. Body:`, text);
      const err = new Error(`HTTP ${res.status} ${method} ${url}`);
      err.status = res.status;
      err.response = json;
      throw err;
    }

    return json;
  } catch (error) {
    console.error(`🔥 NETWORK ERROR en ${method} ${url}:`);
    console.error(`   Mensaje: ${error.message}`);
    if (error.cause) console.error(`   Causa:`, error.cause);
    throw error;
  } finally {
    clearTimeout(t);
  }
};

const parseSeatIds = (seatIdsRaw) => {
  if (Array.isArray(seatIdsRaw)) {
    return seatIdsRaw.map(String).map((s) => s.trim()).filter(Boolean);
  }
  if (typeof seatIdsRaw === "string") {
    return seatIdsRaw.split(",").map((s) => String(s).trim()).filter(Boolean);
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
        eventId: body.eventId,
        paymentProviderId: body.paymentProviderId || "",
      };
    }
    if (body.type && body.data && body.data.object) {
      const o = body.data.object;
      const meta = (o && o.metadata) || {};
      return {
        userId: String(meta.user_id || ""),
        amount: Number(o.amount_total || 0),
        paymentStatus: o.payment_status,
        seatIds: parseSeatIds(meta.seat_ids),
        eventId: meta.event_id || "",
        paymentProviderId: o.payment_intent || o.id || "",
      };
    }
  }
  return null;
};

const ensureBookingService = () => {
  if (!bookingServiceBaseUrl) {
    throw new Error("BOOKING_SERVICE_BASE_URL no configurado y fallback vacío");
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
    console.error("❌ Error verificando idempotencia (no crítico):", e.message);
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

const updateEventAvailability = async (eventID) => {
  if (!eventID) return;
  ensureBookingService();
  return requestJson("PATCH", `${bookingServiceBaseUrl}/api/v1/events/availability/${encodeURIComponent(eventID)}`);
};

const processOne = async (rawBody) => {
  const msg = normalizeMessage(rawBody);
  if (!msg) throw new Error("Mensaje no soportado");
  if (!msg.userId) throw new Error("Falta userId");
  if (!msg.seatIds || msg.seatIds.length === 0) throw new Error("Faltan seatIds");
  
  const shouldCheckIdempotency = String(process.env.ENABLE_IDEMPOTENCY_LOOKUP || "true").toLowerCase() !== "false";
  const already = shouldCheckIdempotency ? await bookingOrderExists(msg.paymentProviderId) : false;

  await markSeatsSold(msg.seatIds);
  
  if (!already) {
    await createBookingOrder(msg);
  }

  if (msg.eventId) {
    try {
      await updateEventAvailability(msg.eventId);
      console.log(`Availability updated for event: ${msg.eventId}`);
    } catch (e) {
      console.error("Error updating availability (non-critical):", e.message);
    }
  }

  return { alreadyProcessed: already };
};

exports.handler = async (event) => {
  console.log("🔔 Webhook/SQS Handler iniciado");

  if (event.Records) {
    const batchItemFailures = [];
    for (const record of event.Records) {
      try {
        const body = typeof record.body === "string" ? JSON.parse(record.body) : record.body;
        console.log("📩 Mensaje SQS recibido:", JSON.stringify(body, null, 2));
        await processOne(body);
      } catch (err) {
        console.error("❌ Error procesando mensaje SQS:", err.message);
        if (record && record.messageId) {
          batchItemFailures.push({ itemIdentifier: record.messageId });
        }
      }
    }
    return { batchItemFailures };
  }

  console.log("Evento directo (no SQS):", JSON.stringify(event));
  try {
    await processOne(event);
  } catch (err) {
    console.error("❌ Error procesando evento directo:", err.message);
    return { statusCode: 500, body: "ERROR" };
  }
  return { statusCode: 200, body: "OK" };
};
