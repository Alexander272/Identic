# Структура проекта Identic Backend

## Общая архитектура (Clean Architecture)

```
backend/
├── cmd/app/           # Точка входа приложения
├── configs/           # Конфигурационные файлы
├── internal/          # Приватный код приложения
│   ├── access/        # Проверка прав доступа (Casbin)
│   ├── config/        # Загрузка конфигурации
│   ├── constants/     # Константы
│   ├── events/        # Event-driven система
│   ├── migrate/       # Миграции базы данных
│   ├── models/        # Модели данных (DTO, доменные объекты)
│   ├── repository/    # Слой доступа к данным
│   ├── server/        # Инициализация сервера
│   ├── services/     # Бизнес-логика
│   └── transport/     # Транспортный слой (HTTP, WebSocket)
└── pkg/              # Публичные пакеты (не зависят от internal)
```

---

## cmd/app/

**`main.go`** — точка входа. Инициализирует конфигурацию, базы данных (PostgreSQL, Redis), сервер.

---

## configs/

Конфигурационные файлы (YAML/JSON) для различных окружений.

---

## internal/access/

Система авторизации на основе **Casbin**:
- `builder.go` — построение политик
- `permissions.go` — проверка прав
- `registry.go` — реестр политик
- `types.go` — типы для авторизации

---

## internal/config/

**`config.go`** — загрузка и парсинг конфигурации из YAML файлов и переменных окружения.

---

## internal/constants/

Константы проекта:
- `context.go` — ключи контекста
- `cookie.go` — настройки cookies
- `import.go` — константы для импорта

---

## internal/events/

Event-driven архитектура для уведомления об изменениях политик:
- `policy.go` — события изменения политик Casbin

---

## internal/migrate/

Миграции базы данных (используется **goose**):
```
migrate/postgres/migrations/
├── YYYYMMDDHHMMSS_name.sql  # Миграции вверх и вниз
└── migrations.go             # Регистрация миграций
```

---

## internal/models/

Модели данных — DTO (Data Transfer Objects) и доменные объекты:

| Файл | Описание |
|------|----------|
| `activity.go` | Логи изменений заказов/позиций |
| `audit.go` | Аудит политик безопасности |
| `audit_log.go` | Модель аудит-лога |
| `errors.go` | Кастомные ошибки приложения |
| `import.go` | Модели для импорта данных |
| `orders.go` | Заказы |
| `params.go` | Параметры запросов |
| `permissions.go` | Роли и разрешения |
| `policies.go` | Политики доступа |
| `positions.go` | Позиции заказов |
| `response.go` | Стандартные ответы API |
| `roles.go` | Роли пользователей |
| `role_hierarchy.go` | Иерархия ролей |
| `search.go` | Запросы и результаты поиска |
| `search_log.go` | Логи поисковых запросов |
| `session.go` | Сессии пользователей |
| `subscribe.go` | Подписки на обновления |
| `user_login.go` | История входов |
| `users.go` | Пользователи |

---

## internal/repository/

Слой доступа к данным. Реализует паттерн **Repository**:

```
repository/
├── repo.go           # Агрегирующий репозиторий (внедряется в services)
├── postgres/          # Реализация для PostgreSQL
│   ├── activity.go          # Логи активности
│   ├── audit_log.go         # Аудит-логи
│   ├── order_events.go       # События заказов
│   ├── orders.go            # Заказы
│   ├── permissions.go       # Разрешения
│   ├── positions.go         # Позиции
│   ├── role_hierarchy.go   # Иерархия ролей
│   ├── roles.go            # Роли
│   ├── search.go           # Поиск (exact/fuzzy)
│   ├── search_log.go       # Логи поиска
│   ├── tables.go           # Имена таблиц
│   ├── transactions.go     # Управление транзакциями
│   ├── user_login.go       # История входов
│   ├── users.go            # Пользователи
│   └── utils.go            # QueryBuilder для динамических запросов
└── redis/
    └── search.go           # Redis-кэш для поиска
```

**`postgres/utils.go`** содержит `QueryBuilder` — построитель SQL-запросов с поддержкой:
- Фильтров (eq, neq, gte, lte, con, like, in, null)
- Пагинации (limit, offset)
- Сортировки (single, multi-sort)
- Курсорной пагинации

---

## internal/server/

**`server.go`** — инициализация HTTP/WebSocket сервера с Gin.

---

## internal/services/

Бизнес-логика приложения. Каждый сервис инкапсулирует логику одной предметной области:

| Файл | Описание |
|------|----------|
| `access_policies.go` | Управление политиками Casbin |
| `activity.go` | Работа с логами активности |
| `adapter.go` | Адаптер для Casbin |
| `audit_log.go` | Аудит-логи |
| `import.go` | Импорт данных из файлов |
| `order_stream.go` | WebSocket стриминг обновлений заказов |
| `orders.go` | Управление заказами |
| `permissions.go` | Разрешения |
| `positions.go` | Управление позициями |
| `role_hierarchy.go` | Иерархия ролей |
| `roles.go` | Роли |
| `runner.go` | Фоновые задачи |
| `search.go` | Поиск заказов (exact/fuzzy) |
| `search_log.go` | Логирование поисковых запросов |
| `search_stream.go` | WebSocket стриминг результатов поиска |
| `services.go` | Агрегирующий сервис (внедряется в transport) |
| `session.go` | Управление сессиями |
| `transaction.go` | Менеджер транзакций |
| `user_login.go` | История входов |
| `users.go` | Пользователи |

---

## internal/transport/

Транспортный слой — HTTP и WebSocket обработчики.

### HTTP API

```
transport/http/
├── handler.go              # Основной HTTP хендлер
├── utils/                  # Утилиты для HTTP
└── handlers/
    ├── audit/              # Аудит-логи
    ├── auth/               # Аутентификация (Keycloak)
    ├── import_file/        # Импорт файлов
    ├── orders/             # Заказы
    ├── permissions/        # Разрешения
    ├── positions/          # Позиции
    ├── roles/              # Роли
    ├── search/             # Поиск
    └── users/              # Пользователи
```

### Middleware

```
transport/middleware/
├── identity.go             # Извлечение identity пользователя
├── middleware.go          # Общие middleware (CORS, логирование)
└── permissions.go          # Проверка прав доступа
```

### WebSocket

```
transport/ws/
├── handler.go             # WebSocket хендлер (Hub, reader/writer)
├── router/                # Маршрутизация WS сообщений
│   └── router.go          # Роутер по action
├── search/                # Поиск через WS
├── search_logs/           # Получение логов поиска
└── subscribe/             # Подписки на обновления
```

---

## pkg/

Публичные, переиспользуемые пакеты:

| Директория | Описание |
|------------|----------|
| `auth/` | Клиент Keycloak |
| `database/postgres/` | Подключение к PostgreSQL |
| `database/redis/` | Подключение к Redis |
| `error_bot/` | Отправка ошибок в Telegram бот |
| `hasher/` | Хеширование паролей (bcrypt) |
| `limiter/` | Rate limiting |
| `logger/` | Логирование |
| `ws_hub/` | WebSocket Hub (менеджер соединений) |

---

## Схема зависимостей (Dependency Flow)

```
HTTP/WS Request
       │
       ▼
transport/handlers     # Парсит запросы, валидирует
       │
       ▼
services/             # Бизнес-логика
       │
       ▼
repository/           # Доступ к данным (БД, кэш)
       │
       ▼
PostgreSQL / Redis     # Хранение данных
```

---

## Основные технологии

- **Web Framework**: Gin
- **База данных**: PostgreSQL (pgx)
- **Кэш**: Redis
- **Авторизация**: Casbin + Keycloak
- **Миграции**: goose
- **Логирование**: zerolog
- **WebSocket**: gorilla/websocket
- **Тестирование**: testify
