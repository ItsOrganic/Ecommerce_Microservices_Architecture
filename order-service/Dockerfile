FROM golang:1.22.5-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o order-service .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/order-service .

EXPOSE 8083

# Correct the CMD path:
CMD ["./order-service"] 