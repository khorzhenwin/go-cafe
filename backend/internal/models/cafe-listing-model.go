package models

import "time"

type CafeListing struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	UserID          uint      `gorm:"not null;index" json:"user_id"`
	User            *User     `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Name            string    `gorm:"not null" json:"name"`
	Address         string    `json:"address"`
	City            string    `json:"city,omitempty"`
	Neighborhood    string    `json:"neighborhood,omitempty"`
	Description     string    `json:"description,omitempty"`
	ImageURL        string    `json:"image_url,omitempty"`
	Latitude        *float64  `json:"latitude,omitempty"`
	Longitude       *float64  `json:"longitude,omitempty"`
	SourceProvider  string    `gorm:"index" json:"source_provider,omitempty"`
	ExternalPlaceID string    `gorm:"index" json:"external_place_id,omitempty"`
	VisitStatus     string    `gorm:"not null;default:to_visit;index" json:"visit_status"`
	SourceCafeID    *uint     `gorm:"index" json:"source_cafe_id,omitempty"`
	AvgRating       float64   `gorm:"-" json:"avg_rating"`
	ReviewCount     int64     `gorm:"-" json:"review_count"`
}
