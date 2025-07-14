# Используем официальный образ Go 1.24 для сборки
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Установка зависимостей, в том числе git для go mod
RUN apk add --no-cache git

# Копируем модули и кэшируем зависимости
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходники
COPY . .

# Собираем бинарник
RUN go build -o app ./cmd/app

# Финальный образ — минимальный alpine с бинарником
FROM alpine:latest

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /app/app .

EXPOSE 50051

CMD ["./app"]
