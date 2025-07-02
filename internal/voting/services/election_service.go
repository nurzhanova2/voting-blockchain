package services

import (
	"context"

	"voting-blockchain/internal/voting/models"
	"voting-blockchain/internal/voting/repositories"
)

type ElectionService interface {
	Create(ctx context.Context, e *models.Election) error
	GetByID(ctx context.Context, id int) (*models.Election, error)
	List(ctx context.Context) ([]*models.Election, error)
	Update(ctx context.Context, e *models.Election) error
	Delete(ctx context.Context, id int) error
	CreateChoices(ctx context.Context, electionID int, choices []string) error
}

type electionService struct {
	electionRepo repositories.ElectionRepository
	choiceRepo   repositories.ChoiceRepository
}

func NewElectionService(
	electionRepo repositories.ElectionRepository,
	choiceRepo repositories.ChoiceRepository,
) ElectionService {
	return &electionService{
		electionRepo: electionRepo,
		choiceRepo:   choiceRepo,
	}
}

func (s *electionService) Create(ctx context.Context, e *models.Election) error {
	return s.electionRepo.Create(ctx, e)
}

func (s *electionService) GetByID(ctx context.Context, id int) (*models.Election, error) {
	return s.electionRepo.GetByID(ctx, id)
}

func (s *electionService) List(ctx context.Context) ([]*models.Election, error) {
	return s.electionRepo.List(ctx)
}

func (s *electionService) Update(ctx context.Context, e *models.Election) error {
	return s.electionRepo.Update(ctx, e)
}

func (s *electionService) Delete(ctx context.Context, id int) error {
	return s.electionRepo.Delete(ctx, id)
}

func (s *electionService) CreateChoices(ctx context.Context, electionID int, choices []string) error {
	return s.choiceRepo.CreateChoices(ctx, electionID, choices)
}
