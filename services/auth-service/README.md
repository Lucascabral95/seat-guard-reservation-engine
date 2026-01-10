<p align="center">
  <a href="http://nestjs.com/" target="blank"><img src="https://nestjs.com/img/logo-small.svg" width="120" alt="Nest Logo" /></a>
</p>

[circleci-image]: https://img.shields.io/circleci/build/github/nestjs/nest/master?token=abc123def456
[circleci-url]: https://circleci.com/gh/nestjs/nest

  <p align="center">A progressive <a href="http://nodejs.org" target="_blank">Node.js</a> framework for building efficient and scalable server-side applications.</p>
    <p align="center">
<a href="https://www.npmjs.com/~nestjscore" target="_blank"><img src="https://img.shields.io/npm/v/@nestjs/core.svg" alt="NPM Version" /></a>
<a href="https://www.npmjs.com/~nestjscore" target="_blank"><img src="https://img.shields.io/npm/l/@nestjs/core.svg" alt="Package License" /></a>
<a href="https://www.npmjs.com/~nestjscore" target="_blank"><img src="https://img.shields.io/npm/dm/@nestjs/common.svg" alt="NPM Downloads" /></a>
<a href="https://circleci.com/gh/nestjs/nest" target="_blank"><img src="https://img.shields.io/circleci/build/github/nestjs/nest/master" alt="CircleCI" /></a>
<a href="https://discord.gg/G7Qnnhy" target="_blank"><img src="https://img.shields.io/badge/discord-online-brightgreen.svg" alt="Discord"/></a>
<a href="https://opencollective.com/nest#backer" target="_blank"><img src="https://opencollective.com/nest/backers/badge.svg" alt="Backers on Open Collective" /></a>
<a href="https://opencollective.com/nest#sponsor" target="_blank"><img src="https://opencollective.com/nest/sponsors/badge.svg" alt="Sponsors on Open Collective" /></a>
  <a href="https://paypal.me/kamilmysliwiec" target="_blank"><img src="https://img.shields.io/badge/Donate-PayPal-ff3f59.svg" alt="Donate us"/></a>
    <a href="https://opencollective.com/nest#sponsor"  target="_blank"><img src="https://img.shields.io/badge/Support%20us-Open%20Collective-41B883.svg" alt="Support us"></a>
  <a href="https://twitter.com/nestframework" target="_blank"><img src="https://img.shields.io/twitter/follow/nestframework.svg?style=social&label=Follow" alt="Follow us on Twitter"></a>
</p>
  <!--[![Backers on Open Collective](https://opencollective.com/nest/backers/badge.svg)](https://opencollective.com/nest#backer)
  [![Sponsors on Open Collective](https://opencollective.com/nest/sponsors/badge.svg)](https://opencollective.com/nest#sponsor)-->

## Auth Service — SeatGuard

Versión profesional y centrada del servicio de autenticación del monorepo SeatGuard. Este servicio gestiona registros, inicio de sesión, emisión y validación de JWT y la integración con la capa de persistencia (Prisma/Postgres).

## Estado

- Stack: NestJS + TypeScript
- ORM: Prisma
- Entrypoint: `src/main.ts`

## Objetivo

Proveer un servicio seguro y escalable para la gestión de identidades que soporte:
- Registro de usuarios
- Inicio de sesión (login) con hashing seguro de contraseñas
- Emisión y verificación de tokens JWT
- Endpoints protegidos por guardas de autorización

## Características principales

- Autenticación basada en JWT con expiración configurable.
- Validaciones estrictas de input usando DTOs.
- Integración directa con Prisma para abstracción de base de datos.
- Modularidad: `auth` y `prisma` son módulos independientes y testeables.
- Buenas prácticas de seguridad (hash de contraseñas, validadores, manejo de errores de Prisma centralizado).

## Estructura relevante

- `src/auth/` — controladores, servicios, guards y DTOs.
- `src/prisma/` — módulo y servicio de Prisma.
- `src/main.ts` — bootstrap de la app.
- `prisma/schema.prisma` — modelo de datos.

## Requisitos locales

- Node.js 18+
- npm o yarn
- Base de datos PostgreSQL accesible (Neon o similar)

## Configuración de entorno

Copiar `.env.template` a `.env` y rellenar las variables necesarias. Variables principales:

- `PORT` — Puerto en el que corre el microservicio.
- `DATABASE_URL` — URL de conexión a Postgres (formato Prisma).
- `MY_FRONTEND_URL` — URL del frontend para configurar CORS.
- `SECRET_JWT` — secreto para firmar JWT.
- `NODE_ENV` — Debe ser 'development' o 'production', segun tu entorno actual.
- `EXPIRED_TOKEN_JWT` — duración por defecto del token (ej. `30d`).

Ejemplo mínimo de `.env`:

```env
PORT=3000
DATABASE_URL="postgresql://user:password@host:5432/auth_db?schema=public"
MY_FRONTEND_URL="http://localhost:4200"
SECRET_JWT="replace_with_a_strong_secret"
NODE_ENV="development"
EXPIRED_TOKEN_JWT="30d"
```

## Scripts útiles

- `npm install` — instala dependencias.
- `npm run start` — ejecuta en modo producción.
- `npm run start:dev` — modo desarrollo (watch).
- `npm run test` — ejecuta tests unitarios.
- `npm run test:e2e` — tests e2e.

Revisa el `package.json` para más detalles.

## Docker

El servicio incluye un `Dockerfile` para construir una imagen lista para ECS/Fargate o pruebas locales.

Ejemplo de build y run local (sin Docker Compose):

```bash
docker build -t seatguard-auth:local .
docker run --env-file .env -p 3000:3000 seatguard-auth:local
```

## API principal (resumen)

Nota: las rutas pueden estar prefijadas por la configuración del router en `src/main.ts`.

- POST `/auth/register` — Registrar usuario (body: `email`, `password`, `name`).
- POST `/auth/login` — Autenticar y recibir `accessToken` (body: `email`, `password`).
- GET `/auth/profile` — Obtener perfil (requiere `Authorization: Bearer <token>`).

Ejemplo cURL de login:

```bash
curl -X POST http://localhost:3000/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"MyP@ssw0rd"}'
```

## Seguridad y buenas prácticas implementadas

- Contraseñas almacenadas usando bcrypt con sal.
- Tokens JWT firmados con `SECRET_JWT` y expiración configurable.
- Guardas (`AuthGuard`) protegiendo rutas privadas.
- Manejo centralizado de errores de Prisma en `src/errors/handler-prisma-error.ts`.
- Validadores personalizados (ej. contraseña fuerte) en `src/common/validators`.

## Desarrollo y pruebas

1. Configurar `.env` con `DATABASE_URL` apuntando a una DB de desarrollo.
2. Ejecutar migraciones/seed de Prisma si aplica:

```bash
npx prisma migrate deploy
npx prisma db seed
```

3. Ejecutar la app en modo desarrollo:

```bash
npm run start:dev
```

4. Ejecutar tests unitarios:

```bash
npm run test
```

## Observabilidad

- Logs estándar via `console` y NestJS logger.
- Instrumentación y métricas pueden agregarse en `src/main.ts` según la plataforma de despliegue.

## Despliegue

- El servicio está preparado para ejecutarse en contenedores y desplegarse en ECS/Fargate.
- Asegúrate de inyectar las variables de entorno sensibles (`DATABASE_URL`, `SECRET_JWT`) en el entorno de producción.

## Contribuir

- Seguir las convenciones de commits (Conventional Commits).
- Añadir pruebas para nueva lógica en `src/auth`.

## Contacto

- Autor del repositorio: Lucas Cabral — lucassimple@hotmail.com

---

Archivo actualizado: [services/auth-service/README.md](services/auth-service/README.md)
