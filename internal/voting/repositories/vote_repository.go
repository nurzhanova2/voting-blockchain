package repositories

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"voting-blockchain/internal/voting/models"
)

type VoteRepository interface {
	Create(ctx context.Context, v *models.Vote) error
	HasVoted(ctx context.Context, userID, electionID int) (bool, error)
	GetByElectionID(ctx context.Context, electionID int) ([]*models.Vote, error)
}

type VotePostgres struct {
	DB *pgxpool.Pool
}

func NewVotePostgres(db *pgxpool.Pool) *VotePostgres {
	return &VotePostgres{DB: db}
}

// Create — сохраняет голос в таблицу votes
func (r *VotePostgres) Create(ctx context.Context, v *models.Vote) error {
	query := `
		INSERT INTO votes (user_id, election_id, vote_hash)
		VALUES ($1, $2, $3)
		RETURNING id, created_at
	`
	return r.DB.QueryRow(ctx, query, v.UserID, v.ElectionID, v.VoteHash).Scan(&v.ID, &v.CreatedAt)
}

// HasVoted — проверяет, голосовал ли уже пользователь
func (r *VotePostgres) HasVoted(ctx context.Context, userID, electionID int) (bool, error) {
	query := `SELECT COUNT(*) FROM votes WHERE user_id = $1 AND election_id = $2`
	var count int
	err := r.DB.QueryRow(ctx, query, userID, electionID).Scan(&count)
	return count > 0, err
}

// GetByElectionID — получает все голоса по ID выборов
func (r *VotePostgres) GetByElectionID(ctx context.Context, electionID int) ([]*models.Vote, error) {
	query := `SELECT id, user_id, election_id, vote_hash, created_at FROM votes WHERE election_id = $1 ORDER BY created_at`
	rows, err := r.DB.Query(ctx, query, electionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var votes []*models.Vote
	for rows.Next() {
		var v models.Vote
		if err := rows.Scan(&v.ID, &v.UserID, &v.ElectionID, &v.VoteHash, &v.CreatedAt); err != nil {
			return nil, err
		}
		votes = append(votes, &v)
	}
	return votes, nil
}
