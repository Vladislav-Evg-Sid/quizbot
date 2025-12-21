# Dockerfile
FROM golang:1.25-alpine

WORKDIR /app

# 1. Устанавливаем air для hot reload
RUN go install github.com/air-verse/air@latest

# 2. Копируем ТОЛЬКО зависимости для кэширования
COPY go.mod go.sum ./

# 3. Скачиваем зависимости
RUN go mod download

# 4. Копируем ВЕСЬ проект (включая сгенерированные protobuf файлы)
COPY . .

# 5. Запускаем с hot reload
CMD ["air", "-c", ".air.toml"]