package cafelisting

import (
	"errors"

	"github.com/khorzhenwin/go-cafe/backend/internal/models"
	"gorm.io/gorm"
)

type Storage interface {
	Create(listing *models.CafeListing) error
	GetByID(id uint) (*models.CafeListing, error)
	GetByUserID(userID uint) ([]models.CafeListing, error)
	GetByUserIDFiltered(userID uint, filter ListFilter) ([]models.CafeListing, error)
	Update(id uint, updated models.CafeListing) error
	Delete(id uint) error
}

type ListFilter struct {
	VisitStatus string
	Sort        string
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
	err := r.db.First(&listing, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &listing, err
}

func (r *Repository) GetByUserID(userID uint) ([]models.CafeListing, error) {
	return r.GetByUserIDFiltered(userID, ListFilter{})
}

func (r *Repository) GetByUserIDFiltered(userID uint, filter ListFilter) ([]models.CafeListing, error) {
	var listings []models.CafeListing
	q := r.db.Where("user_id = ?", userID)
	if filter.VisitStatus != "" {
		q = q.Where("visit_status = ?", filter.VisitStatus)
	}
	orderBy := "updated_at DESC"
	switch filter.Sort {
	case "created_desc":
		orderBy = "created_at DESC"
	case "name_asc":
		orderBy = "name ASC"
	case "name_desc":
		orderBy = "name DESC"
	case "status_asc":
		orderBy = "visit_status ASC, updated_at DESC"
	case "status_desc":
		orderBy = "visit_status DESC, updated_at DESC"
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
	existing.Description = updated.Description
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
