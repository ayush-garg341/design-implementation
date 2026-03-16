package store

import (
	"context"
	"database/sql"
)

type ClickStore struct {
	db *sql.DB
}

type ClickStats struct {
	ShortCode  string
	ClickCount uint64
	LongUrl    string
	UserId     string
}

func NewClickStore(db *sql.DB) *ClickStore {
	return &ClickStore{
		db: db,
	}
}

func (s *ClickStore) RecordClick(ctx context.Context, code string) error {

	_, err := s.db.ExecContext(
		ctx,
		`UPDATE links set click_count = click_count+1 where short_code=$1`,
		code,
	)

	return err
}

func (s *ClickStore) GetClickStats(ctx context.Context, code string) (*ClickStats, error) {
	var cl ClickStats
	err := s.db.QueryRowContext(ctx, `SELECT short_code, click_count, long_url, user_id FROM links where short_code=$1`, code).Scan(&cl.ShortCode, &cl.ClickCount, &cl.LongUrl, &cl.UserId)
	return &cl, err
}
