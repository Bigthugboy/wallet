package repo

import "github.com/Bigthugboy/wallet/pkg/models"

type DBStore interface {
	InsertUser(user models.User) (int64, error)
	SearchUserByEmail(email string) (int64, string, error)
	// GetTransaction(userId int)([]models.Transaction)

}
