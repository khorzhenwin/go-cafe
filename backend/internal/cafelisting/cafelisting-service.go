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
	return s.GetByUserIDFiltered(userID, "", "")
}

func (s *Service) GetByUserIDFiltered(userID uint, visitStatus, sort string) ([]models.CafeListing, error) {
	status, err := normalizeVisitStatus(visitStatus)
	if err != nil {
		return nil, err
	}
	if visitStatus == "" {
		status = ""
	}
	return s.store.GetByUserIDFiltered(userID, ListFilter{
		VisitStatus: status,
		Sort:        sort,
	})
}

func (s *Service) CreateListing(listing *models.CafeListing) error {
	status, err := normalizeVisitStatus(listing.VisitStatus)
	if err != nil {
		return err
	}
	listing.VisitStatus = status
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
	status, err := normalizeVisitStatus(updated.VisitStatus)
	if err != nil {
		return err
	}
	updated.VisitStatus = status
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

func (s *Service) IsListingVisited(id uint) (bool, error) {
	existing, err := s.store.GetByID(id)
	if err != nil {
		return false, err
	}
	if existing == nil {
		return false, gorm.ErrRecordNotFound
	}
	return existing.VisitStatus == VisitStatusVisited, nil
}
