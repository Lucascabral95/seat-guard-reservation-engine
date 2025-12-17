variable "stripe_secret_key" {
description = "Stripe secret key"
type = string
}

variable "url_arn_payment" {
description = "ARN de la queue de pago"
type = string
}

variable "x_internal_secret" {
description = "X-Internal-Secret"
type = string
}