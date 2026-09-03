# Neurodronism

[English version](./README.md)


AI ghostwriter, который клонирует мой стиль письма и автоматически публикует посты в Telegram-канал.

## Идея

Проект собирает историю моих постов из канала [Дронизм | Авторынок "Ереван"](https://t.me/), строит на их основе модель моего стиля письма и автоматически генерирует новые посты в этом стиле для канала-двойника — [Нейродронизм]. Перед публикацией каждый черновик проходит модерацию: бот присылает его мне в личку, и пост публикуется только после ручного одобрения.

Подробное описание архитектуры и обоснование технических решений — в [`architecture.md`](./architecture.md).

## Статус проекта

MVP работает целиком: новые посты собираются в реальном времени, модель генерирует черновик в стиле канала, после одобрения через Telegram он публикуется автоматически.

- [x] Real-time сбор новых постов канала в Postgres (обработка `channel_post` от Telegram Bot API)
- [x] Бэкфилл истории канала из экспорта Telegram Desktop
- [x] Генерация черновиков в стиле канала (базовый few-shot) — пока на небольшом захардкоженном наборе примеров, точность стиля дорабатывается
- [x] Модерация и публикация через inline-кнопки в личных сообщениях
- [ ] Векторный поиск релевантных постов для few-shot (pgvector) — следующий шаг, чтобы уйти от захардкоженных примеров
- [ ] Автоматический запуск генерации по расписанию (cron) — пока намеренно запускается вручную, пока пайплайн ещё дорабатывается

## Стек

- Go
- PostgreSQL
- Telegram Bot API (`go-telegram-bot-api`)
- Google Gemini API — генерация черновиков
- pgvector — планируется

## Как поднять локально

1. Склонировать репозиторий.

2. Поднять Postgres:
   ```bash
   docker compose up -d
   ```

3. Применить схему (один раз, вручную через `psql`):
   ```sql
   CREATE TABLE posts (
       id            BIGSERIAL PRIMARY KEY,
       tg_message_id BIGINT UNIQUE NOT NULL,
       text          TEXT NOT NULL,
       posted_at     TIMESTAMPTZ NOT NULL,
       is_generated  BOOLEAN NOT NULL DEFAULT false,
       created_at    TIMESTAMPTZ NOT NULL DEFAULT now()
   );
   ```

4. Задать переменные окружения:
   ```bash
   export TELEGRAM_BOT_TOKEN=...
   export DATABASE_URL=postgres://dronism:change_me@localhost:5432/dronism?sslmode=disable
   ```

5. Запустить сборщик постов:
   ```bash
   go run ./cmd/collector
   ```

## Структура репозитория

```
neurodronism/
├── cmd/
│   ├── collector/     # точка входа: сбор постов из Telegram в реальном времени
│   ├── generator/     # точка входа: ручной запуск генерации, для тестов
│   └── data/          # экспорт истории канала (result.json)
├── internal/
│   ├── backfill/      # логика импорта истории
│   ├── generator/     # логика few-shot генерации черновиков
│   └── store/         # работа с Postgres
├── docker-compose.yml
├── architecture.md
└── README.md
```

> На доработку: часть точек входа и внутренней логики пока не до конца разнесена (например, в `cmd/collector` сейчас лежит логика, которой место в `internal/collector`, а в `internal/backfill` — точка входа, которой место в `cmd/backfill`). Пока оставлено как есть, поправится по мере стабилизации проекта.

По мере добавления векторного поиска и планировщика структура и этот README будут обновляться.