package query

import (
	"errors"
	"fmt"

	"github.com/Bigthugboy/wallet/pkg/config"
	"github.com/Bigthugboy/wallet/pkg/internal/repo"

	"github.com/Bigthugboy/wallet/pkg/models"
	"github.com/jinzhu/gorm"
)

type WalletDB struct {
	App *config.AppTools
	DB  *gorm.DB
}

func NewWalletDB(app *config.AppTools, db *gorm.DB) repo.DBStore {
	return &WalletDB{
		App: app,
		DB:  db,
	}
}

func (w *WalletDB) InsertUser(user models.User) (int64, error) {
	if w.DB == nil {
		return -1, fmt.Errorf("database connection is not initialized")
	}

	var existingUser models.User
	if err := w.DB.Where("email = ?", user.Email).First(&existingUser).Error; err != nil && err != gorm.ErrRecordNotFound {
		return -1, err
	}
	if existingUser.ID != 0 {
		return -1, fmt.Errorf("user with email '%s' already exists", user.Email)
	}
	result := w.DB.Create(&user)
	if err := result.Error; err != nil {
		return -1, err
	}

	return result.RowsAffected, nil
}

func (w *WalletDB) SearchUserByEmail(email string) (int64, string, error) {
	if w.DB == nil {
		return -1, "", fmt.Errorf("database connection is not initialized")
	}

	user := models.User{}
	if err := w.DB.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return -1, "", nil
		}
		return -1, "", err
	}

	return int64(user.ID), user.FirstName, nil
}
