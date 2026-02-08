package cafelisting

import (
	"errors"
	"testing"

	"github.com/khorzhenwin/go-cafe/backend/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockCafeStorage struct {
	listings   []models.CafeListing
	getByID    *models.CafeListing
	getByIDErr error
	createErr error
	updateErr error
	deleteErr error
}

func (m *mockCafeStorage) Create(c *models.CafeListing) error {
	if m.createErr != nil {
		return m.createErr
	}
	c.ID = uint(len(m.listings) + 1)
	m.listings = append(m.listings, *c)
	return nil
}

func (m *mockCafeStorage) GetByID(id uint) (*models.CafeListing, error) {
	if m.getByIDErr != nil {
		return nil, m.getByIDErr
	}
	return m.getByID, nil
}

func (m *mockCafeStorage) GetByUserID(userID uint) ([]models.CafeListing, error) {
	var out []models.CafeListing
	for _, l := range m.listings {
		if l.UserID == userID {
			out = append(out, l)
		}
	}
	return out, nil
}

func (m *mockCafeStorage) Update(id uint, updated models.CafeListing) error { return m.updateErr }

func (m *mockCafeStorage) Delete(id uint) error { return m.deleteErr }

func TestService_CreateListing(t *testing.T) {
	m := &mockCafeStorage{}
	svc := NewService(m)
	listing := &models.CafeListing{UserID: 1, Name: "Cafe A", Address: "123"}
	err := svc.CreateListing(listing)
	require.NoError(t, err)
	require.Len(t, m.listings, 1)
	assert.Equal(t, "Cafe A", m.listings[0].Name)
}

func TestService_UpdateListing_NotOwner(t *testing.T) {
	m := &mockCafeStorage{getByID: &models.CafeListing{ID: 1, UserID: 10}}
	svc := NewService(m)
	err := svc.UpdateListing(1, 99, models.CafeListing{Name: "X"})
	assert.ErrorIs(t, err, ErrNotOwner)
}

func TestService_UpdateListing_Owner(t *testing.T) {
	m := &mockCafeStorage{getByID: &models.CafeListing{ID: 1, UserID: 10}}
	svc := NewService(m)
	err := svc.UpdateListing(1, 10, models.CafeListing{Name: "New Name"})
	require.NoError(t, err)
	assert.NoError(t, m.updateErr)
}

func TestService_DeleteListing_NotOwner(t *testing.T) {
	m := &mockCafeStorage{getByID: &models.CafeListing{ID: 1, UserID: 10}}
	svc := NewService(m)
	err := svc.DeleteListing(1, 99)
	assert.ErrorIs(t, err, ErrNotOwner)
}

func TestService_DeleteListing_NotFound(t *testing.T) {
	m := &mockCafeStorage{getByID: nil}
	svc := NewService(m)
	err := svc.DeleteListing(1, 10)
	require.Error(t, err)
}

func TestService_DeleteListing_Owner(t *testing.T) {
	m := &mockCafeStorage{getByID: &models.CafeListing{ID: 1, UserID: 10}}
	svc := NewService(m)
	err := svc.DeleteListing(1, 10)
	require.NoError(t, err)
}
