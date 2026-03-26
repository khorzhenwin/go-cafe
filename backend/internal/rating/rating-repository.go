package rating

import (
	"errors"

	"github.com/khorzhenwin/go-cafe/backend/internal/models"
	"gorm.io/gorm"
)

type Storage interface {
	Create(rating *models.Rating) error
	GetByID(id uint) (*models.Rating, error)
	GetByCafeListingID(cafeListingID uint) ([]models.Rating, error)
	GetByExternalPlaceID(externalPlaceID string) ([]models.Rating, error)
	GetByUserID(userID uint) ([]models.Rating, error)
	FindByUserAndCafe(userID uint, cafeListingID uint) (*models.Rating, error)
	Update(id uint, updated models.Rating) error
	Delete(id uint) error
}

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(rt *models.Rating) error {
	return r.db.Create(rt).Error
}

func (r *Repository) GetByID(id uint) (*models.Rating, error) {
	var rating models.Rating
	err := r.db.Preload("User").Preload("CafeListing").First(&rating, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &rating, err
}

func (r *Repository) GetByCafeListingID(cafeListingID uint) ([]models.Rating, error) {
	var cafe models.CafeListing
	if err := r.db.First(&cafe, cafeListingID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []models.Rating{}, nil
		}
		return nil, err
	}

	rootID := cafe.ID
	if cafe.SourceCafeID != nil {
		rootID = *cafe.SourceCafeID
	}

	if cafe.ExternalPlaceID != "" {
		return r.GetByExternalPlaceID(cafe.ExternalPlaceID)
	}

	var ratings []models.Rating
	err := r.db.
		Preload("User").
		Preload("CafeListing").
		Joins("JOIN gocafe_cafe_listings ON gocafe_cafe_listings.id = gocafe_ratings.cafe_listing_id").
		Where("COALESCE(gocafe_cafe_listings.source_cafe_id, gocafe_cafe_listings.id) = ?", rootID).
		Order("gocafe_ratings.visited_at DESC").
		Find(&ratings).Error
	return ratings, err
}

func (r *Repository) GetByExternalPlaceID(externalPlaceID string) ([]models.Rating, error) {
	var ratings []models.Rating
	err := r.db.
		Preload("User").
		Preload("CafeListing").
		Joins("JOIN gocafe_cafe_listings ON gocafe_cafe_listings.id = gocafe_ratings.cafe_listing_id").
		Where("gocafe_cafe_listings.external_place_id = ?", externalPlaceID).
		Order("gocafe_ratings.visited_at DESC").
		Find(&ratings).Error
	return ratings, err
}

func (r *Repository) GetByUserID(userID uint) ([]models.Rating, error) {
	var ratings []models.Rating
	err := r.db.
		Preload("CafeListing").
		Where("user_id = ?", userID).
		Order("visited_at DESC").
		Find(&ratings).Error
	return ratings, err
}

func (r *Repository) FindByUserAndCafe(userID uint, cafeListingID uint) (*models.Rating, error) {
	var rating models.Rating
	err := r.db.Where("user_id = ? AND cafe_listing_id = ?", userID, cafeListingID).First(&rating).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &rating, err
}

func (r *Repository) Update(id uint, updated models.Rating) error {
	var existing models.Rating
	if err := r.db.First(&existing, id).Error; err != nil {
		return err
	}
	existing.VisitedAt = updated.VisitedAt
	existing.Rating = updated.Rating
	existing.Review = updated.Review
	return r.db.Save(&existing).Error
}

func (r *Repository) Delete(id uint) error {
	result := r.db.Delete(&models.Rating{}, id)
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return result.Error
}
