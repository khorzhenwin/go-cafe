package user

import (
	"errors"
	"testing"

	"github.com/khorzhenwin/go-cafe/backend/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockStorage struct {
	users       []models.User
	createErr   error
	getByID     *models.User
	getByIDErr  error
	getByEmail  *models.User
	getByEmailErr error
	updateErr   error
	deleteErr   error
}

func (m *mockStorage) Create(u *models.User) error {
	if m.createErr != nil {
		return m.createErr
	}
	if u.ID == 0 {
		u.ID = uint(len(m.users) + 1)
	}
	m.users = append(m.users, *u)
	return nil
}

func (m *mockStorage) GetAll() ([]models.User, error) { return m.users, nil }

func (m *mockStorage) GetByID(id uint) (*models.User, error) {
	if m.getByIDErr != nil {
		return nil, m.getByIDErr
	}
	return m.getByID, nil
}

func (m *mockStorage) GetByEmail(email string) (*models.User, error) {
	if m.getByEmailErr != nil {
		return nil, m.getByEmailErr
	}
	return m.getByEmail, nil
}

func (m *mockStorage) Update(id uint, updated models.User) error { return m.updateErr }

func (m *mockStorage) Delete(id uint) error { return m.deleteErr }

func TestService_CreateWithPassword(t *testing.T) {
	m := &mockStorage{}
	svc := NewService(m)
	id, err := svc.CreateWithPassword("a@b.com", "Alice", "pass123")
	require.NoError(t, err)
	assert.NotZero(t, id)
	require.Len(t, m.users, 1)
	assert.Equal(t, "a@b.com", m.users[0].Email)
	assert.Equal(t, "Alice", m.users[0].Name)
	assert.NotEmpty(t, m.users[0].PasswordHash)
	assert.NotEqual(t, "pass123", m.users[0].PasswordHash)
}

func TestService_GetByEmailForAuth_NotFound(t *testing.T) {
	m := &mockStorage{getByEmail: nil}
	svc := NewService(m)
	id, hash, err := svc.GetByEmailForAuth("x@y.com")
	require.NoError(t, err)
	assert.Zero(t, id)
	assert.Empty(t, hash)
}

func TestService_GetByEmailForAuth_Found(t *testing.T) {
	m := &mockStorage{getByEmail: &models.User{ID: 3, Email: "u@v.com", PasswordHash: "hashed"}}
	svc := NewService(m)
	id, hash, err := svc.GetByEmailForAuth("u@v.com")
	require.NoError(t, err)
	assert.Equal(t, uint(3), id)
	assert.Equal(t, "hashed", hash)
}

func TestService_CreateWithPassword_PropagatesError(t *testing.T) {
	m := &mockStorage{createErr: errors.New("db error")}
	svc := NewService(m)
	_, err := svc.CreateWithPassword("a@b.com", "A", "p")
	assert.Error(t, err)
}
