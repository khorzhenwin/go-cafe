package user

import (
	"testing"

	"github.com/khorzhenwin/go-cafe/backend/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestRepository_CreateAndGetByID(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&models.User{}))

	repo := NewRepository(db)
	u := &models.User{Email: "a@b.com", Name: "A", PasswordHash: "hash"}
	err = repo.Create(u)
	require.NoError(t, err)
	require.NotZero(t, u.ID)

	got, err := repo.GetByID(u.ID)
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, "a@b.com", got.Email)
	assert.Equal(t, "hash", got.PasswordHash)
}

func TestRepository_GetByID_NotFound(t *testing.T) {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	_ = db.AutoMigrate(&models.User{})
	repo := NewRepository(db)

	got, err := repo.GetByID(999)
	require.NoError(t, err)
	assert.Nil(t, got)
}

func TestRepository_GetByEmail(t *testing.T) {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	_ = db.AutoMigrate(&models.User{})
	repo := NewRepository(db)
	_ = repo.Create(&models.User{Email: "x@y.com", Name: "X", PasswordHash: "h"})

	got, err := repo.GetByEmail("x@y.com")
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, "x@y.com", got.Email)
}

func TestRepository_Update(t *testing.T) {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	_ = db.AutoMigrate(&models.User{})
	repo := NewRepository(db)
	u := &models.User{Email: "old@x.com", Name: "Old", PasswordHash: "h"}
	_ = repo.Create(u)

	err := repo.Update(u.ID, models.User{Email: "new@x.com", Name: "New"})
	require.NoError(t, err)

	got, _ := repo.GetByID(u.ID)
	require.NotNil(t, got)
	assert.Equal(t, "new@x.com", got.Email)
	assert.Equal(t, "New", got.Name)
}

func TestRepository_Delete(t *testing.T) {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	_ = db.AutoMigrate(&models.User{})
	repo := NewRepository(db)
	u := &models.User{Email: "d@d.com", Name: "D", PasswordHash: "h"}
	_ = repo.Create(u)

	err := repo.Delete(u.ID)
	require.NoError(t, err)

	got, _ := repo.GetByID(u.ID)
	assert.Nil(t, got)
}
