# syntax=docker/dockerfile:1
FROM golang:1.23.4-alpine AS builder
WORKDIR /app
COPY ./api-gateway ./api-gateway
COPY ./order-service ./order-service
COPY ./user-service ./user-service
COPY ./product-service ./product-service
COPY ./payment-service ./payment-service
COPY ./shopping-cart-service ./shopping-cart-service
WORKDIR /app/api-gateway
RUN go mod download
RUN go build -o app ./cmd/main.go

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/api-gateway/app .
EXPOSE 8080
CMD ["./app"]
