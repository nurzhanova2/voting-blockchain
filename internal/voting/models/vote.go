package models

import "time"

// Vote представляет голос пользователя в конкретных выборах.
type Vote struct {
	ID         int       // Уникальный ID голоса
	UserID     int       // ID пользователя, который проголосовал
	ElectionID int       // ID выборов, в которых проголосовал
	VoteHash   string    // Хэш голоса (содержимое + подпись)
	CreatedAt  time.Time // Время создания голоса
}
