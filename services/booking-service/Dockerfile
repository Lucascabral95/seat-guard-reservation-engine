# syntax=docker/dockerfile:1

# 1. Etapa de compilación
FROM golang:1.25 AS build
WORKDIR /app

# Cache de dependencias
COPY go.mod go.sum ./
RUN go mod download

# Copiamos TODO el código fuente
COPY . .

# 2. Build apuntando a la ruta real de tu main
# -trimpath: limpia rutas del host para seguridad
# -ldflags="-s -w": reduce tamaño del binario quitando debug info
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-s -w" -o /bin/booking-service ./cmd/api

# 3. Etapa final (Runtime mínimo)
FROM gcr.io/distroless/base-debian12:nonroot

WORKDIR /

# Copiamos el binario desde la etapa 'build'
COPY --from=build /bin/booking-service /booking-service

EXPOSE 4000
ENV PORT=4000

# Entrypoint directo al binario
ENTRYPOINT ["/booking-service"]
