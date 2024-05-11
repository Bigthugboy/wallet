package repo

import (
	"github.com/Bigthugboy/wallet/pkg/config"
	_ "github.com/Bigthugboy/wallet/pkg/internals/repo"

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
