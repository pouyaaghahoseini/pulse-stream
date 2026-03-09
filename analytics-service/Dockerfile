FROM golang:1.26 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o analytics-service ./cmd/worker

FROM debian:bookworm-slim

WORKDIR /app

COPY --from=builder /app/analytics-service .

EXPOSE 8081

CMD ["./analytics-service"]