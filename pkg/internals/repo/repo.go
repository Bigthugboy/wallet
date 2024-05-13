package repo

import "github.com/Bigthugboy/wallet/pkg/models"

type DBStore interface {
	InsertUser(user models.User) (int64, error)
	SearchUserByEmail(email string) (int64, string, error)
	GetUserByID(userId string) (models.User, error)
	SavePayment(transaction models.Transaction) (int64, error)
	CreateWallet(User *models.User) error
	GetAllTransactions(userId string) ([]models.Transaction, error)
	GetTransactionWithID(userID, transactionID string) (models.Transaction, error)
	UpdateWalletBalance(userID uint, amount float64) error
}
