package cafelisting

import (
	"github.com/khorzhenwin/go-cafe/backend/internal/models"
	"gorm.io/gorm"
)

type Service struct {
	store Storage
}

func NewService(store Storage) *Service {
	return &Service{store: store}
}

func (s *Service) GetByID(id uint) (*models.CafeListing, error) {
	return s.store.GetByID(id)
}

func (s *Service) GetByUserID(userID uint) ([]models.CafeListing, error) {
	return s.store.GetByUserID(userID)
}

func (s *Service) CreateListing(listing *models.CafeListing) error {
	return s.store.Create(listing)
}

func (s *Service) UpdateListing(id uint, userID uint, updated models.CafeListing) error {
	existing, err := s.store.GetByID(id)
	if err != nil || existing == nil {
		return err
	}
	if existing.UserID != userID {
		return ErrNotOwner
	}
	return s.store.Update(id, updated)
}

func (s *Service) DeleteListing(id uint, userID uint) error {
	existing, err := s.store.GetByID(id)
	if err != nil {
		return err
	}
	if existing == nil {
		return gorm.ErrRecordNotFound
	}
	if existing.UserID != userID {
		return ErrNotOwner
	}
	return s.store.Delete(id)
}
