package store

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/oklog/ulid/v2"
	"golang.org/x/crypto/bcrypt"
)

// cost for bcrypt hashing
const bcryptCost = 10

// HashPassword returns a bcrypt hash of the password.
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	if err != nil {
		return "", fmt.Errorf("hash password: %w", err)
	}
	return string(bytes), nil
}

// CheckPassword verifies a password against a bcrypt hash.
func CheckPassword(password, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

// Store wraps the PostgreSQL connection pool.
type Store struct {
	Pool *pgxpool.Pool
}

// NewID generates a new ULID public identifier.
func (s *Store) NewID() (string, error) {
	return ulid.Make().String(), nil
}

// Open creates a connection pool and verifies connectivity.
func Open(ctx context.Context, databaseURL string) (*Store, error) {
	cfg, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		return nil, fmt.Errorf("parse pool config: %w", err)
	}
	cfg.MaxConns = 20
	cfg.MinConns = 2
	cfg.MaxConnLifetime = 30 * time.Minute

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("create pool: %w", err)
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping database: %w", err)
	}
	return &Store{Pool: pool}, nil
}

// Close shuts down the pool.
func (s *Store) Close() {
	s.Pool.Close()
}

// EnsureInitialAdmin creates the initial admin user and organization
// if they do not exist. Called once at startup.
func (s *Store) EnsureInitialAdmin(ctx context.Context, email, password string) error {
	var count int
	if err := s.Pool.QueryRow(ctx, `SELECT COUNT(*) FROM users WHERE status = 'active'`).Scan(&count); err != nil {
		return fmt.Errorf("check users: %w", err)
	}
	if count > 0 {
		return nil
	}

	orgID := "01ARZ3NDEKTSV4RRFFQ69G5FAV"
	if _, err := s.Pool.Exec(ctx, `INSERT INTO organizations (id, name) VALUES ($1, $2) ON CONFLICT DO NOTHING`, orgID, "default"); err != nil {
		return fmt.Errorf("create org: %w", err)
	}

	userID := "01ARZ3NDEKTSV4RRFFQ69G5FAW"
	hash, err := HashPassword(password)
	if err != nil {
		return fmt.Errorf("hash password: %w", err)
	}
	if _, err := s.Pool.Exec(ctx, `INSERT INTO users (id, organization_id, email, password_hash, language) VALUES ($1, $2, $3, $4, 'en') ON CONFLICT DO NOTHING`,
		userID, orgID, email, hash); err != nil {
		return fmt.Errorf("create initial admin: %w", err)
	}
	return nil
}