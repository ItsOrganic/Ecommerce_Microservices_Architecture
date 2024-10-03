FROM golang:1.22.5-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o gpql-gateway .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/gpql-gateway .

EXPOSE 8080

# Correct the CMD path:
CMD ["./gpql-gateway"] 