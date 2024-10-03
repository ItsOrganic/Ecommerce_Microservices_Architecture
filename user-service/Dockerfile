FROM golang:1.22.5-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o user-service .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/user-service .

EXPOSE 8081

# Correct the CMD path:
CMD ["./user-service"] 