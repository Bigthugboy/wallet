package query

import (
	"errors"
	"fmt"

	"github.com/Bigthugboy/wallet/pkg/config"
	"github.com/Bigthugboy/wallet/pkg/internals/repo"

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

func (wa *WalletDB) CreateWallet(user *models.User) error {
	wallet := &models.Wallet{UserID: user.ID}
	if err := wa.DB.Create(&wallet).Error; err != nil {
		return err
	}
	user.Wallet = *wallet
	return nil
}

func (wa *WalletDB) GetAllTransactions(userID uint) ([]models.Transaction, error) {
	var user models.User
	if err := wa.DB.Preload("Wallet").First(&user, userID).Error; err != nil {
		return nil, fmt.Errorf("failed to find user: %v", err)
	}

	var transactions []models.Transaction
	if err := wa.DB.Model(&user.Wallet).Related(&transactions).Error; err != nil {
		return nil, fmt.Errorf("failed to get transactions: %v", err)
	}
	return transactions, nil
}

func (wa *WalletDB) GetTransactionWithID(userID, transactionID string) (models.Transaction, error) {
	var user models.User
	if err := wa.DB.Preload("Wallet").First(&user, userID).Error; err != nil {
		return models.Transaction{}, fmt.Errorf("failed to find user: %v", err)
	}
	var transaction models.Transaction
	if err := wa.DB.Where("id = ?", transactionID).First(&transaction).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return models.Transaction{}, fmt.Errorf("transaction not found with ID %d", transactionID)
		}
		return models.Transaction{}, fmt.Errorf("failed to get transaction: %v", err)
	}
	if transaction.WalletID != user.Wallet.ID {
		return models.Transaction{}, fmt.Errorf("transaction with ID %d does not belong to user with ID %d", transactionID, userID)
	}
	return transaction, nil
}

func (w *WalletDB) GetUserByID(userId string) (models.User, error) {
	user := models.User{}
	if w.DB == nil {
		return user, fmt.Errorf("database connection is not initialized")
	}
	if err := w.DB.Where("ID = ?", userId).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return user, fmt.Errorf("user not found with ID %s", userId)
		}
		return user, err
	}

	return user, nil
}
