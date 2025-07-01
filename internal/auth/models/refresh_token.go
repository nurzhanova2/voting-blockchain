package models

import "time"

type RefreshToken struct {
	Token     string
	UserID    int
	ExpiresAt time.Time
	Revoked   bool
}
