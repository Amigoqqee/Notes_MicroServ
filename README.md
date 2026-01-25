# microServ

Небольшой учебный проект с двумя микросервисами на Go:

- **Auth** — регистрация, логин, JWT, bcrypt.
- **Notes** — CRUD заметок с авторизацией через JWT.

Общий вход через **Nginx**. Хранилища: **PostgreSQL** (Auth), **MongoDB** и **Redis** (Notes).

## Архитектура

- **auth** → PostgreSQL
- **notes** → MongoDB + Redis
- **nginx** → проксирует `/auth` и `/notes`

## Быстрый старт (Docker)

1. Скопируйте файл окружения:

```bash
cp .env.example .env
```

2. Убедитесь, что установлен Docker и Docker Compose.
3. Запустите:

```bash
docker-compose up --build -d
```

После запуска:

- Auth: http://localhost/auth/
- Notes: http://localhost/notes/
- PgAdmin: http://localhost:8102
- Mongo Express: http://localhost:8081

Порты и учётные данные можно изменить в [.env](.env).

## Запуск без Docker

> Нужны локально поднятые PostgreSQL, MongoDB и Redis.

1. Скопируйте файл окружения:

```bash
cp .env.example .env
```

- Auth:

```bash
bash auth/start_auth.sh
```

- Notes:

```bash
bash notes/start_notes.sh
```

## Тестирование API

Есть скрипт для полного прогона регистрации/логина и CRUD заметок:

```bash
bash global_test.sh
```

Скрипт использует базовый URL `http://localhost` и ожидает, что прокси Nginx доступен на порту 80.

## Endpoints

### Auth (`/auth`)

| Method | Path | Description |
| --- | --- | --- |
| POST | `/auth/register` | Регистрация пользователя |
| POST | `/auth/login` | Логин, выдача JWT |
| POST | `/auth/refresh` | Обновление access-токена |
| GET | `/auth/user` | Получить профиль (JWT) |
| PUT | `/auth/user` | Обновить профиль (JWT) |
| DELETE | `/auth/user` | Удалить пользователя (JWT) |

### Notes (`/notes`)

| Method | Path | Description |
| --- | --- | --- |
| POST | `/notes/note` | Создать заметку |
| GET | `/notes/note/:id` | Получить заметку |
| PUT | `/notes/note/:id` | Обновить заметку |
| DELETE | `/notes/note/:id` | Удалить заметку |
| GET | `/notes/notes` | Получить все заметки |

### Авторизация

Для защищённых эндпоинтов добавляйте заголовок:

```
Authorization: Bearer <access_token>
```

## Структура репозитория

- [auth](auth) — сервис авторизации
- [notes](notes) — сервис заметок
- [nginx](nginx) — конфигурация прокси
- [pkg/jwtmanager](pkg/jwtmanager) — общий пакет для JWT
