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
- [x] Draft generation using few-shot examples retrieved via vector search (pgvector)
- [x] Moderation and publishing via inline buttons in private messages
- [ ] Scheduled automatic generation (cron) — currently triggered manually on purpose, while the pipeline is still being tuned
- [ ] Dockerize the Go services themselves (only Postgres runs in Docker so far) — deferred for now

## Stack

- Go
- PostgreSQL
- Telegram Bot API (`go-telegram-bot-api`)
- Google Gemini API — draft generation and embeddings
- pgvector — similarity search over past posts

## Running locally

1. Clone the repository.

2. Start Postgres:
   ```bash
   docker compose up -d
   ```

3. Apply the schema (one-time, via `psql`):
   ```sql
   -- 1. Main table with Telegram post history
   CREATE TABLE IF NOT EXISTS posts (
       tg_message_id BIGINT PRIMARY KEY,
       text TEXT NOT NULL,
       posted_at TIMESTAMP WITH TIME ZONE NOT NULL
   );

   -- 2. Draft table (the moderation queue)
   CREATE TABLE IF NOT EXISTS drafts (
       id SERIAL PRIMARY KEY,
       text TEXT NOT NULL,
       status VARCHAR(20) DEFAULT 'draft',
       created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
   );

   -- 3. pgvector extension and embeddings table
   CREATE EXTENSION IF NOT EXISTS vector;
   CREATE TABLE IF NOT EXISTS post_embeddings (
       post_id BIGINT PRIMARY KEY REFERENCES posts(tg_message_id) ON DELETE CASCADE,
       embedding vector(768) NOT NULL
   );
   ```

4. Set environment variables:
   ```bash
   export TELEGRAM_BOT_TOKEN=...
   export DATABASE_URL=postgres://dronism:change_me@localhost:5432/dronism?sslmode=disable
   ```

5. Run the bot (real-time collection + moderation):
   ```bash
   go run .
   ```

6. One-off commands, run as needed:
   ```bash
   go run ./cmd/backfill    # import historical posts from a Telegram Desktop export
   go run ./cmd/embedder    # embed any posts that don't have a vector yet
   go run ./cmd/generator   # manually trigger a draft generation, for testing
   ```

## Repository structure

```
neurodronism/
├── cmd/
│   ├── backfill/      # entry point: one-time historical import
│   ├── embedder/      # entry point: embeds posts that don't have a vector yet
│   └── generator/     # entry point: manual generation trigger, for testing
├── data/              # exported channel history (result.json)
├── internal/
│   ├── collector/     # real-time post collection logic
│   ├── generator/     # few-shot draft generation logic
│   └── store/         # Postgres access layer
├── main.go            # entry point: the running bot (real-time collection + moderation)
├── docker-compose.yml
├── architecture.md
└── README.md
```

Structure and this README will be updated as the scheduler gets added and the services get containerized.