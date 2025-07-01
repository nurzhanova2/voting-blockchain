package repositories

import (
	"context"

	"voting-blockchain/internal/auth/models"
)

type RefreshTokenRepository interface {
	Save(ctx context.Context, token *models.RefreshToken) error
	FindByToken(ctx context.Context, token string) (*models.RefreshToken, error)
}
