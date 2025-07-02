package services

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"voting-blockchain/internal/voting/models"
	"voting-blockchain/internal/voting/repositories"
)

type VoteService interface {
	CastVote(ctx context.Context, userID, electionID int, choice string) error
	GetBlockchain(ctx context.Context, electionID int) ([]*models.Block, error)
	GetResults(ctx context.Context, electionID int) (map[string]int, error)
	GetChoices(ctx context.Context, electionID int) ([]*models.Choice, error)
}

type voteService struct {
	voteRepo  repositories.VoteRepository
	blockRepo repositories.BlockchainRepository
}

func NewVoteService(
	voteRepo repositories.VoteRepository,
	blockRepo repositories.BlockchainRepository,
) VoteService {
	return &voteService{
		voteRepo:  voteRepo,
		blockRepo: blockRepo,
	}
}

func (s *voteService) CastVote(ctx context.Context, userID, electionID int, choice string) error {
	exists, err := s.voteRepo.HasVoted(ctx, userID, electionID)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("пользователь уже голосовал в этом голосовании")
	}

	hash := generateVoteHash(userID, electionID, choice)

	vote := &models.Vote{
		UserID:     userID,
		ElectionID: electionID,
		Choice:     choice,
		VoteHash:   hash,
	}

	if err := s.voteRepo.Create(ctx, vote); err != nil {
		return err
	}

	lastBlock, err := s.blockRepo.GetLastBlock(ctx, electionID)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return err
	}

	newBlock := &models.Block{
		Timestamp:  time.Now(),
		VoteHash:   vote.VoteHash,
		ElectionID: electionID,
	}

	if lastBlock != nil {
		newBlock.PrevHash = lastBlock.Hash
	}

	newBlock.Hash = generateBlockHash(newBlock)

	return s.blockRepo.AddBlock(ctx, newBlock)
}

func (s *voteService) GetBlockchain(ctx context.Context, electionID int) ([]*models.Block, error) {
	return s.blockRepo.GetAllBlocks(ctx, electionID)
}

// Подсчет количества голосов по каждому варианту
func (s *voteService) GetResults(ctx context.Context, electionID int) (map[string]int, error) {
	blocks, err := s.blockRepo.GetAllBlocks(ctx, electionID)
	if err != nil {
		return nil, err
	}

	results := make(map[string]int)

	for _, block := range blocks {
		vote, err := s.voteRepo.GetByHash(ctx, block.VoteHash)
		if err != nil {
			return nil, err
		}
		results[vote.Choice]++
	}

	return results, nil
}

// Возврат списка уникальных вариантов (Choices)
func (s *voteService) GetChoices(ctx context.Context, electionID int) ([]*models.Choice, error) {
	return s.voteRepo.GetResults(ctx, electionID)
}

func generateVoteHash(userID, electionID int, choice string) string {
	raw := fmt.Sprintf("%d|%d|%s", userID, electionID, choice)
	hash := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(hash[:])
}

func generateBlockHash(b *models.Block) string {
	data := []byte(
		b.Timestamp.String() + b.VoteHash + b.PrevHash + fmt.Sprintf("%d", b.ElectionID),
	)
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])
}
