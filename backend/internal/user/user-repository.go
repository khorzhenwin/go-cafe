package user

import (
	"errors"

	"github.com/khorzhenwin/go-cafe/backend/internal/models"
	"gorm.io/gorm"
)

type Storage interface {
	Create(user *models.User) error
	GetAll() ([]models.User, error)
	GetByID(id uint) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	Update(id uint, updated models.User) error
	Delete(id uint) error
}

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(u *models.User) error {
	return r.db.Create(u).Error
}

func (r *Repository) GetAll() ([]models.User, error) {
	var users []models.User
	err := r.db.Find(&users).Error
	return users, err
}

func (r *Repository) GetByID(id uint) (*models.User, error) {
	var u models.User
	err := r.db.First(&u, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &u, err
}

func (r *Repository) GetByEmail(email string) (*models.User, error) {
	var u models.User
	err := r.db.Where("email = ?", email).First(&u).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &u, err
}

func (r *Repository) Update(id uint, updated models.User) error {
	var existing models.User
	if err := r.db.First(&existing, id).Error; err != nil {
		return err
	}
	existing.Email = updated.Email
	existing.Name = updated.Name
	return r.db.Save(&existing).Error
}

func (r *Repository) Delete(id uint) error {
	result := r.db.Delete(&models.User{}, id)
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return result.Error
}
