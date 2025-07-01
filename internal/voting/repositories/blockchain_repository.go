package repositories

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"voting-blockchain/internal/voting/models"
)

// BlockchainRepository — интерфейс для работы с блоками.
type BlockchainRepository interface {
	AddBlock(ctx context.Context, block *models.Block) error
	GetLastBlock(ctx context.Context, electionID int) (*models.Block, error)
	GetAllBlocks(ctx context.Context, electionID int) ([]*models.Block, error)
}

// BlockchainPostgres — реализация BlockchainRepository через PostgreSQL.
type BlockchainPostgres struct {
	DB *pgxpool.Pool
}

// Конструктор
func NewBlockchainPostgres(db *pgxpool.Pool) *BlockchainPostgres {
	return &BlockchainPostgres{DB: db}
}

// AddBlock — сохраняет новый блок в таблицу blockchain.
func (r *BlockchainPostgres) AddBlock(ctx context.Context, block *models.Block) error {
	query := `
		INSERT INTO blockchain (vote_hash, previous_hash, current_hash, election_id)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at
	`
	return r.DB.QueryRow(ctx, query,
		block.VoteHash,
		block.PrevHash,
		block.Hash,
		block.ElectionID,
	).Scan(&block.Index, &block.Timestamp)
}

// GetLastBlock — получает последний блок по голосованию.
func (r *BlockchainPostgres) GetLastBlock(ctx context.Context, electionID int) (*models.Block, error) {
	query := `
		SELECT id, created_at, vote_hash, previous_hash, current_hash, election_id
		FROM blockchain
		WHERE election_id = $1
		ORDER BY id DESC
		LIMIT 1
	`

	var b models.Block
	err := r.DB.QueryRow(ctx, query, electionID).Scan(
		&b.Index,
		&b.Timestamp,
		&b.VoteHash,
		&b.PrevHash,
		&b.Hash,
		&b.ElectionID,
	)
	if err != nil {
		return nil, err
	}
	return &b, nil
}

// GetAllBlocks — возвращает полную цепочку блоков для голосования.
func (r *BlockchainPostgres) GetAllBlocks(ctx context.Context, electionID int) ([]*models.Block, error) {
	query := `
		SELECT id, created_at, vote_hash, previous_hash, current_hash, election_id
		FROM blockchain
		WHERE election_id = $1
		ORDER BY id ASC
	`

	rows, err := r.DB.Query(ctx, query, electionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var blocks []*models.Block
	for rows.Next() {
		var b models.Block
		if err := rows.Scan(
			&b.Index,
			&b.Timestamp,
			&b.VoteHash,
			&b.PrevHash,
			&b.Hash,
			&b.ElectionID,
		); err != nil {
			return nil, err
		}
		blocks = append(blocks, &b)
	}
	return blocks, nil
}