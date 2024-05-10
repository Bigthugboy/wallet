package models

import (
	middeware "github.com/Bigthugboy/wallet/cmd/middleware"
	"github.com/jinzhu/gorm"
)

var db *gorm.DB

type User struct {
	gorm.Model
	FirstName string `gorm: ""json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type Transaction struct {
	gorm.Model
	ID        uint
	UserID    uint
	Amount    float64
	Timestamp string
}
type ExchangeRates struct {
	USD float64 `json:"usd"`
	EUR float64 `json:"eur"`
}

type PaymentRequest struct {
	PublicKey     string  `json:"public_key"`
	EncryptionKey string  `json:"encryption_key"`
	Currency      string  `json:"currency"`
	Amount        float64 `json:"amount"`
	Email         string  `json:"email"`
	FullName      string  `json:"full_name"`
	TxRef         string  `json:"tx_ref"`
}

func init() {
	middeware.Connect()
	db = middeware.GetDB()
	db.AutoMigrate(&User{})
}
