# football58 local dev

Первый шаг: поднимаем локальный PostgreSQL в Docker. Go API, Dockerfile для API, миграции и CI/CD подключим следующими шагами.

## Что уже подготовлено

- `docker-compose.yml` - локальный PostgreSQL с volume для данных и healthcheck.
- `.env.example` - пример переменных для базы и будущего API.
- `.dockerignore` - исключает мусор из Docker build context.
- `migrations/` - место для будущих миграций базы.

## Структура проекта

Репозиторий GitHub будет инициализирован в корне проекта:

```text
football58/
```

То есть `.git`, `docker-compose.yml`, `.env.example`, README и будущие CI/CD-файлы будут лежать в корне одного репозитория.

Код Go API будем держать отдельно от инфраструктурных файлов:

```text
football58/
  backend/
    cmd/
      api/
    internal/
    go.mod
    go.sum

  migrations/
  docker-compose.yml
  .env.example
  .dockerignore
  README.md
```

Когда появится минимальный API, Dockerfile для него лучше положить рядом с кодом:

```text
backend/Dockerfile
backend/Dockerfile.dev
```

Так будет понятно, что эти Dockerfile относятся именно к Go API, а `docker-compose.yml` в корне отвечает за локальное окружение целиком.

## Деплой без Docker Hub

Образы в Docker Hub публиковать не будем.

План CI/CD такой:

1. GitHub хранит исходный код проекта.
2. GitHub Actions запускает проверки: тесты, линтеры, валидацию конфигов.
3. После успешных проверок GitHub Actions подключается к VPS по SSH.
4. На VPS выполняется `git pull` в директории проекта.
5. На VPS запускается локальная сборка и перезапуск контейнеров через `docker compose up -d --build`.

То есть Docker-образ будет собираться прямо на VPS из исходников, которые приехали из GitHub. Внешний registry на этом этапе не нужен.

## Как поднять базу локально

Можно стартовать сразу с дефолтными значениями:

```bash
docker compose up -d postgres
```

Или сначала создать свой `.env`:

```bash
cp .env.example .env
docker compose up -d postgres
```

Проверить состояние:

```bash
docker compose ps
docker compose logs -f postgres
```

Остановить:

```bash
docker compose down
```

Подключение к базе с хоста:

```text
postgres://football58:football58@localhost:5432/football58?sslmode=disable
```

## Что делаем дальше

1. Инициализируем GitHub-репозиторий в корне `football58/`.
2. Создаем папку `backend/`.
3. Создаем минимальный Go module и `backend/cmd/api`.
4. Добавляем простой HTTP endpoint `/healthz`.
5. Добавляем `backend/Dockerfile.dev` для локальной разработки API.
6. Подключаем API service в `docker-compose.yml`.
7. Добавляем миграции и выбираем инструмент: `golang-migrate`, `goose` или другой.
8. После локальной проверки добавляем GitHub Actions: test, SSH deploy на VPS, `git pull`, `docker compose up -d --build`.
