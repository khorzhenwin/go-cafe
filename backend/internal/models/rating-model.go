package models

import "time"

type Rating struct {
	ID            uint        `gorm:"primaryKey" json:"id"`
	CreatedAt     time.Time   `json:"created_at"`
	UpdatedAt     time.Time   `json:"updated_at"`
	UserID        uint        `gorm:"not null;index" json:"user_id"`
	User          *User       `gorm:"foreignKey:UserID" json:"user,omitempty"`
	CafeListingID uint        `gorm:"not null;index" json:"cafe_listing_id"`
	CafeListing   *CafeListing `gorm:"foreignKey:CafeListingID" json:"cafe_listing,omitempty"`
	VisitedAt     time.Time   `gorm:"not null" json:"visited_at"`
	Rating        int         `gorm:"not null" json:"rating"` // e.g. 1-5
	Review        string      `json:"review,omitempty"`
}
