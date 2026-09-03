# Neurodronism

[Русская версия](./README.ru.md)

An AI ghostwriter that clones my writing style and automatically publishes posts to a Telegram channel.

## Idea

The project collects the history of my posts from [Дронизм | Авторынок "Ереван"](https://t.me/), builds a model of my writing style from them, and automatically generates new posts in that style for a companion channel, [Нейродронизм]. Every draft goes through moderation before publishing: the bot sends it to me in a private message, and it only gets posted after manual approval.

Full architecture and reasoning behind the technical decisions: [`architecture.md`](./architecture.md).

## Status

MVP is working end-to-end: new posts are collected in real time, the model generates a draft in the channel's style, and after approval via Telegram it gets published automatically.

- [x] Real-time ingestion of new channel posts into Postgres (handling `channel_post` updates from the Telegram Bot API)
- [x] Backfill of channel history from a Telegram Desktop export
- [x] Draft generation in the channel's style (basic few-shot) — currently limited to a handful of hardcoded examples, style accuracy is being improved
- [x] Moderation and publishing via inline buttons in private messages
- [ ] Vector search for relevant few-shot examples (pgvector) — next up, to scale beyond a handful of hardcoded posts
- [ ] Scheduled automatic generation (cron) — currently triggered manually on purpose, while the pipeline is still being tuned

## Stack

- Go
- PostgreSQL
- Telegram Bot API (`go-telegram-bot-api`)
- Google Gemini API — draft generation
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
│   ├── collector/     # entry point: real-time post collection from Telegram
│   ├── generator/     # entry point: manual generation trigger, for testing
│   └── data/          # exported channel history (result.json)
├── internal/
│   ├── backfill/      # historical import logic
│   ├── generator/     # few-shot draft generation logic
│   └── store/         # Postgres access layer
├── docker-compose.yml
├── architecture.md
└── README.md
```

> Known cleanup item: some entry points and internal logic aren't cleanly split yet (e.g. `cmd/collector` currently holds logic that belongs in `internal/collector`, and `internal/backfill` currently holds an entry point that belongs in `cmd/backfill`). Left as-is for now, will be straightened out as the project stabilizes.

Structure and this README will be updated as the vector search and scheduler get added.