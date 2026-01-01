const Stripe = require("stripe");

const stripeSecretKey = process.env.STRIPE_SECRET_KEY;
const stripe = stripeSecretKey ? new Stripe(stripeSecretKey) : null;

const internalApiKey = process.env.INTERNAL_API_KEY || process.env.SECRET_X_INTERNAL_SECRET || "";
const defaultUrl = "http://monorepo-prod-alb-895613190.us-east-1.elb.amazonaws.com:8080"; 
const bookingServiceBaseUrl = (process.env.BOOKING_SERVICE_URL || defaultUrl).replace(/\/$/, "");

console.log("🌐 Booking Service URL:", bookingServiceBaseUrl);

// ============================================
// HTTP REQUEST (MEJORADO PARA NO FALLAR CON 404)
// ============================================
const requestJson = async (method, url, body) => {
  const controller = new AbortController();
  const timeoutMs = Number(process.env.HTTP_TIMEOUT_MS || 15000);
  const t = setTimeout(() => controller.abort(), timeoutMs);

  const headers = { "content-type": "application/json" };
  if (internalApiKey) headers["X-Internal-Secret"] = internalApiKey;

  console.log(`📡 [${method}] ${url}`);
  if (body) console.log(`📦`, JSON.stringify(body, null, 2));

  try {
    const res = await fetch(url, { method, headers, body: body ? JSON.stringify(body) : undefined, signal: controller.signal });
    const text = await res.text();
    
    // ✅ FIX: Parsing seguro. Si falla (ej: 404 text), guarda el texto crudo.
    let json;
    try {
        json = text ? JSON.parse(text) : {};
    } catch (e) {
        json = { rawResponse: text };
    }

    if (!res.ok) {
      console.error(`❌ HTTP ${res.status}:`, text);
      const err = new Error(`HTTP_ERROR_${res.status}`);
      err.status = res.status;
      err.response = json;
      throw err; 
    }

    return json;
  } catch (error) {
    console.error(`🔥 ERROR: ${error.message}`);
    throw error;
  } finally {
    clearTimeout(t);
  }
};

// ============================================
// OBTENER DATOS DEL CLIENTE DESDE STRIPE
// ============================================
const getCustomerDataFromStripe = async (paymentIntentId) => {
  if (!stripe || !paymentIntentId) {
    console.warn("⚠️ Stripe no disponible");
    return { email: "noreply@booking.com", name: "Cliente" };
  }

  try {
    console.log(`🔍 Consultando PaymentIntent: ${paymentIntentId}`);
    
    const pi = await stripe.paymentIntents.retrieve(paymentIntentId, {
      expand: ['customer', 'latest_charge']
    });
    
    let email = null;
    let name = null;

    if (pi.latest_charge?.billing_details) {
      email = pi.latest_charge.billing_details.email;
      name = pi.latest_charge.billing_details.name;
    }

    if (!email && pi.customer && typeof pi.customer === 'object') {
      email = pi.customer.email;
      name = pi.customer.name;
    }

    if (!email) email = pi.receipt_email;

    if (!name && email) name = email.split('@')[0].replace(/[._-]/g, ' ');

    console.log(`✅ Cliente: ${name} <${email}>`);

    return {
      email: email || "noreply@booking.com",
      name: name || "Cliente",
      customerId: typeof pi.customer === 'string' ? pi.customer : null
    };

  } catch (error) {
    console.error(`❌ Error Stripe: ${error.message}`);
    return { email: "noreply@booking.com", name: "Cliente", customerId: null };
  }
};

// ============================================
// UTILS
// ============================================
const parseSeatIds = (raw) => {
  if (Array.isArray(raw)) return raw.map(String).filter(Boolean);
  if (typeof raw === "string") return raw.split(",").map(s => s.trim()).filter(Boolean);
  return [];
};

const toPaymentStatus = (raw) => {
  const s = String(raw || "").toLowerCase();
  if (["paid", "complete", "completed", "succeeded"].includes(s)) return "COMPLETED";
  if (["failed", "canceled", "cancelled"].includes(s)) return "FAILED";
  return "PENDING";
};

// ============================================
// OPERACIONES
// ============================================
const updateBookingOrder = async (orderId, status, paymentProviderId) => {
  await requestJson("PATCH", `${bookingServiceBaseUrl}/api/v1/booking-orders/${orderId}`, {
    status, paymentProviderId
  });
};

const markSeatsSold = async (seatIds) => {
  if (!seatIds.length) return;
  await Promise.all(seatIds.map(id => 
    requestJson("PATCH", `${bookingServiceBaseUrl}/api/v1/seats/${id}`, { status: "SOLD" })
  ));
};

const updateEventAvailability = async (eventId) => {
  if (!eventId) return;
  try {
    await requestJson("PATCH", `${bookingServiceBaseUrl}/api/v1/events/availability/${eventId}`);
  } catch (e) {
    console.warn("⚠️ Availability update failed (non-critical)");
  }
};

const createCheckout = async (orderId, data) => {
  await requestJson("POST", `${bookingServiceBaseUrl}/api/v1/checkouts`, {
    orderId,
    paymentProvider: "STRIPE",
    paymentIntentId: data.paymentProviderId,
    currency: data.currency || "usd",
    amount: data.amount,
    customerEmail: data.customerEmail,
    customerName: data.customerName,
    customerId: data.customerId
  });
  console.log(`✅ Checkout: ${data.customerName} <${data.customerEmail}>`);
};

// ✅ FIX: RUTA CORREGIDA A /tickets/create
const createTicket = async (orderId) => {
  await requestJson("POST", `${bookingServiceBaseUrl}/api/v1/tickets`, { orderId });
  console.log(`✅ Ticket creado para orden ${orderId}`);
};

// ============================================
// PROCESAMIENTO
// ============================================
const processOne = async (rawBody) => {
  const msg = {
    userId: rawBody.userId,
    amount: Number(rawBody.amount || 0),
    paymentStatus: rawBody.status,
    seatIds: parseSeatIds(rawBody.seatIds),
    eventId: rawBody.eventId,
    paymentProviderId: rawBody.paymentProviderId,
    orderId: rawBody.orderId,
    currency: rawBody.currency || "usd"
  };

  if (!msg.orderId || !msg.seatIds.length) {
    console.error("❌ Mensaje inválido");
    return;
  }

  const status = toPaymentStatus(msg.paymentStatus);
  console.log(`📊 Status: ${status}`);

  if (status !== "COMPLETED") {
    console.log("ℹ️ Pago no completado, ignorando");
    return;
  }

  const customer = await getCustomerDataFromStripe(msg.paymentProviderId);

  await updateBookingOrder(msg.orderId, status, msg.paymentProviderId);
  await markSeatsSold(msg.seatIds);
  await updateEventAvailability(msg.eventId);
  
  await createCheckout(msg.orderId, {
    paymentProviderId: msg.paymentProviderId,
    currency: msg.currency,
    amount: msg.amount,
    customerEmail: customer.email,
    customerName: customer.name,
    customerId: customer.customerId
  });
  
  await createTicket(msg.orderId);

  console.log(`✅ Orden ${msg.orderId} procesada completamente`);
};

// ============================================
// HANDLER
// ============================================
exports.handler = async (event) => {
  console.log("🔔 Lambda iniciada");
  
  if (event.Records) {
    const failures = [];
    for (const record of event.Records) {
      try {
        const body = typeof record.body === "string" ? JSON.parse(record.body) : record.body;
        await processOne(body);
      } catch (err) {
        console.error(`❌ Error:`, err.message);
        failures.push({ itemIdentifier: record.messageId });
      }
    }
    return { batchItemFailures: failures };
  }

  await processOne(event);
  return { statusCode: 200, body: "OK" };
};
