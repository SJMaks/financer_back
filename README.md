# Financer API

REST API сервер для приложения учёта личных финансов: счета, категории, доходы и расходы, помесячная статистика. Pet-проект, написанный для закрепления навыков backend-разработки на Go: чистая архитектура, работа с PostgreSQL на уровне триггеров и представлений, JWT-аутентификация, тесты с моками и автодокументирование API через Swagger.

Сайт доступен по адресу: www.financerapp.ru

Схема API доступна по адресу: [api.financerapp.ru/swagger/index.html](https://api.financerapp.ru/swagger/index.html)

---

## Содержание

- [Возможности](#возможности)
- [Стек](#стек)
- [Архитектура](#архитектура)
- [Схема базы данных](#схема-базы-данных)
- [Аутентификация](#аутентификация)
- [API](#api)
- [Переменные окружения](#переменные-окружения)
- [Локальный запуск](#локальный-запуск)
- [Миграции](#миграции)
- [Тесты](#тесты)
- [Swagger](#swagger)
- [Планы развития](#планы-развития)

---

## Возможности

- **Пользователи** - регистрация с валидацией (логин, пароль со сложностью, ФИО, валюта), логин, получение/обновление/удаление профиля.
- **Счета** - несколько счетов на пользователя (например, наличные и карта); баланс счёта автоматически пересчитывается триггерами при любых изменениях операций.
- **Категории операций** - пользовательские категории доходов и расходов; при регистрации каждого пользователя автоматически создаются системные категории `no_category`, которые нельзя удалить или изменить.
- **Операции** - доходы и расходы, привязанные к счёту и категории; фильтрация (по типу, категории, диапазону дат), пагинация (`limit`/`offset`), сортировка (по `exec_at` или `sum`, `asc`/`desc`).
- **Статистика** - помесячная статистика доходов/расходов/прибыли по всем операциям пользователя и по отдельному счёту (выполняется SQL-представлениями на стороне БД).
- **Оптимистическая блокировка** - все сущности имеют поле `version`; обновление возможно только при совпадении версии, что исключает потерю изменений при конкурентных запросах.
- **Консистентность данных на уровне БД** - инварианты баланса, знак суммы и перенос операций при удалении категории обеспечиваются триггерами PostgreSQL.
- **Аутентификация и авторизация** - JWT (HS256, срок жизни 24 часа), middleware проверки токена, изоляция данных по пользователю.
- **CORS** - белый список origin'ов, обработка предварительных OPTIONS-запросов.
- **Swagger UI** - интерактивная документация, генерируется аннотациями `swaggo` из кода хендлеров.
- **Graceful shutdown** - корректное завершение по SIGINT/SIGTERM с таймаутом 10 секунд.
- **Логирование** - `zap` с двумя синками: stdout (Info) и файл в `_out/logs/` (Debug).

## Стек

| Компонент | Технология |
|---|---|
| Язык | Go 1.26 |
| HTTP-фреймворк | [Gin](https://github.com/gin-gonic/gin) v1.12 |
| База данных | PostgreSQL 18 |
| Драйвер БД | [pgx/v5](https://github.com/jackc/pgx) + `pgxpool` (пул соединений: min 5 / max 20) |
| Миграции | [golang-migrate](https://github.com/golang/migrate) v4.20.1 (запуск через Docker) |
| Аутентификация | [golang-jwt/jwt/v5](https://github.com/golang-jwt/jwt) (HS256), [bcrypt](https://pkg.go.dev/golang.org/x/crypto/bcrypt) |
| Логирование | [uber-go/zap](https://go.uber.org/zap) |
| Документация API | [swaggo/swag](https://github.com/swaggo/swag) v1.16.6 + Swagger UI |
| Тесты | [testify](https://github.com/stretchr/testify), [mockery](https://github.com/vektra/mockery) v3 (`.mockery.yml`) |
| Инфраструктура | Docker, Docker Compose, Makefile |
| CI/CD | GitHub (feature-branches → PR → `main`) |

## Архитектура

Код разбит на **фичи (feature slices)**; внутри каждой фичи - слои `domain → repository → service → transport` с зависимостями только внутрь (transport зависит от service, service - от интерфейсов репозиториев, репозитории - от domain):

```
cmd/financer/          - точка входа (main.go) и Dockerfile приложения
internal/
  core/                - инфраструктурный код, не привязанный к бизнес-логике
    app/               - сборка приложения: роутер, middleware, все фичи, Swagger
    apperrors/         - типизированные ошибки приложения
    auth/              - JWT TokenManager (Generate/Parse), TokenPayload
    config/            - загрузка конфигурации из переменных окружения
    database/          - инициализация pgxpool с пулом соединений
    httpserver/        - обёртка над http.Server (Start/Stop)
    logger/            - инициализация zap (stdout + файл)
    middleware/        - CORS и проверка JWT (CheckAuth, UserFromContext)
    utils/             - единый формат ошибок HTTP ({error, code})
  features/            - бизнес-функции, изолированные друг от друга
    users/             - регистрация, логин, профиль
    accounts/          - счета пользователя
    categories/        - категории доходов/расходов
    operation_types/   - справочник типов операций (income/expense)
    operations/        - доходы/расходы, фильтры, пагинация, сортировка
    statistics/        - помесячная статистика (SQL-представления)
      domain/          - модели и бизнес-константы
      repository/      - SQL-запросы (pgx), реализация интерфейсов
      service/         - бизнес-логика, валидация, работа с интерфейсами
      transport/       - Gin-хендлеры, DTO, маппинг, роутер
migrations/            - SQL-миграции (up/down)
docs/                  - сгенерированный Swagger (docs.go, swagger.json, swagger.yaml)
docker-compose.yaml    - локальное окружение: Postgres, миграции, Swagger, port-forwarder
Makefile               - все команды запуска и обслуживания
```

Принципы, заложенные в архитектуру:

- **Инверсия зависимостей** - сервисы зависят только от интерфейсов репозиториев и токен-менеджера, что позволяет тестировать их моками без БД.
- **DTO на транспортном слое** - доменные модели не сериализуются наружу напрямую; маппинг вынесен в `mapping.go`.
- **Единый формат ошибок** - все ошибки отдаются как `{"error": "...", "code": "..."}` через `utils.WriteHTTPError`; ошибки 500 логируются, детали не утекают наружу.
- **Тонкие хендлеры, вся логика в сервисах** - хендлеры парсят запрос и делегируют; сервисы валидируют и оркестрируют.

## Схема базы данных

**Таблицы:** `users`, `accounts`, `operation_types`, `operation_categories`, `operations` - с внешними ключами, `ON DELETE CASCADE`, индексами под частые запросы и полем `version BIGINT` под оптимистическую блокировку.

**Триггеры** (инварианты на уровне БД):

1. `trg_update_balance` (`INSERT/UPDATE/DELETE` на `operations`) - автоматически пересчитывает `accounts.balance` и инкрементирует `version` счёта при любых изменениях операций; смена `account_id` у операции запрещена.
2. `trg_operations_signed_sum` (`BEFORE INSERT/UPDATE`) - нормализует знак суммы: расходы хранятся со знаком минус, доходы - с плюсом.
3. `trg_reassign_operations_before_category_delete` (`BEFORE DELETE` на `operation_categories`) - перед удалением категории переносит все её операции в `no_category` соответствующего типа, чтобы не нарушить внешний ключ.

**Представления** для статистики: `v_account_monthly_stats` и `v_user_monthly_stats` - агрегация доходов, расходов и прибыли по месяцам одним запросом на стороне БД.

Бизнес-правила (нельзя удалять/изменять `no_category`, владение ресурсом проверяется по `user_id` в каждом SQL-запросе) дополнительно контролируются на уровне сервисов и репозиториев.

## Аутентификация

- Регистрация: пароль хешируется `bcrypt` (DefaultCost) и валидируется по правилам: 5–50 символов, латиница/цифры, обязательны заглавная буква и цифра.
- Логин выдаёт JWT (HS256) со сроком жизни **24 часа**; payload содержит `id` и `login` пользователя.
- Защищённые эндпоинты требуют заголовок `Authorization: Bearer <token>`; middleware кладёт payload в контекст Gin, а каждый репозиторий фильтрует данные по `user_id` - чужие данные недоступны даже по прямому ID.

## API

Базовый путь: `/api/v1`. Полная интерактивная документация - **[Swagger UI](http://localhost:8000/swagger/index.html)** после запуска.

| Метод | Путь | Описание | Доступ |
|---|---|---|---|
| `POST` | `/auth/register` | Регистрация пользователя | публичный |
| `POST` | `/auth/login` | Логин, возврат JWT | публичный |
| `GET` | `/users/me` | Текущий пользователь | Bearer |
| `PATCH` | `/users/me` | Обновить профиль | Bearer |
| `DELETE` | `/users/me` | Удалить аккаунт | Bearer |
| `POST` | `/accounts` | Создать счёт | Bearer |
| `GET` | `/accounts` | Список счетов | Bearer |
| `GET` | `/accounts/{id}` | Счёт по ID | Bearer |
| `PATCH` | `/accounts/{id}` | Обновить счёт (optimistic lock) | Bearer |
| `DELETE` | `/accounts/{id}` | Удалить счёт | Bearer |
| `POST` | `/categories` | Создать категорию | Bearer |
| `GET` | `/categories` | Список категорий | Bearer |
| `GET` | `/categories/{id}` | Категория по ID | Bearer |
| `PATCH` | `/categories/{id}` | Обновить категорию (optimistic lock) | Bearer |
| `DELETE` | `/categories/{id}` | Удалить категорию | Bearer |
| `GET` | `/operation_types` | Справочник типов операций | Bearer |
| `POST` | `/accounts/{accountID}/operations` | Создать операцию | Bearer |
| `GET` | `/accounts/{accountID}/operations` | Список операций (фильтры, пагинация, сортировка) | Bearer |
| `GET` | `/accounts/{accountID}/operations/{operationID}` | Операция по ID | Bearer |
| `PATCH` | `/accounts/{accountID}/operations/{operationID}` | Обновить операцию (optimistic lock) | Bearer |
| `DELETE` | `/accounts/{accountID}/operations/{operationID}` | Удалить операцию | Bearer |
| `GET` | `/statistics` | Помесячная статистика по всем операциям пользователя | Bearer |
| `GET` | `/accounts/{accountID}/statistics` | Помесячная статистика по счёту | Bearer |

Пример параметров `GET /accounts/{accountID}/operations`: `limit` (по умолчанию 50), `offset`, `type` (`income`|`expense`), `category_id`, `min_date`, `max_date`, `sort` (`exec_at`|`sum`), `order` (`asc`|`desc`).

## Переменные окружения

| Переменная | Описание | Пример |
|---|---|---|
| `HTTP_ADDRESS` | Адрес и порт HTTP-сервера (`host:port`) | `localhost:8000`, в Docker - `0.0.0.0:8000` |
| `CORS_ALLOWED_ORIGINS` | Разрешённые CORS Origins через запятую | `http://localhost:5173` |
| `JWT_SECRET` | Секрет для подписи JWT (HS256) | случайная строка ≥ 32 символов |
| `DB_HOST` | Хост PostgreSQL | `localhost` / `financer_postgres` |
| `DB_PORT` | Порт PostgreSQL | `5432` |
| `DB_USER` | Пользователь БД | `financer` |
| `DB_PASSWORD` | Пароль БД | - |
| `DB_NAME` | Имя базы данных | `financer` |

Все переменные **обязательны** - приложение не стартует без любой из них. Шаблон: [.env.example](.env.example). Файл `.env` добавлен в `.gitignore`.


## Локальный запуск

Требования: Docker с Docker Compose и Make.

```bash
# 1. Клонировать репозиторий и перейти в каталог
git clone https://github.com/SJMaks/financer_back.git
cd financer_back

# 2. Создать файл окружения
cp .env.example .env
# заполнить .env: HTTP_ADDRESS=localhost:8000, JWT_SECRET=<случайная строка>,
# DB_HOST=localhost, DB_PORT=<свободный порт, например 5433>,
# DB_USER=financer, DB_PASSWORD=<пароль>, DB_NAME=financer

# 3. Поднять PostgreSQL
make env-up

# 4. Применить миграции
make migrate-up

# 5. Пробросить порт Postgres на localhost
make env-port-forward

# 6. Запустить приложение
make run
# сервер: http://localhost:8000, Swagger: http://localhost:8000/swagger/index.html
```

Опционально:

```bash
make env-down           # остановить окружение
make env-cleanup        # остановить и удалить данные БД
make financer-deploy    # запустить приложение в Docker-контейнере (порт 8000)
make docker-ps          # статус контейнеров
```

## Миграции

SQL-миграции лежат в [`migrations/`](migrations/) и управляются [golang-migrate](https://github.com/golang/migrate) через Docker (локальная установка CLI не нужна):

```bash
make migrate-create seq=<name>       # создать пару up/down файлов
make migrate-up                      # применить все миграции
make migrate-down                    # откатить последнюю миграцию
make migrate-action action=<action>  # выполнить определённое действие
```

Текущая схема: `000001_init` (таблицы, индексы, представления, справочник типов), `000002_account_balance_trigger` (пересчёт баланса), `000004_delete_insert_triggers` (перенос операций при удалении категории, нормализация знака суммы).

## Тесты

```bash
go test ./...   # или go test ./internal/features/users/...
```

- **Сервисный слой** (`internal/features/users/service/test/`) - юнит-тесты `Register`/`Login`/CRUD на моках `UserRepository` и `TokenManager`.
- **Транспортный слой** (`internal/features/users/transport/test/`) - тесты Gin-хендлеров (статус-коды, валидация входа, формат ошибок) на моках сервиса и auth-middleware.

Моки генерируются [mockery](https://github.com/vektra/mockery) v3 по конфигурации [`.mockery.yml`](.mockery.yml):

```bash
mockery   # сгенерировать/обновить моки в пакетах */test/
```

Проверка перед коммитом: `go build ./... && go vet ./... && go test ./...`.

## Swagger

Документация генерируется из аннотаций `swaggo` над хендлерами и метаданных в `main.go`:

```bash
make swagger-gen
```

Результат - папка [`docs/`](docs/) (`docs.go`, `swagger.json`, `swagger.yaml`), которая подключается пустым импортом в `main.go` и отдаётся по `/swagger/index.html`. После изменения хендлеров документацию нужно перегенерировать и закоммитить.

### Дальнейшие улучшения (по желанию)

- Вынести список разрешённых CORS-origin'ов из кода (`internal/core/config/config.go`) в переменные окружения, чтобы указывать домен развёрнутого фронтенда без пересборки.

## Планы развития

- Refresh-токены (сейчас access-токен живёт 24 часа без продления).
- Регулярные (recurring) операции и бюджеты с уведомлениями.
- Операция перевода между счетами.
- Автоматическое распределение доходов между специально настроенными счетами (метод конвертов).
- Экспорт статистики (CSV/PDF) и графики по месяцам.
- Интеграционные тесты на реальной БД (testcontainers).
