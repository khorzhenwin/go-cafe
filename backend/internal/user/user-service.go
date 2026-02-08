package user

import (
	"fmt"

	"github.com/khorzhenwin/go-cafe/backend/internal/auth"
	"github.com/khorzhenwin/go-cafe/backend/internal/models"
)

type Service struct {
	store Storage
}

func NewService(store Storage) *Service {
	return &Service{store: store}
}

func (s *Service) FindAll() ([]models.User, error) {
	return s.store.GetAll()
}

func (s *Service) GetByID(id uint) (*models.User, error) {
	return s.store.GetByID(id)
}

func (s *Service) GetByEmail(email string) (*models.User, error) {
	return s.store.GetByEmail(email)
}

// GetByEmailForAuth returns user id and password hash for auth (implements auth.LoginFinder).
func (s *Service) GetByEmailForAuth(email string) (id uint, passwordHash string, err error) {
	u, err := s.store.GetByEmail(email)
	if err != nil || u == nil {
		return 0, "", err
	}
	return u.ID, u.PasswordHash, nil
}

// CreateWithPassword creates a user with hashed password (implements auth.RegisterCreator).
func (s *Service) CreateWithPassword(email, name, password string) (id uint, err error) {
	hash, err := auth.HashPassword(password)
	if err != nil {
		return 0, fmt.Errorf("hash password: %w", err)
	}
	u := &models.User{Email: email, Name: name, PasswordHash: hash}
	if err := s.store.Create(u); err != nil {
		return 0, err
	}
	return u.ID, nil
}

func (s *Service) CreateUser(user *models.User) error {
	return s.store.Create(user)
}

func (s *Service) UpdateUser(id uint, updated models.User) error {
	return s.store.Update(id, updated)
}

func (s *Service) DeleteUser(id uint) error {
	return s.store.Delete(id)
}
