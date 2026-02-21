package auth

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrEmailAlreadyExists   = errors.New("email already exists")
	ErrInvalidCredentials   = errors.New("invalid credentials")
	ErrRefreshTokenNotFound = errors.New("refresh token not found")
	ErrRefreshTokenRevoked  = errors.New("refresh token revoked")
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

type User struct {
	ID           string
	Email        string
	PasswordHash string
	Mobile       string
	Role         string
	CreatedAt    time.Time
}

type RefreshTokenRow struct {
	ID        string
	UserID    string
	RevokedAt *time.Time
	ExpiresAt time.Time
	CreatedAt time.Time
}

func (r *Repository) CreateUser(ctx context.Context, email, passwordHash, mobile, role string) (User, error) {
	const q = `
		INSERT INTO users (email, password_hash, mobile, role)
		VALUES ($1, $2, $3, $4)
		RETURNING id, email, password_hash, mobile, role, created_at
	`

	var u User
	err := r.db.QueryRow(ctx, q, email, passwordHash, mobile, role).
		Scan(&u.ID, &u.Email, &u.PasswordHash, &u.Mobile, &u.Role, &u.CreatedAt)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
			// Check the error detail to see which column caused the violation
			if strings.Contains(pgErr.Detail, "email") {
				return User{}, ErrEmailAlreadyExists
			}
			if strings.Contains(pgErr.Detail, "mobile") {
				return User{}, errors.New("mobile already exists")
			}
			return User{}, errors.New("duplicate value")
		}

		return User{}, err
	}

	return u, nil
}

func (r *Repository) GetUserByEmail(ctx context.Context, email string) (User, error) {
	const q = `
		SELECT id, email, password_hash, mobile, role, created_at
		FROM users
		WHERE email = $1
	`
	var u User
	err := r.db.QueryRow(ctx, q, email).
		Scan(&u.ID, &u.Email, &u.PasswordHash, &u.Mobile, &u.Role, &u.CreatedAt)
	if err != nil {
		return User{}, ErrInvalidCredentials
	}
	return u, nil
}

func (r *Repository) InsertRefreshToken(ctx context.Context, tokenID, userID string, expiresAt time.Time) error {
	const q = `
		INSERT INTO refresh_tokens (id, user_id, expires_at)
		VALUES ($1, $2, $3)
	`
	_, err := r.db.Exec(ctx, q, tokenID, userID, expiresAt)
	return err
}

func (r *Repository) GetRefreshToken(ctx context.Context, tokenID string) (RefreshTokenRow, error) {
	const q = `
		SELECT id, user_id, revoked_at, expires_at, created_at
		FROM refresh_tokens
		WHERE id = $1
	`
	var rt RefreshTokenRow
	err := r.db.QueryRow(ctx, q, tokenID).Scan(&rt.ID, &rt.UserID, &rt.RevokedAt, &rt.ExpiresAt, &rt.CreatedAt)
	if err != nil {
		return RefreshTokenRow{}, ErrRefreshTokenNotFound
	}
	return rt, nil
}

func (r *Repository) RevokeRefreshToken(ctx context.Context, tokenID string, revokedAt time.Time) error {
	const q = `
		UPDATE refresh_tokens
		SET revoked_at = $2
		WHERE id = $1 AND revoked_at IS NULL
	`
	cmd, err := r.db.Exec(ctx, q, tokenID, revokedAt)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return ErrRefreshTokenRevoked
	}
	return nil
}

func (r *Repository) GetUserByID(ctx context.Context, userID string) (User, error) {
	const q = `
		SELECT id, email, password_hash, mobile, role, created_at
		FROM users
		WHERE id = $1
	`
	var u User
	err := r.db.QueryRow(ctx, q, userID).
		Scan(&u.ID, &u.Email, &u.PasswordHash, &u.Mobile, &u.Role, &u.CreatedAt)
	if err != nil {
		return User{}, errors.New("user not found")
	}
	return u, nil
}
