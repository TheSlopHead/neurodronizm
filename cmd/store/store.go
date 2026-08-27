package store

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type Store struct {
	db *sql.DB
}

func New(connString string) (*Store, error) {
	db, err := sql.Open("pgx", connString)

	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &Store{db: db}, nil
}

func (s *Store) SavePost(ctx context.Context, tgMessageID int64, text string, postedAt time.Time) error {
	_, err := s.db.ExecContext(ctx, `
	INSERT INTO posts (tg_message_id, text, posted_at)
	VALUES ($1, $2, $3)
	ON CONFLICT (tg_message_id) DO NOTHING`, tgMessageID, text, postedAt)

	return err
}

//constring = postgres://dronism:change_me@localhost:5432/dronism?sslmode=disable
