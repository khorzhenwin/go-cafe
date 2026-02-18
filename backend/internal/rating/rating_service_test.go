package rating

import (
	"testing"

	"github.com/khorzhenwin/go-cafe/backend/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockCafeLookup struct {
	visited bool
	err     error
}

func (m *mockCafeLookup) IsListingVisited(id uint) (bool, error) {
	if m.err != nil {
		return false, m.err
	}
	return m.visited, nil
}

type mockRatingStorage struct {
	ratings    []models.Rating
	getByID    *models.Rating
	createErr  error
	updateErr  error
	deleteErr  error
}

func (m *mockRatingStorage) Create(r *models.Rating) error {
	if m.createErr != nil {
		return m.createErr
	}
	r.ID = uint(len(m.ratings) + 1)
	m.ratings = append(m.ratings, *r)
	return nil
}

func (m *mockRatingStorage) GetByID(id uint) (*models.Rating, error) { return m.getByID, nil }

func (m *mockRatingStorage) GetByCafeListingID(id uint) ([]models.Rating, error) {
	var out []models.Rating
	for _, r := range m.ratings {
		if r.CafeListingID == id {
			out = append(out, r)
		}
	}
	return out, nil
}

func (m *mockRatingStorage) GetByUserID(id uint) ([]models.Rating, error) {
	var out []models.Rating
	for _, r := range m.ratings {
		if r.UserID == id {
			out = append(out, r)
		}
	}
	return out, nil
}

func (m *mockRatingStorage) Update(id uint, updated models.Rating) error { return m.updateErr }

func (m *mockRatingStorage) Delete(id uint) error { return m.deleteErr }

func TestService_CreateRating(t *testing.T) {
	m := &mockRatingStorage{}
	svc := NewService(m, &mockCafeLookup{visited: true})
	r := &models.Rating{UserID: 1, CafeListingID: 2, Rating: 5}
	err := svc.CreateRating(r)
	require.NoError(t, err)
	require.Len(t, m.ratings, 1)
	assert.Equal(t, 5, m.ratings[0].Rating)
}

func TestService_CreateRating_RequiresVisitedCafe(t *testing.T) {
	m := &mockRatingStorage{}
	svc := NewService(m, &mockCafeLookup{visited: false})
	r := &models.Rating{UserID: 1, CafeListingID: 2, Rating: 5}
	err := svc.CreateRating(r)
	assert.ErrorIs(t, err, ErrCafeNotVisited)
	require.Len(t, m.ratings, 0)
}

func TestService_UpdateRating_NotOwner(t *testing.T) {
	m := &mockRatingStorage{getByID: &models.Rating{ID: 1, UserID: 10}}
	svc := NewService(m, nil)
	err := svc.UpdateRating(1, 99, models.Rating{Rating: 4})
	assert.ErrorIs(t, err, ErrNotOwner)
}

func TestService_DeleteRating_NotOwner(t *testing.T) {
	m := &mockRatingStorage{getByID: &models.Rating{ID: 1, UserID: 10}}
	svc := NewService(m, nil)
	err := svc.DeleteRating(1, 99)
	assert.ErrorIs(t, err, ErrNotOwner)
}
