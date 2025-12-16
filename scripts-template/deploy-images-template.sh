# #!/bin/bash
AWS_REGION="[AWS_REGION]"
ACCOUNT_ID="[ACCOUNT_ID]"

# Nombres de Repos ECR
REPO_AUTH="$ACCOUNT_ID.dkr.ecr.$AWS_REGION.amazonaws.com/auth-service"
REPO_BOOKING="$ACCOUNT_ID.dkr.ecr.$AWS_REGION.amazonaws.com/booking-service"

# Nombres de Cluster y Servicios (Extraídos de tu Terraform)
# Terraform: app_name = "monorepo-prod" -> cluster = "monorepo-prod-cluster"
CLUSTER_NAME="[CLUSTER_NAME]"  

# Terraform: resource "aws_ecs_service" "auth" { name = "auth-service" }
SERVICE_AUTH="[SERVICE_AUTH]"

# Terraform: resource "aws_ecs_service" "booking" { name = "booking-service" }
SERVICE_BOOKING="[SERVICE_BOOKING]"


# --- 1. LOGIN EN ECR ---
echo "🔐 [1/4] Logueándose en ECR..."
aws ecr get-login-password --region $AWS_REGION | docker login --username AWS --password-stdin $ACCOUNT_ID.dkr.ecr.$AWS_REGION.amazonaws.com


# --- 2. BUILD & PUSH AUTH ---
echo "🚀 [2/4] Procesando AUTH SERVICE..."
docker build -t auth-service ./services/auth-service
docker tag auth-service:latest $REPO_AUTH:latest
docker push $REPO_AUTH:latest
echo "   -> Imagen Auth subida."


# --- 3. BUILD & PUSH BOOKING ---
echo "🚀 [3/4] Procesando BOOKING SERVICE..."
docker build -t booking-service ./services/booking-service
docker tag booking-service:latest $REPO_BOOKING:latest
docker push $REPO_BOOKING:latest
echo "   -> Imagen Booking subida."


# --- 4. ACTUALIZAR FARGATE (Force Deployment) ---
echo "🔄 [4/4] Forzando actualización en ECS Fargate..."

# Update Auth
echo "   -> Reiniciando Auth Service ($SERVICE_AUTH)..."
aws ecs update-service \
    --cluster $CLUSTER_NAME \
    --service $SERVICE_AUTH \
    --force-new-deployment \
    --region $AWS_REGION > /dev/null

# Update Booking
echo "   -> Reiniciando Booking Service ($SERVICE_BOOKING)..."
aws ecs update-service \
    --cluster $CLUSTER_NAME \
    --service $SERVICE_BOOKING \
    --force-new-deployment \
    --region $AWS_REGION > /dev/null

echo "✅ DESPLIEGUE EXITOSO! Fargate está bajando las nuevas imágenes."
