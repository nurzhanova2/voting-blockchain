package repositories

import (
    "context"

    "github.com/jackc/pgx/v5/pgxpool"
    "voting-blockchain/internal/voting/models"
)

type ChoicePostgres struct {
    DB *pgxpool.Pool
}

func NewChoicePostgres(db *pgxpool.Pool) *ChoicePostgres {
    return &ChoicePostgres{DB: db}
}

func (r *ChoicePostgres) CreateChoices(ctx context.Context, electionID int, choices []string) error {
    for _, text := range choices {
        _, err := r.DB.Exec(ctx,
            `INSERT INTO choices (election_id, text) VALUES ($1, $2)`,
            electionID, text,
        )
        if err != nil {
            return err
        }
    }
    return nil
}

func (r *ChoicePostgres) GetChoices(ctx context.Context, electionID int) ([]*models.Choice, error) {
    rows, err := r.DB.Query(ctx,
        `SELECT id, election_id, text FROM choices WHERE election_id = $1 ORDER BY id`, electionID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var res []*models.Choice
    for rows.Next() {
        var c models.Choice
        if err := rows.Scan(&c.ID, &c.ElectionID, &c.Text); err != nil {
            return nil, err
        }
        res = append(res, &c)
    }
    return res, nil
}
