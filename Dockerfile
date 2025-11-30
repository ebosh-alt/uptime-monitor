FROM golang:1.24-bullseye AS builder

WORKDIR /app

# Ускоряем скачивание зависимостей.
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Собираем статический бинарник.
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o admin_historium ./cmd

FROM gcr.io/distroless/base-debian12

WORKDIR /app

COPY --from=builder /app/uptime-monitor /app/uptime-monitor
COPY --chown=nonroot:nonroot config /app/config

EXPOSE 3000

USER nonroot

ENTRYPOINT ["/app/uptime-monitor"]
