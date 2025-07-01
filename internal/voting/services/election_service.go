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
}

type electionService struct {
    repo repositories.ElectionRepository
}

func NewElectionService(repo repositories.ElectionRepository) ElectionService {
    return &electionService{repo: repo}
}

func (s *electionService) Create(ctx context.Context, e *models.Election) error {
    return s.repo.Create(ctx, e)
}

func (s *electionService) GetByID(ctx context.Context, id int) (*models.Election, error) {
    return s.repo.GetByID(ctx, id)
}

func (s *electionService) List(ctx context.Context) ([]*models.Election, error) {
    return s.repo.List(ctx)
}

func (s *electionService) Update(ctx context.Context, e *models.Election) error {
    return s.repo.Update(ctx, e)
}

func (s *electionService) Delete(ctx context.Context, id int) error {
    return s.repo.Delete(ctx, id)
}
