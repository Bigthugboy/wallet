package db

import (
	"github.com/Bigthugboy/wallet/pkg/config"
	"github.com/Bigthugboy/wallet/pkg/internal/repo"

	"github.com/jinzhu/gorm"
)

type wallet struct {
	App *config.AppTools
	DB  *gorm.DB
}

func NewWalletDB(app *config.AppTools, db *gorm.DB) repo.DBStore {
	return &wallet{
		App: app,
		DB:  db,
	}
}
