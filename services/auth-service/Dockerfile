
FROM node:22-alpine AS base
# Instalar OpenSSL (OBLIGATORIO para Prisma en Alpine)
RUN apk add --no-cache openssl libc6-compat

WORKDIR /usr/src/app

# --- ETAPA DEVELOPMENT ---
FROM base AS development

COPY package*.json ./
COPY prisma ./prisma/

RUN npm install

# Generamos el cliente inicial
RUN npx prisma generate

COPY . .

# Comando de inicio
CMD ["sh", "-c", "npx prisma generate && npm run start:dev"]

# Stage 2: Builder
FROM base AS builder

COPY package*.json ./
COPY prisma ./prisma/

# Instalar dependencias (npm ci es más estricto y rápido para builds)
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

# Instalo deps de producción Y prisma CLI para poder regenerar
RUN npm ci --omit=dev && npm install -D prisma

# Regenero el cliente ESPECÍFICAMENTE para el entorno Alpine de producción
RUN npx prisma generate

# Copio el build limpio desde el builder
COPY --from=builder /usr/src/app/dist ./dist

USER node
EXPOSE 3000

# En producción uso el build compilado
CMD ["node", "dist/src/main.js"]
