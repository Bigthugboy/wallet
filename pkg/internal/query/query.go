package query

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Bigthugboy/wallet/pkg/models"
	"github.com/jinzhu/gorm"
)

var db *gorm.DB

func (w *wallet) InsertUser(user models.User) (int64, error) {
	_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if db == nil {
		return -1, fmt.Errorf("database connection is not initialized")
	}
	var existingUser models.User
	if err := db.Where("email = ?", user.Email).First(&existingUser).Error; err != nil && err != gorm.ErrRecordNotFound {
		return -1, err
	}
	if existingUser.ID != 0 {
		return -1, fmt.Errorf("user with email '%s' already exists", user.Email)
	}
	result := db.Create(&user)
	if err := result.Error; err != nil {
		return -1, err
	}

	return result.RowsAffected, nil
}

// get customer by email
func (w *wallet) SearchUserByEmail(email string) (int64, string, error) {
	_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if w.DB == nil {
		return -1, "", fmt.Errorf("database connection is not initialized")
	}

	user := models.User{}
	if err := w.DB.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			w.App.ErrorLogger.Println("no document found for this query")
			return -1, "", nil
		}
		return -1, "", err
	}

	return int64(user.ID), user.FirstName, nil
}
