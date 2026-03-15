FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o order-service .

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/order-service .
EXPOSE 8082
CMD ["./order-service"]