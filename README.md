# Neurodronism

[Русская версия](./README.ru.md)

An AI ghostwriter that clones my writing style and automatically publishes posts to a Telegram channel.

## Idea

The project collects the history of my posts from [Дронизм | Авторынок "Ереван"](https://t.me/dronizm), builds a model of my writing style from them, and automatically generates new posts in that style for a companion channel, [Нейродронизм]. Every draft goes through moderation before publishing: the bot sends it to me in a private message, and it only gets posted after manual approval.

Full architecture and reasoning behind the technical decisions: [`architecture.md`](./architecture.md).

## Status

- [x] Real-time ingestion of new channel posts into Postgres (handling `channel_post` updates from the Telegram Bot API)
- [x] Backfill of channel history from a Telegram Desktop export
- [ ] Draft generation in the channel's style (few-shot / RAG)
- [ ] Moderation and publishing via inline buttons in private messages
- [ ] Vector search for relevant few-shot examples (pgvector)

## Stack

- Go
- PostgreSQL
- Telegram Bot API (`go-telegram-bot-api`)
- LLM API (OpenAI / Anthropic) — planned
- pgvector — planned

## Running locally

1. Clone the repository.

2. Start Postgres:
   ```bash
   docker compose up -d
   ```

3. Apply the schema (one-time, via `psql`):
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

4. Set environment variables:
   ```bash
   export TELEGRAM_BOT_TOKEN=...
   export DATABASE_URL=postgres://dronism:change_me@localhost:5432/dronism?sslmode=disable
   ```

5. Run the collector:
   ```bash
   go run ./cmd/collector
   ```

## Repository structure

```
neurodronism/
├── cmd/
│   └── collector/     # entry point: real-time post collection from Telegram
├── internal/
│   └── store/         # Postgres access layer
├── docker-compose.yml
├── architecture.md
└── README.md
```

Structure and this README will be updated as generation, moderation, and the scheduler get added.
