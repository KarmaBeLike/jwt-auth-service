# Используем официальный образ Go
FROM golang:1.23

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем файлы модулей
COPY go.mod go.sum ./

# Загружаем зависимости
RUN go mod download

# Копируем остальной проект
COPY . .

# Собираем приложение из папки cmd/
RUN go build -o auth-service ./cmd

# Запускаем приложение
CMD ["./auth-service"]
