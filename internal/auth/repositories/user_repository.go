package repositories

import (
	"context"
	"errors"
	"time"

	"voting-blockchain/db"
	"voting-blockchain/internal/auth/models"
)

type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	FindByEmail(ctx context.Context, email string) (*models.User, error)
	FindByID(ctx context.Context, id int) (*models.User, error)
}

type userRepository struct{}

func NewUserRepository() UserRepository {
	return &userRepository{}
}

func (r *userRepository) FindByID(ctx context.Context, id int) (*models.User, error) {
	query := `
		SELECT id, email, password_hash, created_at
		FROM users
		WHERE id = $1
	`
	row := db.DB.QueryRow(ctx, query, id)

	var user models.User
	err := row.Scan(&user.ID, &user.Email, &user.PasswordHash, &user.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &user, nil
}


// Создание нового пользователя
func (r *userRepository) Create(ctx context.Context, user *models.User) error {
	query := `
		INSERT INTO users (email, password_hash, created_at)
		VALUES ($1, $2, $3)
		RETURNING id
	`

	err := db.DB.QueryRow(ctx, query, user.Email, user.PasswordHash, time.Now()).Scan(&user.ID)
	return err
}

// Поиск по email (для входа)
func (r *userRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `
		SELECT id, email, password_hash, created_at, role
		FROM users
		WHERE email = $1
	`

	row := db.DB.QueryRow(ctx, query, email)

	var user models.User
	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.Role,
	)
	if err != nil {
		return nil, errors.New("пользователь не найден")
	}

	return &user, nil
}

