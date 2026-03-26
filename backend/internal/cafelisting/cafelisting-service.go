package cafelisting

import (
	"strings"

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

func (s *Service) ListDiscovery(query, city, sort string, limit int) ([]models.CafeListing, error) {
	return s.store.ListDiscovery(DiscoveryFilter{
		Query: strings.TrimSpace(query),
		City:  strings.TrimSpace(city),
		Sort:  strings.TrimSpace(sort),
		Limit: limit,
	})
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
	if err := sanitizeListing(listing); err != nil {
		return err
	}

	status, err := normalizeVisitStatus(listing.VisitStatus)
	if err != nil {
		return err
	}
	listing.VisitStatus = status

	if listing.SourceCafeID != nil {
		source, err := s.store.GetByID(*listing.SourceCafeID)
		if err != nil {
			return err
		}
		if source == nil {
			return gorm.ErrRecordNotFound
		}
		if source.SourceCafeID != nil {
			listing.SourceCafeID = source.SourceCafeID
		}
	}

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
	if err := sanitizeListing(&updated); err != nil {
		return err
	}
	updated.VisitStatus = status
	updated.SourceCafeID = existing.SourceCafeID
	updated.SourceProvider = existing.SourceProvider
	updated.ExternalPlaceID = existing.ExternalPlaceID
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

func sanitizeListing(listing *models.CafeListing) error {
	listing.Name = strings.TrimSpace(listing.Name)
	listing.Address = strings.TrimSpace(listing.Address)
	listing.City = strings.TrimSpace(listing.City)
	listing.Neighborhood = strings.TrimSpace(listing.Neighborhood)
	listing.Description = strings.TrimSpace(listing.Description)
	listing.ImageURL = strings.TrimSpace(listing.ImageURL)
	listing.SourceProvider = strings.TrimSpace(listing.SourceProvider)
	listing.ExternalPlaceID = strings.TrimSpace(listing.ExternalPlaceID)

	if listing.Name == "" {
		return ErrInvalidCafeName
	}

	if (listing.Latitude == nil) != (listing.Longitude == nil) {
		return ErrInvalidCoordinates
	}

	if listing.Latitude != nil && (*listing.Latitude < -90 || *listing.Latitude > 90) {
		return ErrInvalidCoordinates
	}

	if listing.Longitude != nil && (*listing.Longitude < -180 || *listing.Longitude > 180) {
		return ErrInvalidCoordinates
	}

	return nil
}
