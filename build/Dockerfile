# Build
FROM golang:1.24-alpine AS builder
WORKDIR /app

RUN apk add --no-cache git gcc musl-dev sqlite-dev
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=1 GOOS=linux go build -o /app/worker ./cmd/worker

# Run
FROM alpine:latest
WORKDIR /app

COPY --from=builder /app/worker .

RUN adduser -D -g '' appuser \
    && chown -R appuser:appuser /app \
    && chmod -R 700 /app

USER appuser
