package store

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type Store struct {
	db *sql.DB
}

const MAXCONNS = 5
const CONLIFETIME time.Duration = 5 * time.Minute
const ConIdleTime time.Duration = 10 * time.Second

func New(connString string) (*Store, error) {
	db, err := sql.Open("pgx", connString)

	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(MAXCONNS)
	db.SetMaxIdleConns(MAXCONNS)
	db.SetConnMaxLifetime(CONLIFETIME)
	db.SetConnMaxIdleTime(ConIdleTime)

	return &Store{db: db}, nil
}

func (s *Store) SavePost(ctx context.Context, tgMessageID int64, text string, postedAt time.Time) error {
	_, err := s.db.ExecContext(ctx, `
	INSERT INTO posts (tg_message_id, text, posted_at)
	VALUES ($1, $2, $3)
	ON CONFLICT (tg_message_id) DO NOTHING`, tgMessageID, text, postedAt)

	return err
}

func (s *Store) GetLastPost(ctx context.Context, limit int) ([]string, error) {
	query := "SELECT tg_message_id, text, posted_at FROM posts ORDER BY tg_message_id ASC LIMIT $1;"
	rows, err := s.db.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("Error or timeout: %v", err)
	}

	defer rows.Close()

	var allPostsText []string
	for rows.Next() {
		var id int64
		var text string
		var postedAt time.Time

		err := rows.Scan(&id, &text, &postedAt)
		if err != nil {
			return nil, fmt.Errorf("Cannot read string: %v", err)
		}
		allPostsText = append(allPostsText, text)

	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("Error with strings:%v", err)
	}
	return allPostsText, err
}

func (s *Store) SaveDraft(ctx context.Context, text string) (int, error) {
	var id int
	query := "INSERT INTO drafts (text) VALUES ($1) RETURNING id;"

	row := s.db.QueryRowContext(ctx, query, text)
	err := row.Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("Cannot scan row id: %v", err)
	}

	return id, nil
}

func (s *Store) GetDraftByID(ctx context.Context, id int) (string, error) {
	var text string
	query := "SELECT text FROM drafts WHERE id = $1"
	row := s.db.QueryRowContext(ctx, query, id)
	err := row.Scan(&text)
	if err != nil {
		return "", fmt.Errorf("Cannot scan text from db: %v", err)
	}
	return text, nil
}

func (s *Store) ChangeStatus(ctx context.Context, id int) error {
	query := "UPDATE drafts SET status = 'published' WHERE id = $1"
	_, err := s.db.ExecContext(ctx, query)
	if err != nil {
		return fmt.Errorf("Cannot update status: %v", err)
	}
	return nil

}
