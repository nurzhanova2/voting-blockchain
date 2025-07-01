package services

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"voting-blockchain/internal/voting/models"
	"voting-blockchain/internal/voting/repositories"
)

// BlockchainService описывает методы работы с блокчейном
type BlockchainService interface {
	AddBlock(ctx context.Context, electionID int, voteHash string) (*models.Block, error)
	GetChain(ctx context.Context, electionID int) ([]*models.Block, error)
}

type blockchainService struct {
	blockRepo repositories.BlockchainRepository
}

// NewBlockchainService — конструктор сервиса
func NewBlockchainService(blockRepo repositories.BlockchainRepository) BlockchainService {
	return &blockchainService{blockRepo: blockRepo}
}

// AddBlock — добавляет блок в блокчейн голосования
func (s *blockchainService) AddBlock(ctx context.Context, electionID int, voteHash string) (*models.Block, error) {
	// Получаем последний блок
	prevBlock, err := s.blockRepo.GetLastBlock(ctx, electionID)
	var prevHash string
	if err == nil && prevBlock != nil {
		prevHash = prevBlock.Hash
	}

	// Строим хеш нового блока
	timestamp := time.Now()
	raw := fmt.Sprintf("%s|%s|%d|%s", voteHash, prevHash, electionID, timestamp.UTC().String())
	hash := sha256.Sum256([]byte(raw))
	hashStr := hex.EncodeToString(hash[:])

	newBlock := &models.Block{
		Timestamp:  timestamp,
		VoteHash:   voteHash,
		PrevHash:   prevHash,
		Hash:       hashStr,
		ElectionID: electionID,
	}

	err = s.blockRepo.AddBlock(ctx, newBlock)
	if err != nil {
		return nil, err
	}
	return newBlock, nil
}

// GetChain — получить всю цепочку блоков для голосования
func (s *blockchainService) GetChain(ctx context.Context, electionID int) ([]*models.Block, error) {
	return s.blockRepo.GetAllBlocks(ctx, electionID)
}
