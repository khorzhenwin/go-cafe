package cafelisting

import (
	"errors"
	"strings"

	"github.com/khorzhenwin/go-cafe/backend/internal/models"
	"gorm.io/gorm"
)

type Storage interface {
	Create(listing *models.CafeListing) error
	GetByID(id uint) (*models.CafeListing, error)
	ListDiscovery(filter DiscoveryFilter) ([]models.CafeListing, error)
	GetByUserID(userID uint) ([]models.CafeListing, error)
	GetByUserIDFiltered(userID uint, filter ListFilter) ([]models.CafeListing, error)
	Update(id uint, updated models.CafeListing) error
	Delete(id uint) error
}

type ListFilter struct {
	VisitStatus string
	Sort        string
}

type DiscoveryFilter struct {
	Query string
	City  string
	Sort  string
	Limit int
}

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(c *models.CafeListing) error {
	return r.db.Create(c).Error
}

func (r *Repository) GetByID(id uint) (*models.CafeListing, error) {
	var listing models.CafeListing
	err := r.baseListingQuery().
		Where("gocafe_cafe_listings.id = ?", id).
		First(&listing).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &listing, err
}

func (r *Repository) ListDiscovery(filter DiscoveryFilter) ([]models.CafeListing, error) {
	var listings []models.CafeListing
	q := r.baseListingQuery().
		Where("gocafe_cafe_listings.source_cafe_id IS NULL")

	if query := strings.TrimSpace(filter.Query); query != "" {
		likeQuery := "%" + strings.ToLower(query) + "%"
		q = q.Where(
			`LOWER(gocafe_cafe_listings.name) LIKE ? OR LOWER(gocafe_cafe_listings.address) LIKE ? OR LOWER(gocafe_cafe_listings.city) LIKE ? OR LOWER(gocafe_cafe_listings.neighborhood) LIKE ? OR LOWER(gocafe_cafe_listings.description) LIKE ?`,
			likeQuery,
			likeQuery,
			likeQuery,
			likeQuery,
			likeQuery,
		)
	}

	if city := strings.TrimSpace(filter.City); city != "" {
		q = q.Where("LOWER(gocafe_cafe_listings.city) = ?", strings.ToLower(city))
	}

	orderBy := "COALESCE(stats.review_count, 0) DESC, COALESCE(stats.avg_rating, 0) DESC, gocafe_cafe_listings.updated_at DESC"
	switch filter.Sort {
	case "newest":
		orderBy = "gocafe_cafe_listings.created_at DESC"
	case "name_asc":
		orderBy = "gocafe_cafe_listings.name ASC"
	case "rating_desc":
		orderBy = "COALESCE(stats.avg_rating, 0) DESC, COALESCE(stats.review_count, 0) DESC, gocafe_cafe_listings.updated_at DESC"
	}

	limit := filter.Limit
	if limit <= 0 {
		limit = 18
	}
	if limit > 60 {
		limit = 60
	}

	err := q.Order(orderBy).Limit(limit).Find(&listings).Error
	return listings, err
}

func (r *Repository) GetByUserID(userID uint) ([]models.CafeListing, error) {
	return r.GetByUserIDFiltered(userID, ListFilter{})
}

func (r *Repository) GetByUserIDFiltered(userID uint, filter ListFilter) ([]models.CafeListing, error) {
	var listings []models.CafeListing
	q := r.baseListingQuery().Where("gocafe_cafe_listings.user_id = ?", userID)
	if filter.VisitStatus != "" {
		q = q.Where("gocafe_cafe_listings.visit_status = ?", filter.VisitStatus)
	}
	orderBy := "gocafe_cafe_listings.updated_at DESC"
	switch filter.Sort {
	case "created_desc":
		orderBy = "gocafe_cafe_listings.created_at DESC"
	case "name_asc":
		orderBy = "gocafe_cafe_listings.name ASC"
	case "name_desc":
		orderBy = "gocafe_cafe_listings.name DESC"
	case "status_asc":
		orderBy = "gocafe_cafe_listings.visit_status ASC, gocafe_cafe_listings.updated_at DESC"
	case "status_desc":
		orderBy = "gocafe_cafe_listings.visit_status DESC, gocafe_cafe_listings.updated_at DESC"
	}
	err := q.Order(orderBy).Find(&listings).Error
	return listings, err
}

func (r *Repository) Update(id uint, updated models.CafeListing) error {
	var existing models.CafeListing
	if err := r.db.First(&existing, id).Error; err != nil {
		return err
	}
	existing.Name = updated.Name
	existing.Address = updated.Address
	existing.City = updated.City
	existing.Neighborhood = updated.Neighborhood
	existing.Description = updated.Description
	existing.ImageURL = updated.ImageURL
	existing.Latitude = updated.Latitude
	existing.Longitude = updated.Longitude
	existing.SourceProvider = updated.SourceProvider
	existing.ExternalPlaceID = updated.ExternalPlaceID
	existing.VisitStatus = updated.VisitStatus
	return r.db.Save(&existing).Error
}

func (r *Repository) Delete(id uint) error {
	result := r.db.Delete(&models.CafeListing{}, id)
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return result.Error
}

func (r *Repository) baseListingQuery() *gorm.DB {
	statsQuery := r.db.Table("gocafe_cafe_listings AS stats_cafes").
		Select(`
			COALESCE(stats_cafes.source_cafe_id, stats_cafes.id) AS root_id,
			COALESCE(ROUND(AVG(CAST(gocafe_ratings.rating AS numeric)), 2), 0) AS avg_rating,
			COUNT(gocafe_ratings.id) AS review_count
		`).
		Joins("LEFT JOIN gocafe_ratings ON gocafe_ratings.cafe_listing_id = stats_cafes.id").
		Group("COALESCE(stats_cafes.source_cafe_id, stats_cafes.id)")

	return r.db.
		Table("gocafe_cafe_listings").
		Select(`
			gocafe_cafe_listings.*,
			COALESCE(stats.avg_rating, 0) AS avg_rating,
			COALESCE(stats.review_count, 0) AS review_count
		`).
		Joins("LEFT JOIN (?) AS stats ON stats.root_id = COALESCE(gocafe_cafe_listings.source_cafe_id, gocafe_cafe_listings.id)", statsQuery)
}
