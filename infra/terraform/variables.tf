variable "auth_service_envs" {
description = "Variables de entorno de auth service"
type = map(string)
default = {}
}

variable "booking_service_envs" {
description = "Variables de entorno de booking service"
type = map(string)
default = {}
}