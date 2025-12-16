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


# Stage 1: Base & Dependencies
FROM node:22-alpine AS base
WORKDIR /usr/src/app
# Necesario para que prisma funcione en Alpine si usa binarios nativos de SSL
RUN apk add --no-cache openssl 

COPY package*.json ./
COPY prisma ./prisma/

# Instalamos todo para el build
RUN npm ci 
# Generamos el cliente aquí para asegurar que 'builder' lo tenga
RUN npx prisma generate

# Stage 2: Development
FROM base AS development
CMD ["npm", "run", "start:dev"]

# Stage 3: Builder
FROM base AS builder
COPY . .
# Importante: NestJS usa el cliente generado en node_modules
RUN npm run build

# Stage 4: Production (CORREGIDO)
FROM node:22-alpine AS production
ENV NODE_ENV=production
WORKDIR /usr/src/app

# Importante para Alpine en prod también
RUN apk add --no-cache openssl

COPY package*.json ./
COPY prisma ./prisma/

RUN npm ci --omit=dev && npm i -D prisma

RUN npx prisma generate

RUN npm cache clean --force

COPY --from=builder /usr/src/app/dist ./dist

USER node
EXPOSE 3000
CMD ["node", "dist/src/main.js"]
