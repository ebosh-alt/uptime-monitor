FROM golang:1.24-bullseye AS builder

WORKDIR /app

# Ускоряем скачивание зависимостей.
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Собираем статический бинарник. Имя совпадает с тем, что копируем во второй стадии.
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o uptime-monitor ./cmd

FROM gcr.io/distroless/base-debian12

WORKDIR /app

COPY --from=builder /app/uptime-monitor /app/uptime-monitor
COPY --chown=nonroot:nonroot config /app/config
COPY --chown=nonroot:nonroot db /app/db

EXPOSE 8080

USER nonroot

ENTRYPOINT ["/app/uptime-monitor"]
