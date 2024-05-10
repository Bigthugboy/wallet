package controller

import (
	"github.com/Bigthugboy/wallet/internal"
	"github.com/Bigthugboy/wallet/pkg/config"
	"github.com/Bigthugboy/wallet/pkg/repo"
	"github.com/jinzhu/gorm"
)

type Wallet struct {
	App *config.AppTools
	DB  repo.DBStore
}

func NewWallet(app *config.AppTools, db *gorm.DB) internal.MainStore {
	return &Wallet{
		App: app,
		DB:  repo.NewTravasDB(app, db),
	}
}
