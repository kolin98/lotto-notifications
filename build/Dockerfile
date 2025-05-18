# Build
FROM golang:1.24-alpine AS builder
WORKDIR /app

RUN apk add --no-cache git
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/worker ./cmd/worker

# Run
FROM alpine:latest
WORKDIR /app

COPY --from=builder /app/worker .

RUN adduser -D -g '' appuser \
    && chown -R appuser:appuser /app \
    && chmod -R 700 /app

USER appuser
