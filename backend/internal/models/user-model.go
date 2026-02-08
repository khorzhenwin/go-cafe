package models

import "time"

type User struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Email        string    `gorm:"uniqueIndex;not null" json:"email"`
	Name         string    `json:"name"`
	PasswordHash string    `json:"-"` // empty for legacy users; required for login
}
