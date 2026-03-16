package store

import (
	"context"
	"database/sql"
)

type Link struct {
	ID         string
	Shortcode  string
	Longurl    string
	ClickCount uint64
	CreatedBy  string
	CreatedAt  string
	ExpiresAt  string
}

type LinkStore interface {
	SaveShortLink(ctx context.Context, link Link) (Link, error)
	GetAllLinks(ctx context.Context, userId string) ([]Link, error)
	GetRedirectUrl(ctx context.Context, shortUrl string) (*string, error)
}

type PostgresLinkStore struct {
	db *sql.DB
}

func NewPostgresLinkStore(db *sql.DB) *PostgresLinkStore {
	return &PostgresLinkStore{db}
}

func (s *PostgresLinkStore) SaveShortLink(ctx context.Context, l Link) (Link, error) {
	err := s.db.QueryRowContext(ctx, `INSERT into links (short_code, long_url, user_id, click_count)
VALUES ($1, $2, $3, $4) RETURNING id, created_at`,
		l.Shortcode, l.Longurl, l.CreatedBy, l.ClickCount).Scan(&l.ID, &l.CreatedAt)
	return l, err
}

func (s *PostgresLinkStore) GetAllLinks(ctx context.Context, userID string) ([]Link, error) {
	var links []Link

	rows, err := s.db.QueryContext(ctx, `SELECT short_code, long_url from links where user_id=$1`, userID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var l Link
		if err := rows.Scan(&l.Shortcode, &l.Longurl); err != nil {
			return nil, err
		}

		links = append(links, l)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return links, nil

}

func (s *PostgresLinkStore) GetRedirectUrl(ctx context.Context, shortCode string) (*string, error) {
	var longurl string
	err := s.db.QueryRowContext(ctx, `SELECT long_url from links where short_code=$1`, shortCode).Scan(&longurl)
	if err != nil {
		return nil, err
	}

	return &longurl, nil
}
