package rating

import "github.com/khorzhenwin/go-cafe/backend/internal/models"

type Service struct {
	store Storage
}

func NewService(store Storage) *Service {
	return &Service{store: store}
}

func (s *Service) GetByID(id uint) (*models.Rating, error) {
	return s.store.GetByID(id)
}

func (s *Service) GetByCafeListingID(cafeListingID uint) ([]models.Rating, error) {
	return s.store.GetByCafeListingID(cafeListingID)
}

func (s *Service) GetByUserID(userID uint) ([]models.Rating, error) {
	return s.store.GetByUserID(userID)
}

func (s *Service) CreateRating(rating *models.Rating) error {
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
