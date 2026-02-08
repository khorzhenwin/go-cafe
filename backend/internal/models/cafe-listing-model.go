package models

import "time"

type CafeListing struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	UserID      uint       `gorm:"not null;index" json:"user_id"`
	User        *User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Name        string     `gorm:"not null" json:"name"`
	Address     string     `json:"address"`
	Description string     `json:"description,omitempty"`
}
