package repositories

import (
    "context"

    "github.com/jackc/pgx/v5/pgxpool"
    "voting-blockchain/internal/voting/models"
)

type ElectionPostgres struct {
    DB *pgxpool.Pool
}

func NewElectionPostgres(db *pgxpool.Pool) *ElectionPostgres {
    return &ElectionPostgres{DB: db}
}

func (r *ElectionPostgres) Create(ctx context.Context, e *models.Election) error {
    query := `
        INSERT INTO elections (title, description, created_by, is_active)
        VALUES ($1, $2, $3, $4)
        RETURNING id, created_at
    `
    return r.DB.QueryRow(ctx, query,
        e.Title,
        e.Description,
        e.CreatedBy,
        e.IsActive,
    ).Scan(&e.ID, &e.CreatedAt)
}

func (r *ElectionPostgres) GetByID(ctx context.Context, id int) (*models.Election, error) {
    query := `
        SELECT id, title, description, created_by, created_at, is_active
        FROM elections
        WHERE id = $1
    `
    var e models.Election
    err := r.DB.QueryRow(ctx, query, id).Scan(
        &e.ID,
        &e.Title,
        &e.Description,
        &e.CreatedBy,
        &e.CreatedAt,
        &e.IsActive,
    )
    if err != nil {
        return nil, err
    }
    return &e, nil
}

func (r *ElectionPostgres) List(ctx context.Context) ([]*models.Election, error) {
    query := `
        SELECT id, title, description, created_by, created_at, is_active
        FROM elections
        ORDER BY created_at DESC
    `
    rows, err := r.DB.Query(ctx, query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var elections []*models.Election
    for rows.Next() {
        var e models.Election
        if err := rows.Scan(
            &e.ID,
            &e.Title,
            &e.Description,
            &e.CreatedBy,
            &e.CreatedAt,
            &e.IsActive,
        ); err != nil {
            return nil, err
        }
        elections = append(elections, &e)
    }
    return elections, nil
}

func (r *ElectionPostgres) Update(ctx context.Context, e *models.Election) error {
    query := `
        UPDATE elections
        SET title = $1, description = $2, is_active = $3
        WHERE id = $4
    `
    _, err := r.DB.Exec(ctx, query,
        e.Title,
        e.Description,
        e.IsActive,
        e.ID,
    )
    return err
}

func (r *ElectionPostgres) Delete(ctx context.Context, id int) error {
    query := `DELETE FROM elections WHERE id = $1`
    _, err := r.DB.Exec(ctx, query, id)
    return err
}
