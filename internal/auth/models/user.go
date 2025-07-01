package models

import "time"

// User — структура пользователя системы
type User struct {
	ID           int       `db:"id"`
	Email        string    `db:"email"`
	PasswordHash string    `db:"password_hash"`
	Role         string    `db:"role"` // "admin" или "user"
	CreatedAt    time.Time `db:"created_at"`
}
