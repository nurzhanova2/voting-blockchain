package repositories

import (
	"context"
	"voting-blockchain/internal/voting/models"
)

type ElectionRepository interface {
    Create(ctx context.Context, e *models.Election) error
    GetByID(ctx context.Context, id int) (*models.Election, error)
    List(ctx context.Context) ([]*models.Election, error)
    Update(ctx context.Context, e *models.Election) error
    Delete(ctx context.Context, id int) error
}
