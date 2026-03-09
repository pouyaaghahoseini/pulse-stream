FROM golang:1.26 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o api-service ./cmd/api

FROM debian:bookworm-slim

WORKDIR /app

COPY --from=builder /app/api-service .

EXPOSE 8080

CMD ["./api-service"]