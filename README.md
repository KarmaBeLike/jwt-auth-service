# JWT Auth Service

Сервис аутентификации, предоставляющий REST API для работы с JWT токенами доступа и токенами обновления.

## Технологии

- **Go**: основной язык программирования
- **JWT**: формат токенов доступа
- **PostgreSQL**: хранение данных о токенах обновления
- **bcrypt**: хеширование токенов обновления
- **Gin**: HTTP веб-фреймворк

## Основная функциональность

Сервис предоставляет два основных REST маршрута:

1. **POST /token**: Генерирует новую пару токенов (Access + Refresh) для пользователя с указанным GUID
2. **POST /refresh**: Обновляет пару токенов, используя действующий Refresh токен

## Характеристики токенов

### Access Token

- Формат: JWT
- Алгоритм подписи: SHA512
- Не хранится в базе данных
- Содержит информацию о пользователе и IP-адресе клиента

### Refresh Token

- Случайно сгенерированная строка
- Формат передачи: base64
- Хранится в базе данных только в виде bcrypt хеша
- Защищен от изменения на стороне клиента и повторного использования
- Связан с конкретным Access токеном

## Особенности безопасности

- Обоюдная связь Access и Refresh токенов
- Проверка IP-адреса клиента при операции Refresh
- Автоматическое уведомление по email при обнаружении изменения IP-адреса
- Защита от повторного использования Refresh токенов

## Запуск и использование

### Предварительные требования

- Go 1.22+
- PostgreSQL 14+
- Настроенный .env файл (пример в .env.example)

### Установка

```bash
# Клонирование репозитория
git clone https://github.com/KarmaBeLike/jwt-auth-service.git
cd jwt-auth-service

```

### Конфигурация

Создайте файл `.env` на основе примера в .env.example

### Запуск

```
docker-compose up --build
```

## API Endpoints

r.POST("/token")
r.POST("/refresh")

### Генерация токенов

```
POST /token?user_id={guid}
```

#### Параметры запроса:
- `user_id`: GUID пользователя (обязательный параметр)

#### Ответ:
```json
{
  "access_token": "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "ZjJkNjQ5MzAtNGZhYS00ZGZmLWJmNDEtMTNhZjVlOWJhOWQ1"
}
```

### Обновление токенов

```
POST /refresh
```

#### Тело запроса:
```json
{
  "refresh_token": "ZjJkNjQ5MzAtNGZhYS00ZGZmLWJmNDEtMTNhZjVlOWJhOWQ1"
}
```

#### Ответ:
```json
{
  "access_token": "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "YzMwMmQ3ZTAtMzJlYy00YTY2LWFlMTYtNWZmMzE2ZWVjNmIy"
}
```

## Структура проекта

```
jwt-auth-service/
├── cmd/
│   └── main.go/        # Точка входа приложения
├── config/             # Конфигурация приложения
├── internal/
│   ├── dto/            # Объекты передачи данных
│   ├── handler/        # HTTP обработчики
│   ├── database/       # Подключение к базе данных
│   ├── model/          # Модели данных
│   ├── repository/     # Слой доступа к базе данных
│   └── service/        # Бизнес-логика
├── pkg/                # Общие пакеты
│             
├── go.mod
├── go.sum
└── README.md
```
