package models

import "time"

// Election — структура голосования (создаётся админом)
type Election struct {
	ID          int       `db:"id"`
	Title       string    `db:"title"`
	Description string    `db:"description"`
	CreatedBy   int       `db:"created_by"` // ID администратора
	CreatedAt   time.Time `db:"created_at"`
	IsActive    bool      `db:"is_active"`  // Можно ли ещё голосовать
}
