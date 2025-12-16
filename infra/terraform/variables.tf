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

// Para Github Actions
variable "names_images_ecr" {
  description = "Variables de entorno de nombres de las 2 imágenes ECR y mono repo"
  type = map(string)
   default = {
    "auth-service-image"    = "560765037562.dkr.ecr.us-east-1.amazonaws.com/auth-service:latest"
    "booking-service-image" = "560765037562.dkr.ecr.us-east-1.amazonaws.com/booking-service:latest"
    "app_name"              = "monorepo-prod"
  }
}