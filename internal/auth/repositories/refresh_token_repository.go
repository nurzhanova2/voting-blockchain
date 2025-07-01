package repositories

import (
	"context"

	"voting-blockchain/db"
	"voting-blockchain/internal/auth/models"
)

type refreshTokenRepository struct{}

func NewRefreshTokenRepository() RefreshTokenRepository {
	return &refreshTokenRepository{}
}

// Сохраняет refresh токен
func (r *refreshTokenRepository) Save(ctx context.Context, t *models.RefreshToken) error {
	query := `
		INSERT INTO refresh_tokens (token, user_id, expires_at, revoked)
		VALUES ($1, $2, $3, $4)
	`
	_, err := db.DB.Exec(ctx, query, t.Token, t.UserID, t.ExpiresAt, t.Revoked)
	return err
}

// Находит refresh токен по строке
func (r *refreshTokenRepository) FindByToken(ctx context.Context, token string) (*models.RefreshToken, error) {
	query := `
		SELECT token, user_id, expires_at, revoked
		FROM refresh_tokens
		WHERE token = $1
	`
	row := db.DB.QueryRow(ctx, query, token)

	var t models.RefreshToken
	err := row.Scan(&t.Token, &t.UserID, &t.ExpiresAt, &t.Revoked)
	if err != nil {
		return nil, err
	}

	return &t, nil
}

// Помечает refresh токен как отозванный
func (r *refreshTokenRepository) Revoke(ctx context.Context, token string) error {
	query := `UPDATE refresh_tokens SET revoked = true WHERE token = $1`
	_, err := db.DB.Exec(ctx, query, token)
	return err
}
