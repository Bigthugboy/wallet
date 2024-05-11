package repo

import "github.com/Bigthugboy/wallet/pkg/models"

type DBStore interface {
	InsertUser(user models.User) (int64, error)
	SearchUserByEmail(email string) (int64, string, error)
	CreateWallet(User *models.User) error
	GetAllTransactions(userId uint) ([]models.Transaction, error)
	GetTransactionWithID(userID, transactionID uint) (models.Transaction, error)
}
