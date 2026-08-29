package store

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/oklog/ulid/v2"
	"golang.org/x/crypto/bcrypt"
)

const bcryptCost = 10

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	if err != nil {
		return "", fmt.Errorf("hash password: %w", err)
	}
	return string(bytes), nil
}

func CheckPassword(password, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

// DefaultRoles and their permissions (mirrors middleware.RolePermissions).
var DefaultRoles = map[string][]string{
	"Super Admin":     {"app:admin", "gateway:write", "policy:admin", "event:read", "audit:read", "user:admin", "system:admin", "cert:write", "ip:write", "ratelimit:write"},
	"Platform Admin":  {"app:admin", "gateway:write", "policy:admin", "event:read", "audit:read", "user:read", "cert:write", "ip:write", "ratelimit:write"},
	"Security Admin":  {"app:write", "policy:write", "event:read", "audit:read", "ip:write", "ratelimit:write", "cert:read"},
	"App Owner":       {"app:read", "app:write", "gateway:read", "policy:read", "event:read"},
	"SOC Analyst":     {"event:read", "app:read", "gateway:read", "policy:read"},
	"Auditor":         {"audit:read", "event:read", "app:read", "gateway:read", "policy:read"},
	"Read Only":       {"app:read", "gateway:read", "policy:read", "event:read"},
}

// SeedRoles inserts the default roles and permissions if they don't exist.
func (s *Store) SeedRoles(ctx context.Context) error {
	for name, perms := range DefaultRoles {
		_, err := s.Pool.Exec(ctx,
			`INSERT INTO roles (id, name, permissions) VALUES ($1, $2, $3) ON CONFLICT (name) DO NOTHING`,
			ulid.Make().String(), name, perms)
		if err != nil {
			return fmt.Errorf("seed role %q: %w", name, err)
		}
	}
	return nil
}

// LookupUserRole returns the user's role name. Uses user_roles (direct mapping)
// with fallback to the old group-based lookup for backward compatibility.
func (s *Store) LookupUserRole(ctx context.Context, userID string) (string, error) {
	var role string
	err := s.Pool.QueryRow(ctx,
		`SELECT r.name FROM user_roles ur
		 JOIN roles r ON r.id = ur.role_id
		 WHERE ur.user_id = $1
		 LIMIT 1`, userID).Scan(&role)
	if err == nil {
		return role, nil
	}
	// Fallback: group-based lookup (legacy, for existing data).
	err = s.Pool.QueryRow(ctx,
		`SELECT r.name FROM roles r
		 JOIN user_group_memberships ugm ON ugm.group_id = r.id
		 WHERE ugm.user_id = $1
		 LIMIT 1`, userID).Scan(&role)
	if err == nil {
		return role, nil
	}
	return "Read Only", nil
}

// AssignRole assigns a role to a user.
func (s *Store) AssignRole(ctx context.Context, userID, roleID string) error {
	_, err := s.Pool.Exec(ctx,
		`INSERT INTO user_roles (user_id, role_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`,
		userID, roleID)
	return err
}

// RoleID returns the role id for a role name. Errors if not seeded.
func (s *Store) RoleID(ctx context.Context, name string) (string, error) {
	var roleID string
	err := s.Pool.QueryRow(ctx, `SELECT id FROM roles WHERE name = $1`, name).Scan(&roleID)
	return roleID, err
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

	// Seed roles and assign Super Admin to the initial user.
	if err := s.SeedRoles(ctx); err != nil {
		return fmt.Errorf("seed roles: %w", err)
	}
	var roleID string
	if err := s.Pool.QueryRow(ctx, `SELECT id FROM roles WHERE name = 'Super Admin'`).Scan(&roleID); err == nil {
		_ = s.AssignRole(ctx, userID, roleID)
	}
	return nil
}