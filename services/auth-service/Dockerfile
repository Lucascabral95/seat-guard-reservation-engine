# # Stage 1: Base & Dependencies
# FROM node:22-alpine AS base
# WORKDIR /usr/src/app
# COPY package*.json ./
# COPY prisma ./prisma/
# # Instalamos TODO (incluido devDependencies)
# RUN npm ci 
# RUN npx prisma generate

# # Stage 2: Development (NUEVO - Para usar en local)
# FROM base AS development
# # No copiamos todo el código aquí, lo haremos con volúmenes en el docker-compose
# CMD ["npm", "run", "start:dev"]

# # Stage 3: Builder (Para compilar a prod)
# FROM base AS builder
# COPY . .
# RUN npm run build

# # Stage 4: Production (Tu stage actual, optimizado)
# FROM node:22-alpine AS production
# ENV NODE_ENV=production
# WORKDIR /usr/src/app

# COPY package*.json ./
# COPY prisma ./prisma/
# RUN npm ci --omit=dev \
#   && npx prisma generate \
#   && npm cache clean --force

# COPY --from=builder /usr/src/app/dist ./dist

# USER node
# EXPOSE 3000
# CMD ["node", "dist/src/main.js"]

# Stage 1: Base
FROM node:22-alpine AS base
# Instalar OpenSSL (OBLIGATORIO para Prisma en Alpine)
RUN apk add --no-cache openssl libc6-compat

WORKDIR /usr/src/app

# Stage 2: Builder
FROM base AS builder

COPY package*.json ./
COPY prisma ./prisma/

# Instalar dependencias
RUN npm ci

# Generar cliente (ahora se guarda en node_modules/@prisma/client)
RUN npx prisma generate

COPY . .

# Compilar NestJS
RUN npm run build

# Stage 3: Production
FROM base AS production

ENV NODE_ENV=production

COPY package*.json ./
COPY prisma ./prisma/

# Instalamos deps de producción Y prisma CLI para poder regenerar
RUN npm ci --omit=dev && npm install -D prisma

# Regeneramos el cliente ESPECÍFICAMENTE para el entorno Alpine de producción
RUN npx prisma generate

# Copiamos el build limpio
COPY --from=builder /usr/src/app/dist ./dist

USER node
EXPOSE 3000

CMD ["node", "dist/src/main.js"]
