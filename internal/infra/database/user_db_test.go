package database

import (
	"testing"

	"github.com/ivandersr/products-api-go/internal/entity"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestCreateUser(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	if err != nil {
		t.Error(err)
	}
	db.AutoMigrate(&entity.User{})
	user, _ := entity.NewUser("John", "j@j.com", "123456")
	userDB := NewUserDB(db)
	err = userDB.Create(user)
	assert.Nil(t, err)

	var userFound entity.User
	err = db.First(&userFound, "id = ?", user.ID).Error
	assert.Nil(t, err)
	assert.Equal(t, "John", userFound.Name)
	assert.Equal(t, "j@j.com", userFound.Email)
	assert.NotEmpty(t, userFound.Password)
}

func TestFindByEmail(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	if err != nil {
		t.Error(err)
	}
	db.AutoMigrate(&entity.User{})
	user, _ := entity.NewUser("John", "j@j.com", "123456")
	userDB := NewUserDB(db)
	db.Create(user)
	userFound, err := userDB.FindByEmail("j@j.com")
	assert.Nil(t, err)
	assert.NotNil(t, userFound)
	assert.Equal(t, user.ID, userFound.ID)
	assert.Equal(t, "j@j.com", userFound.Email)
}
