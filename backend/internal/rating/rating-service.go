package rating

import (
	"strings"

	"github.com/khorzhenwin/go-cafe/backend/internal/models"
)

type Service struct {
	store      Storage
	cafeLookup interface {
		IsListingVisited(id uint) (bool, error)
	}
}

func NewService(store Storage, cafeLookup interface {
	IsListingVisited(id uint) (bool, error)
}) *Service {
	return &Service{store: store, cafeLookup: cafeLookup}
}

func (s *Service) GetByID(id uint) (*models.Rating, error) {
	return s.store.GetByID(id)
}

func (s *Service) GetByCafeListingID(cafeListingID uint) ([]models.Rating, error) {
	return s.store.GetByCafeListingID(cafeListingID)
}

func (s *Service) GetByExternalPlaceID(externalPlaceID string) ([]models.Rating, error) {
	return s.store.GetByExternalPlaceID(externalPlaceID)
}

func (s *Service) GetByUserID(userID uint) ([]models.Rating, error) {
	return s.store.GetByUserID(userID)
}

func (s *Service) CreateRating(rating *models.Rating) error {
	if err := validateRating(rating); err != nil {
		return err
	}
	if s.cafeLookup != nil {
		visited, err := s.cafeLookup.IsListingVisited(rating.CafeListingID)
		if err != nil {
			return err
		}
		if !visited {
			return ErrCafeNotVisited
		}
	}
	existing, err := s.store.FindByUserAndCafe(rating.UserID, rating.CafeListingID)
	if err != nil {
		return err
	}
	if existing != nil {
		return ErrDuplicateRating
	}
	return s.store.Create(rating)
}

func (s *Service) UpdateRating(id uint, userID uint, updated models.Rating) error {
	existing, err := s.store.GetByID(id)
	if err != nil || existing == nil {
		return err
	}
	if existing.UserID != userID {
		return ErrNotOwner
	}
	if err := validateRating(&updated); err != nil {
		return err
	}
	return s.store.Update(id, updated)
}

func (s *Service) DeleteRating(id uint, userID uint) error {
	existing, err := s.store.GetByID(id)
	if err != nil || existing == nil {
		return err
	}
	if existing.UserID != userID {
		return ErrNotOwner
	}
	return s.store.Delete(id)
}

func validateRating(rating *models.Rating) error {
	if rating.Rating < 1 || rating.Rating > 5 {
		return ErrInvalidRatingValue
	}
	rating.Review = strings.TrimSpace(rating.Review)
	return nil
}
