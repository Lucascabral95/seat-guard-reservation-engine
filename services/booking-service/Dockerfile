# syntax=docker/dockerfile:1

FROM golang:1.25 AS build
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-s -w" -o /bin/booking-service ./cmd/api

FROM gcr.io/distroless/base-debian12:nonroot

WORKDIR /

COPY --from=build /bin/booking-service /booking-service

EXPOSE 4000
ENV PORT=4000

ENTRYPOINT ["/booking-service"]
