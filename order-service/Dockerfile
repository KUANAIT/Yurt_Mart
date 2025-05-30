# syntax=docker/dockerfile:1
FROM golang:1.23.0-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o app ./cmd/main.go

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/app .
EXPOSE 50051
CMD ["./app"]

