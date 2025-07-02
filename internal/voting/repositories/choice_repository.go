package repositories

import (
    "context"
    "voting-blockchain/internal/voting/models"
)

type ChoiceRepository interface {
    CreateChoices(ctx context.Context, electionID int, choices []string) error
    GetChoices(ctx context.Context, electionID int) ([]*models.Choice, error)
}
