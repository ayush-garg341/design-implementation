package store

import (
	"context"
	"database/sql"
	"errors"
	"github.com/lib/pq"
)

var (
	ErrUserNotFound  = errors.New("user not found")
	ErrDuplicateUser = errors.New("user already exists")
	ErrForeignKey    = errors.New("invalid reference")
)

type User struct {
	ID        string
	Email     string
	Name      string
	Password  string
	CreatedAt string
}

// UserStore Interface

type UserStore interface {
	Create(ctx context.Context, user User) (*User, error)
	GetByEmail(ctx context.Context, email string) (User, error)
	GetByID(ctx context.Context, id string) (User, error)
	Login(ctx context.Context, email string) (*User, error)
}

// PostgresUserStore
type PostgresUserStore struct {
	db *sql.DB
}

func NewPostgresUserStore(db *sql.DB) *PostgresUserStore {
	return &PostgresUserStore{db}
}

func (s *PostgresUserStore) Create(ctx context.Context, u User) (*User, error) {
	err := s.db.QueryRowContext(ctx, `INSERT into users (email, name, password)
                                    VALUES ($1, $2, $3) RETURNING id, created_at`,
		u.Email, u.Name, u.Password).Scan(&u.ID, &u.CreatedAt)

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code {
			case "23505":
				return nil, ErrDuplicateUser
			}
		}

	}
	return &u, nil
}

func (s *PostgresUserStore) GetByEmail(ctx context.Context, email string) (User, error) {
	var u User
	err := s.db.QueryRowContext(ctx, `SELECT id, email, name, password, created_at from users where email=$1`, email).Scan(&u.ID, &u.Email, &u.Name, &u.Password, &u.CreatedAt)
	return u, err
}

func (s *PostgresUserStore) GetByID(ctx context.Context, id string) (User, error) {
	var u User
	err := s.db.QueryRowContext(ctx, `SELECT id, email, name, password, created_at from users where id=$1`, id).Scan(&u.ID, &u.Email, &u.Name, &u.Password, &u.CreatedAt)
	return u, err
}

func (s *PostgresUserStore) Login(ctx context.Context, email string) (*User, error) {
	var u User
	err := s.db.QueryRowContext(ctx, `SELECT id, email, name, password, created_at from users where email=$1`, email).Scan(&u.ID, &u.Email, &u.Name, &u.Password, &u.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &u, nil
}
