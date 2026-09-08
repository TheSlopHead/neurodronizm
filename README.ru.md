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
- [x] Генерация черновиков с few-shot примерами через векторный поиск (pgvector)
- [x] Модерация и публикация через inline-кнопки в личных сообщениях
- [ ] Автоматический запуск генерации по расписанию (cron) — пока намеренно запускается вручную, пока пайплайн ещё дорабатывается
- [ ] Обернуть сами Go-сервисы в Docker (сейчас в Docker только Postgres) — пока отложено

## Стек

- Go
- PostgreSQL
- Telegram Bot API (`go-telegram-bot-api`)
- Google Gemini API — генерация черновиков и эмбеддинги
- pgvector — поиск похожих постов

## Как поднять локально

1. Склонировать репозиторий.

2. Поднять Postgres:
   ```bash
   docker compose up -d
   ```

3. Применить схему (один раз, вручную через `psql`):
   ```sql
   -- 1. Главная таблица истории постов из Telegram
   CREATE TABLE IF NOT EXISTS posts (
       tg_message_id BIGINT PRIMARY KEY,
       text TEXT NOT NULL,
       posted_at TIMESTAMP WITH TIME ZONE NOT NULL
   );

   -- 2. Таблица для черновиков (очередь модерации)
   CREATE TABLE IF NOT EXISTS drafts (
       id SERIAL PRIMARY KEY,
       text TEXT NOT NULL,
       status VARCHAR(20) DEFAULT 'draft',
       created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
   );

   -- 3. Активация pgvector и таблица векторов
   CREATE EXTENSION IF NOT EXISTS vector;
   CREATE TABLE IF NOT EXISTS post_embeddings (
       post_id BIGINT PRIMARY KEY REFERENCES posts(tg_message_id) ON DELETE CASCADE,
       embedding vector(768) NOT NULL
   );
   ```

4. Задать переменные окружения:
   ```bash
   export TELEGRAM_BOT_TOKEN=...
   export DATABASE_URL=postgres://dronism:change_me@localhost:5432/dronism?sslmode=disable
   ```

5. Запустить бота (real-time сбор + модерация):
   ```bash
   go run .
   ```

6. Разовые команды, по необходимости:
   ```bash
   go run ./cmd/backfill    # импорт истории из экспорта Telegram Desktop
   go run ./cmd/embedder    # эмбеддинг постов, у которых ещё нет вектора
   go run ./cmd/generator   # ручной запуск генерации черновика, для тестов
   ```

## Структура репозитория

```
neurodronism/
├── cmd/
│   ├── backfill/      # точка входа: разовый импорт истории
│   ├── embedder/      # точка входа: эмбеддинг постов без вектора
│   └── generator/     # точка входа: ручной запуск генерации, для тестов
├── data/              # экспорт истории канала (result.json)
├── internal/
│   ├── collector/     # логика сбора постов в реальном времени
│   ├── generator/     # логика few-shot генерации черновиков
│   └── store/         # работа с Postgres
├── main.go            # точка входа: сам работающий бот (сбор + модерация)
├── docker-compose.yml
├── architecture.md
└── README.md
```

По мере добавления планировщика и упаковки сервисов в Docker структура и этот README будут обновляться.