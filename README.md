curl para obtener link de pago de Stripe:

curl https://api.stripe.com/v1/checkout/sessions \
  -u "[STRIPE_SECRET_KEY]:" \
  -d mode=payment \
  -d "success_url=https://a2xbcmw326.execute-api.us-east-1.amazonaws.com/pago-ok?session_id={CHECKOUT_SESSION_ID}" \
  -d "cancel_url=https://a2xbcmw326.execute-api.us-east-1.amazonaws.com/pago-cancelado" \
  -d "line_items[0][price_data][currency]=usd" \
  -d "line_items[0][price_data][product_data][name]=Compra de asientos" \
  -d "line_items[0][price_data][unit_amount]=15000" \
  -d "line_items[0][quantity]=1" \
  -d "metadata[userId]=u_123" \
  -d "metadata[paymentProviderId]=stripe" \
  -d "metadata[seatIds]=A1,A2"
