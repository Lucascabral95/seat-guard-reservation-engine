const Stripe = require("stripe");

const stripeSecretKey = process.env.STRIPE_SECRET_KEY;
const stripe = stripeSecretKey ? new Stripe(stripeSecretKey) : null;

const internalApiKey = process.env.INTERNAL_API_KEY || process.env.SECRET_X_INTERNAL_SECRET || "";

const defaultUrl = "http://monorepo-prod-alb-895613190.us-east-1.elb.amazonaws.com:8080"; 
const bookingServiceBaseUrl = (process.env.BOOKING_SERVICE_URL || defaultUrl).replace(/\/$/, "");

console.log("🌐 URL base del servicio:", bookingServiceBaseUrl);

const requestJson = async (method, url, body) => {
  const controller = new AbortController();
  const timeoutMs = Number(process.env.HTTP_TIMEOUT_MS || 15000);
  const t = setTimeout(() => controller.abort(), timeoutMs);

  const headers = { "content-type": "application/json" };
  if (internalApiKey) {
    headers["X-Internal-Secret"] = internalApiKey;
  }

  console.log(`📡 [${method}] ${url}`);
  if (body) console.log(`📦 Payload:`, JSON.stringify(body));

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
      json = text ? JSON.parse(text) : {};
    } catch {
      json = { raw: text };
    }

    if (!res.ok) {
      console.error(`❌ SERVER ERROR ${res.status} en ${url}. Respuesta:`, text);
      const err = new Error(`HTTP_ERROR_${res.status}`);
      err.status = res.status;
      err.response = json;
      throw err; 
    }

    return json;
  } catch (error) {
    console.error(`🔥 NETWORK/FETCH ERROR en ${method} ${url}: ${error.message}`);
    throw error;
  } finally {
    clearTimeout(t);
  }
};

const parseSeatIds = (seatIdsRaw) => {
  if (Array.isArray(seatIdsRaw)) return seatIdsRaw.map(String).map(s => s.trim()).filter(Boolean);
  if (typeof seatIdsRaw === "string") return seatIdsRaw.split(",").map(s => String(s).trim()).filter(Boolean);
  return [];
};

const toPaymentStatus = (paymentStatusRaw) => {
  const s = String(paymentStatusRaw || "").toLowerCase();
  if (["paid", "complete", "completed", "succeeded"].includes(s)) return "COMPLETED";
  if (["failed", "canceled", "cancelled", "unpaid"].includes(s)) return "FAILED";
  return "PENDING";
};

const normalizeMessage = (body) => {
  if (body && body.userId) {
    return {
      userId: String(body.userId),
      amount: Number(body.amount || 0),
      paymentStatus: body.status, 
      seatIds: parseSeatIds(body.seatIds),
      eventId: body.eventId,
      paymentProviderId: body.paymentProviderId || "",
      orderId: body.orderId || "", 
    };
  }
  if (body && body.type && body.data && body.data.object) {
    const o = body.data.object;
    const meta = (o && o.metadata) || {};
    return {
      userId: String(meta.user_id || ""),
      amount: Number(o.amount_total || 0),
      paymentStatus: o.payment_status,
      seatIds: parseSeatIds(meta.seat_ids),
      eventId: meta.event_id || "",
      paymentProviderId: o.payment_intent || o.id || "",
      orderId: meta.order_id || "", 
    };
  }
  return null;
};

const markSeatsSold = async (seatIds) => {
  if (!seatIds.length) return;
  const promises = seatIds.map(seatId => 
    requestJson("PATCH", `${bookingServiceBaseUrl}/api/v1/seats/${encodeURIComponent(seatId)}`, { status: "SOLD" })
  );
  await Promise.all(promises);
};

const updateEventAvailability = async (eventID) => {
  if (!eventID) return;
  try {
      await requestJson("PATCH", `${bookingServiceBaseUrl}/api/v1/events/availability/${encodeURIComponent(eventID)}`);
      console.log(`Availability updated for event: ${eventID}`);
  } catch (e) {
    console.warn("⚠️ Warning: Error updating availability (non-critical):", e.message);
  }
};

const updateBookingOrder = async (orderId, status, paymentProviderId) => {
  console.log(`🔄 Actualizando Orden ${orderId} -> Status: ${status}, ProviderID: ${paymentProviderId}`);
  
  const payload = {
    status: status,
    paymentProviderId: paymentProviderId 
  };

  return requestJson("PATCH", `${bookingServiceBaseUrl}/api/v1/booking-orders/${encodeURIComponent(orderId)}`, payload);
};

const processOne = async (rawBody) => {
  const msg = normalizeMessage(rawBody);
  if (!msg) {
    console.error("Mensaje ignorado (formato desconocido):", JSON.stringify(rawBody));
    return; 
  }

  console.log("📨 Procesando mensaje normalizado:", JSON.stringify(msg));

  if (!msg.seatIds || msg.seatIds.length === 0) {
    console.error("❌ Mensaje sin seatIds, no se puede procesar.");
    return;
  }
  
  const currentStatus = toPaymentStatus(msg.paymentStatus);

  if (currentStatus === "COMPLETED") {
      console.log(`🎟️ Marcando ${msg.seatIds.length} asientos como SOLD...`);
      await markSeatsSold(msg.seatIds);
      await updateEventAvailability(msg.eventId);
  } else {
      console.log(`ℹ️ Estado de pago: ${currentStatus}. No se modifican asientos.`);
  }

  if (msg.orderId) {
      await updateBookingOrder(msg.orderId, currentStatus, msg.paymentProviderId);
      console.log(`✅ Orden ${msg.orderId} actualizada correctamente.`);
  } else {
      console.error("⚠️ NO SE RECIBIÓ ORDER ID. Esto es inesperado si la orden se creó antes del pago.");
  }

  return { success: true };
};

exports.handler = async (event) => {
  console.log("🔔 Webhook/SQS Handler iniciado");
  
  if (event.Records) {
    const batchItemFailures = [];
    
    for (const record of event.Records) {
      try {
        const body = typeof record.body === "string" ? JSON.parse(record.body) : record.body;
        await processOne(body);
      } catch (err) {
        console.error(`❌ Error procesando mensaje ${record.messageId}:`, err.message);
        batchItemFailures.push({ itemIdentifier: record.messageId });
      }
    }
    
    return { batchItemFailures };
  }

  try {
    await processOne(event);
    return { statusCode: 200, body: "OK" };
  } catch (err) {
    console.error("❌ Error fatal en invocación directa:", err);
    return { statusCode: 500, body: err.message };
  }
};
