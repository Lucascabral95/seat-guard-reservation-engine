<p align="center">
  <img src="https://uxwing.com/wp-content/themes/uxwing/download/festival-culture-religion/tickets-icon.svg"
       alt="Concert seat reservation / tickets"
       width="180"/>
</p>

<h1 align="center">SeatGuard: Reservation Engine</h1>

<p align="center">
  Plataforma de microservicios para venta de entradas de conciertos con bloqueo de asientos en tiempo real y compras en Stripe.
</p>
      
***
## Table of contents

- [Descripción general](#descripción-general)
- [⚙️ Características principales](#️caracteristicas-principales)
- [🏛️ Arquitectura del sistema](#️arquitectura-del-sistema)
  - [Flujo de datos](#flujo-de-datos)
- [Estructura del proyecto](#estructura-del-proyecto)
- [🛠️ Catálogo de microservicios](#️catalogo-de-microservicios)
  - [🔐 Auth Service](#auth-service)
  - [🎟️ Booking Service](#booking-service)
  - [⚡ Payment Lambda](#payment-lambda)
- [🧪 Guía de pruebas de integración con Stripe](#guía-de-pruebas-de-integración-con-stripe)
  - [💳 Tarjetas de prueba](#tarjetas-de-prueba)
  - [🔄 Flujo de Prueba Completo](#flujo-de-prueba-completo)
- [🚀 Guía de instalación y ejecución (local)](#guía-de-instalación-y-ejecución-local)
- [☁️ Guía de despliegue (Terraform)](#️guia-de-despliegue-terraform)
- [🛠️ Scripts avanzados para automatizaciones](#️scripts-avanzados-para-automatizaciones)
- [Contribuciones](#contribuciones)
  - [Convenciones de Commits](#convenciones-de-commits)
- [Licencia](#licencia)
- [📬 Contacto](#contact-anchor)

## Descripción general

**SeatGuard** es una plataforma de microservicios avanzada diseñada para la venta de entradas de conciertos, cuya característica distintiva es un sistema de **bloqueo de asientos en tiempo real**. Implementada como un monorepo robusto, esta solución garantiza una experiencia de compra sin conflictos, permitiendo a los usuarios reservar y bloquear butacas temporalmente por 15 minutos mientras completan su transacción.

El backend está construido con un stack tecnológico moderno que incluye **Go (Gin)** para el alto rendimiento del motor de reservas, **NestJS** para una autenticación segura y escalable, y una arquitectura **Serverless (AWS Lambda + SQS)** para el procesamiento asíncrono de pagos, todo orquestado mediante **Terraform** en **AWS ECS Fargate**.

***

<a id="️caracteristicas-principales"></a>
## ⚙️ Características principales

- Bloqueo en tiempo real: Sistema de concurrencia optimista que bloquea asientos por 15 minutos, previniendo la sobreventa.
- Microservicios desacoplados: Separación clara de responsabilidades entre Autenticación (`auth-service`) y Reservas (`booking-service`).
- Pagos asíncronos: Integración robusta con Stripe mediante Webhooks, colas SQS y funciones Lambda para garantizar la consistencia eventual.
- Infraestructura como código (IaC): Despliegue automatizado y reproducible en AWS utilizando Terraform.
- Seguridad enterprise: Autenticación JWT, validación de datos estricta y protección de endpoints sensibles.
- Base de datos serverless: Integración con Neon Tech (PostgreSQL) para escalabilidad automática y gestión eficiente de conexiones.
- Arquitectura limpia: Diseño de software mantenible siguiendo los principios de Clean Architecture en los servicios principales.

***

<a id="️arquitectura-del-sistema"></a>
## 🏛️ Arquitectura del sistema

El siguiente diagrama ilustra el flujo de datos y la interacción entre los componentes de la plataforma:

```mermaid
graph TD
    subgraph "Cliente"
        Client[Usuario Final]
    end

    subgraph "AWS Cloud"
        ALB[Application Load Balancer]

        subgraph "Amazon ECS (Fargate)"
            AuthService[Auth Service (NestJS)]
            BookingService[Booking Service (Go)]
        end

        subgraph "Cloud Database"
            NeonDB[Neon Tech PostgreSQL]
        end

        subgraph "Procesamiento Asíncrono"
            SQS[SQS Queue]
            Lambda[Payment Processor Lambda]
        end

        subgraph "Pasarela de Pagos Externa"
            Stripe[Stripe API]
            StripeWebhook[Stripe Webhook]
        end
    end

    Client -- HTTPS --> ALB
    ALB -- :80/auth/* --> AuthService
    ALB -- :8080/api/v1/* --> BookingService

    AuthService <--> NeonDB
    BookingService <--> NeonDB

    BookingService -- Crea Sesión de Pago --> Stripe
    StripeWebhook -- Evento de Pago --> Booking Service
    Lambda -- Lee Mensaje --> SQS
    Lambda -- Actualiza Estado --> BookingService
    BookingService -- Recalcula Disponibilidad --> NeonDB
``` 

## Flujo de datos
- **Autenticación**: El cliente obtiene un token JWT a través del Auth Service.

- **Reserva**: El usuario selecciona y bloquea asientos en tiempo real mediante el Booking Service.

- **Pago**: Se inicia una sesión de checkout en Stripe.

- **Confirmación asíncrona**:
  - **Notificación de pago**: Stripe envía un webhook al Booking Service cuando se completa un pago.
  - **Cola de mensajes**: El mensaje se encola automáticamente en SQS para garantizar la entrega.
  - **Procesamiento**: La Lambda consume el mensaje de la cola y:
    - Valida el pago con Stripe
    - Marca los asientos como `SOLD` en la base de datos
    - Crea la orden de compra en el Booking Service
    - Actualiza la disponibilidad del evento
  - **Notificación al usuario**: Se envía confirmación al usuario (opcional, vía email/websocket)
  - **Manejo de errores**: Si falla algún paso, el mensaje vuelve a la cola para reintento.

## Estructura del proyecto
```
seatguard-monorepo/
├── infra/                           # Infraestructura como Código
│   └── terraform/                   # Scripts de Terraform (ECS, ALB, VPC)
├── lambdas/                         # Funciones Serverless
│   └── payment-processor/           # Lógica de procesamiento de pagos (Node.js)
├── services/                        # Microservicios
│   ├── auth-service/                # Servicio de Autenticación (NestJS)
│   │   ├── src/
│   │   ├── prisma/                  # Esquema de base de datos
│   │   └── Dockerfile
│   └── booking-service/             # Motor de Reservas (Go)
│       ├── cmd/                     # Punto de entrada
│       ├── internal/                # Lógica de negocio (Dominios, Servicios, Repos)
│       └── Dockerfile
├── scripts/                         # Scripts de utilidad (Deploy, Build)
├── docker-compose.yml               # Orquestación local
└── README2.md                       # Documentación
```

<a id="️catalogo-de-microservicios"></a>
## 🛠️ Catálogo de microservicios

<a id="auth-service"></a>
### 🔐 Auth Service

- **Ruta**: `services/auth-service`
- **Stack**: NestJS, TypeScript, Prisma ORM, JWT
- **Función**: Gestión de identidades, registro y emisión de tokens JWT.

<a id="booking-service"></a>
### 🎟️ Booking Service

- **Ruta**: `services/booking-service`
- **Stack**: Go (Golang), Gin Framework, GORM
- **Función**:
  - Gestión de inventario (eventos, asientos, órdenes).
  - Bloqueo temporal de asientos para evitar sobreventa.
  - Creación de sesiones de checkout con Stripe.

<a id="payment-lambda"></a>
### ⚡ Payment Lambda

- **Ruta**: `lambdas/payment-processor`
- **Stack**: Node.js
- **Función**: Worker asíncrono que procesa webhooks/mensajes, marca asientos como `SOLD`, crea la orden y actualiza disponibilidad.

<a id="guía-de-pruebas-de-integración-con-stripe"></a>
## 🧪 Guía de pruebas de integración con Stripe

Esta sección detalla cómo validar el flujo de pagos completo en el entorno de desarrollo utilizando las tarjetas de prueba de Stripe.

---

<a id="tarjetas-de-prueba"></a>
## 💳 Tarjetas de prueba

### ✅ Pagos exitosos

| Campo | Valor |
| :--- | :--- |
| **Número de tarjeta** | `4242 4242 4242 4242` |
| **Red** | Visa |
| **Fecha (MM/YY)** | Cualquier fecha futura (ej: `12/30`) |
| **CVC** | Cualquier 3 dígitos (ej: `123`) |
| **Nombre del titular** | Cualquier nombre |
| **Código postal** | Cualquier código válido |

### ❌ Simulación de errores

| Escenario | Número de tarjeta | Código de error | Código de rechazo |
| :--- | :--- | :--- | :--- |
| Rechazo genérico | `4000 0000 0000 0002` | `card_declined` | `generic_decline` |
| Fondos insuficientes | `4000 0000 0000 9995` | `card_declined` | `insufficient_funds` |
| Tarjeta robada | `4000 0000 0000 9979` | `card_declined` | `stolen_card` |
| Tarjeta caducada | `4000 0000 0000 0069` | `expired_card` | — |
| CVC incorrecto | `4000 0000 0000 0127` | `card_declined` | `incorrect_cvc` |
| Error de procesamiento | `4000 0000 0000 0119` | `processing_error` | — |
| Número incorrecto | `4242 4242 4242 4241` | `incorrect_number` | — |
| Límite de velocidad | `4000 0000 0000 6975` | `card_declined` | `card_velocity_exceeded` |

**Referencia**: documentación oficial de Stripe:
https://stripe.com/docs/testing

---

<a id="flujo-de-prueba-completo"></a>
## 🔄 Flujo de Prueba Completo

### Pasos para Realizar una Prueba

1. **Crear checkout**: enviar un `POST /api/v1/stripe/create/checkout/session` con los `seatIds` seleccionados.
```
curl -X POST "http://localhost:8080/api/v1/stripe/create/checkout/session" \
  -H "Content-Type: application/json" \
  -d '{
    "userId": "user_3011",
    "currency": "usd",
    "items": [
      {
        "name": "Dua lipa",
        "amount": 35,
        "seatIds": {
          "id": "6e7af213-8941-47e6-8850-39e99c1a4073"
        }
      }
    ]
  }'
```
2. **Completar el pago**: abrir la `url` devuelta por Stripe y pagar con la tarjeta de prueba `4242 4242 4242 4242`.
3. **Verificar procesamiento**:
   - Revisar logs de la Lambda en CloudWatch (o local si lo emulas).
   - Confirmar en el Booking Service:
     - Los asientos quedaron en estado `SOLD`.
     - Se creó la orden con estado `COMPLETED`.

---

<a id="guía-de-instalación-y-ejecución-local"></a>
## 🚀 Guía de instalación y ejecución (local)

### Prerrequisitos

- Docker y Docker Compose
- Go 1.21+
- Node.js 18+

### 1) Clonar el repositorio

```bash
git clone https://github.com/Lucascabral95/seatguard-reservation-engine.git
cd seatguard-reservation-engine
```

### 2) Configuración de entorno

Crear los archivos `.env` en cada servicio basándote en los templates (`.env.template`).

#### Variables principales

| Variable | Servicio(s) | Descripción |
| :--- | :--- | :--- |
| `DATABASE_URL` | Auth, Booking | URL de conexión a PostgreSQL (Neon u otra). |
| `JWT_SECRET` | Booking | Secreto para validar JWT (debe coincidir con el del Auth Service). |
| `SECRET_JWT` | Auth | Secreto para firmar JWT. |
| `STRIPE_SECRET_KEY` | Booking, Lambda | API key secreta de Stripe. |
| `SQS_QUEUE_URL` | Booking, Lambda | URL de la cola SQS para mensajería. |

### 3) Ejecución local (Docker)

```bash
docker-compose up --build
```

- `auth-service`: `http://localhost:3000`
- `booking-service`: `http://localhost:4000`

### 4) Seeding de datos (opcional)

```bash
cd services/booking-service && go run cmd/api/main.go -seed
```

---

<a id="️guia-de-despliegue-terraform"></a>
## ☁️ Guía de despliegue (Terraform)

#### El despliegue en AWS se gestiona desde `infra/terraform`.

Para desplegar VPC, ALB, ECS (FARGATE) con dos Task Definitons (Auth :80 (pueto predeterminado) y Booking :8080):

```bash
cd infra/terraform 
terraform init 
terrform plan
terraform apply -auto-approve
```

Para desplesgar la lambda que escuche a SQS: 

```bash
cd infra/terraform/lambda-webhook
terraform init 
terrform plan
terraform apply -auto-approve
```

<a id="️scripts-avanzados-para-automatizaciones"></a>
## 🛠️ Scripts avanzados para automatizacions:

#### Activar la semilla de la base de datos de Booking-service (Go): 
```
npm run go:seed
```

#### Levantar microservicio de Booking-service (Go): 
```
npm run go:run
```

#### Levantar microservicio de Auth-service (NestJS): 
```
npm run nest:run
```

#### Hacer deploy de los dos microservicios a ECS (Fargate): 
1) Crear variables de entorno para levantar la ECS (Fargate) de los dos microservicios:
```
# Crear archivo terraform.tfvars en /infra/terraform y escribir:
auth_service_envs = {
  # Auth service
  PORT      = "3000"
  DATABASE_URL       = "postgresql://[your_db_user]:[your_db_password]@[your_db_host]:[your_db_port]/auth_db?sslmode=require&channel_binding=require"
  MY_FRONTEND_URL      = "http://localhost:3000"
  SECRET_JWT  = "yoursecret"
  EXPIRED_TOKEN_JWT     = "30d"
}

booking_service_envs = {
# Booking service
  PORT = "4000"
  JWT_SECRET = "yoursecret"
  DB_HOST = "your_db_host"
  DB_USER = "your_db_user"
  DB_PASSWORD = "your_db_password"
  DB_NAME = "your_db_name"
  DB_PORT = "your_db_port"

  DB_URL = "postgresql://[your_db_user]:[your_db_password]@[your_db_host]:[your_db_port]/booking_db?sslmode=require&channel_binding=require"
  
  ENV_MODE = "development"

  AWS_REGION = "your_aws_region"
  SQS_QUEUE_URL = "https://sqs.your_aws_region.amazonaws.com/your_aws_account_id/payment-queue-reservation"
  AWS_ACCESS_KEY_ID = "your_aws_access_key_id"
  AWS_SECRET_ACCESS_KEY = "your_aws_secret_access_key"

  STRIPE_SECRET_KEY = "your_stripe_secret_key"
  STRIPE_SUCCESS_URL = "http://localhost:4000/api/v1/events"
  STRIPE_CANCEL_URL = "http://localhost:4000/api/v1/events"
  STRIPE_WEBHOOK_SECRET = "http://localhost:4000/api/v1/stripe/webhook"
}
```

2) Ejecutar comandos de init (solo primera vez) y de ejecución: 
``` 
# Esto se ejecuta sólo la primera vez
cd infra/terraform init 

npm run deploy:fargate
```

#### Hacer deploy de la funcion Lambda encargada de Procesar los pagos: 
1) Crear variables de entorno para el levantar la funcion Lambda mediante Terraform: 
```
# Crear archivo terraform.tfvars en /infra/terraform/lambda-webhook y escribir:
stripe_secret_key = "sk_test_[your_stripe_secret_key]"
```

2) Ejecutar comandos de init (solo primera vez) y de ejecución:
```
# Esto se ejecuta sólo la primera vez
cd lambda/payment-processor && npm install

# Esto se ejecuta sólo la primera vez
cd infra/terraform/lambda-webhook && terraform init

npm run deploy:lambda
```


#### Hacer deploy de tanto de los dos microservicios a ECS (Fargate) y de la funcion Lambda:
```
npm run deploy:all
```

#### Hacer build de las dos imagenes a ECR, subirlas a AWS y ademas actualizar las dos Task Definition de ECS (Fargate):
- (Esto tambien actualiza los dos Task Definition que esta escuchando a las dos imagenes ECR de los microservicios):
```
bash scripts-template/deploy-images-template.sh
```

### ✅ Hacer build local/remota de microservicios:
```
npm run docker:compose
```

#### Hacer build tanto local como a AWS de los dos micrservicios:
```
npm run docker:compose && bash scripts-template/deploy-images-template.sh
```

### ❌ Destruir recursos del proyecto: 

#### Destruir recursos de terraform de ECS (Fargate):
```
npm run destroy:fargate
```

#### Destruir recursos de terraform de funcion Lambda:
```
npm run destroy:lambda
```

#### Destruir recursos de terraform de todos los recursos:
```
npm run destroy:all
```

***

## Contribuciones

¡Las contribuciones son bienvenidas! Seguí estos pasos:

1. Hacé un fork del repositorio.
2. Creá una rama para tu feature o fix (`git checkout -b feature/nueva-funcionalidad`).
3. Realizá tus cambios y escribí pruebas si es necesario.
4. Hacé commit y push a tu rama (`git commit -m "feat: agrega nueva funcionalidad"`).
5. Abrí un Pull Request describiendo tus cambios.

### Convenciones de Commits

Este proyecto sigue [Conventional Commits](https://www.conventionalcommits.org/):

- `feat:` Nueva funcionalidad
- `fix:` Corrección de bugs
- `docs:` Cambios en documentación
- `style:` Cambios de formato (no afectan la lógica)
- `refactor:` Refactorización de código
- `test:` Añadir o modificar tests
- `chore:` Tareas de mantenimiento

---

## Licencia

Este proyecto está bajo la licencia **MIT**.

---

<a id="contact-anchor"></a>
## 📬 Contacto

- **Autor:** Lucas Cabral
- **Email:** lucassimple@hotmail.com
- **LinkedIn:** [https://www.linkedin.com/in/lucas-gastón-cabral/](https://www.linkedin.com/in/lucas-gastón-cabral/)
- **Portfolio:** [https://portfolio-web-dev-git-main-lucascabral95s-projects.vercel.app/](https://portfolio-web-dev-git-main-lucascabral95s-projects.vercel.app/)
- **Github:** [https://github.com/Lucascabral95](https://github.com/Lucascabral95/)

---

<p align="center">
  Desarrollado con ❤️ por Lucas Cabral
</p>