package db

import (
	"github.com/Bigthugboy/wallet/pkg/config"
	"github.com/Bigthugboy/wallet/pkg/internal"

	"github.com/jinzhu/gorm"
)

type wallet struct {
	App *config.AppTools
	DB  *gorm.DB
}

func NewWallet(app *config.AppTools, db *gorm.DB) internal.mainstore {
	return &wallet{
		App: app,
		DB:  db,
	}
}
