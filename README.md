# Task Manager

Программа для управления задачами с тремя интерфейсами: консольный CLI, REST API и Telegram-бот.

Пользователи могут создавать задачи, назначать их другим пользователям, менять статусы и получать уведомления в Telegram.

## Стек

- **Go 1.25**
- **PostgreSQL 17** — хранение задач и пользователей
- **Redis 7** — логирование действий
- **Gin** — HTTP-фреймворк
- **JWT** — авторизация API
- **Swagger** — документация API
- **Telegram Bot API** — управление задачами через бота

## Требования

- Docker и docker-compose — для запуска через контейнеры
- Go 1.25+ — для локального запуска

## Структура проекта

```
.
├── cmd/
│   ├── main.go          # Точка входа: REST API + Telegram Bot
│   └── cli/main.go      # Точка входа: консольный интерфейс
├── internal/
│   ├── handler/         # Транспортный слой (HTTP хендлеры, middleware)
│   ├── service/         # Бизнес-логика
│   ├── repository/      # Слой хранения (PostgreSQL, Redis)
│   ├── bot/             # Telegram бот
│   ├── model/           # Модели данных
│   └── db/              # Подключение к БД
├── migrations/          # SQL миграции
├── docs/                # Swagger документация
├── Dockerfile
├── docker-compose.yml
└── .env
```

## Настройка

Создай `.env` файл в корне проекта:

```env
# PostgreSQL
POSTGRES_USER=admin
POSTGRES_PASSWORD=yourpassword
POSTGRES_DB=otus

# DSN для docker-compose — хост "postgres" (имя сервиса)
POSTGRES_DSN=postgres://admin:yourpassword@postgres:5432/otus?sslmode=disable

# Redis
REDIS_ADDR=redis:6379

# JWT авторизация API
JWT_SECRET=your_secret_key

# Логин/пароль для получения JWT токена
LOGIN=admin
PASSWORD=secret

# Telegram бот (опционально)
TELEGRAM_BOT_TOKEN=your_bot_token
ADMIN_TELEGRAM_ID=your_telegram_id
```

> Для локального запуска без Docker замени хосты в DSN и REDIS_ADDR на `localhost`.

## Запуск

### Через Docker (рекомендуется)

```bash
docker-compose up -d
```

Автоматически поднимает PostgreSQL, Redis, применяет миграции и запускает сервер.

Сервер доступен на `http://localhost:8080`.

### CLI — консольный интерфейс

```bash
# Для локального запуска в .env укажи:
# POSTGRES_DSN=postgres://admin:yourpassword@localhost:5432/otus?sslmode=disable
# REDIS_ADDR=localhost:6379

go run cmd/cli/main.go
```

Пример меню:

```
=== Task Manager ===
1. Список задач
2. Создать задачу
3. Обновить задачу
4. Изменить статус задачи
5. Фильтр по статусу
6. Удалить задачу
7. Список пользователей
0. Выход
```

### Сервер (REST API + Telegram Bot)

```bash
go run cmd/main.go
```

Запускает HTTP-сервер на `:8080` и Telegram-бота одновременно.

## Статусы задач

| Статус | Описание |
|---|---|
| `pending` | Задача создана, ожидает выполнения |
| `in_progress` | Задача взята в работу |
| `done` | Задача выполнена |
| `cancelled` | Задача отменена |

## API

Swagger UI: `http://localhost:8080/swagger/index.html`

### Авторизация

Защищённые эндпоинты требуют JWT токен в заголовке:
```
Authorization: Bearer <token>
```

Получить токен:
```bash
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"login": "admin", "password": "secret"}'
```

### Эндпоинты

| Метод | Путь | Авторизация | Описание |
|---|---|---|---|
| POST | /api/login | — | Получить JWT токен |
| GET | /api/tasks | — | Список всех задач |
| GET | /api/task/:id | — | Получить задачу по ID |
| GET | /api/tasks/status?status= | — | Фильтр задач по статусу |
| GET | /api/users | — | Список пользователей |
| GET | /api/user/:id | — | Получить пользователя по ID |
| GET | /api/user/:id/tasks | — | Задачи пользователя |
| POST | /api/task | ✅ | Создать задачу |
| PUT | /api/task/:id | ✅ | Обновить задачу |
| PUT | /api/task/:id/status | ✅ | Изменить статус задачи |
| DELETE | /api/task/:id | ✅ | Удалить задачу |
| POST | /api/user | ✅ | Создать пользователя |
| PUT | /api/user/:id | ✅ | Обновить пользователя |
| DELETE | /api/user/:id | ✅ | Удалить пользователя |

### Пример создания задачи

```bash
curl -X POST http://localhost:8080/api/task \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"title": "Написать тесты"}'
```

## Telegram Bot

Бот позволяет управлять задачами прямо из Telegram.

### Регистрация пользователей

1. Администратор добавляет пользователя командой: `/adduser @username`
2. Пользователь пишет боту `/start` — доступ активируется автоматически

### Возможности

- Создать задачу и назначить исполнителя
- Выбрать дедлайн (сегодня / завтра / +3 дня / +неделя / без срока)
- Просмотр задач: назначенных мной, назначенных мне
- Смена статуса: взять в работу, закрыть, отменить
- Уведомления исполнителю при назначении задачи
- Уведомления автору при изменении статуса
- Архив выполненных и отменённых задач

## Тесты

```bash
# Юнит-тесты сервисного слоя
go test -race -count=100 ./internal/service/...

# Интеграционные тесты репозитория (нужен PostgreSQL)
POSTGRES_DSN=postgres://admin:yourpassword@localhost:5432/otus?sslmode=disable \
  go test -race -v ./internal/repository/postgres/...
```
