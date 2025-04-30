package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/KarmaBeLike/jwt-auth-service/internal/model"
)

type AuthRepository struct {
	db *sql.DB
}

func NewAuthRepository(db *sql.DB) *AuthRepository {
	return &AuthRepository{db: db}
}

func (r *AuthRepository) SaveRefreshToken(ctx context.Context, token model.RefreshToken) error {
	query := `
        INSERT INTO refresh_tokens (id, user_id, refresh_token_hash, access_token_id, client_ip, created_at, expires_at)
        VALUES ($1, $2, $3, $4, $5, $6,$7)
    `
	_, err := r.db.ExecContext(ctx, query,
		token.ID, token.UserID, token.TokenHash, token.AccessTokenID, token.IPAddress, token.CreatedAt, token.ExpiresAt,
	)
	return err
}

func (r *AuthRepository) GetRefreshToken(ctx context.Context, id string) (*model.RefreshToken, error) {
	query := `
        SELECT id, user_id,refresh_token_hash, access_token_id,client_ip, created_at, expires_at
        FROM refresh_tokens
        WHERE id = $1
    `
	row := r.db.QueryRowContext(ctx, query, id)

	var token model.RefreshToken
	if err := row.Scan(&token.ID, &token.UserID, &token.TokenHash, &token.AccessTokenID, &token.IPAddress, &token.CreatedAt, &token.ExpiresAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &token, nil
}

func (r *AuthRepository) DeleteRefreshToken(ctx context.Context, id string) error {
	query := `
        DELETE FROM refresh_tokens
        WHERE id = $1
    `
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
