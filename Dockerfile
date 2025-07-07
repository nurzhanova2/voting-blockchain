# 1. Используем официальный образ Go
FROM golang:1.24.3-alpine AS builder

# 2. Устанавливаем необходимые утилиты
RUN apk add --no-cache git

# 3. Устанавливаем рабочую директорию
WORKDIR /app

# 4. Копируем go.mod и go.sum и ставим зависимости
COPY go.mod go.sum ./
RUN go mod download

# 5. Копируем весь проект
COPY . .

# 6. Собираем бинарник
RUN go build -o server ./cmd/main.go

# === Минималистичный образ для запуска ===
FROM alpine:latest

# 7. Устанавливаем certs (для https)
RUN apk --no-cache add ca-certificates

# 8. Рабочая директория
WORKDIR /root/

# 9. Копируем бинарник из стадии builder
COPY --from=builder /app/server .

# 10. Копируем .env файл
COPY .env .

# 11. Устанавливаем переменные среды (опционально)
ENV GIN_MODE=release

# 12. Команда запуска
CMD ["./server"]

#13

RUN go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.59.0
