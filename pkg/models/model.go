package models

import (
	"time"

	middeware "github.com/Bigthugboy/wallet/cmd/middlewares"
	"github.com/jinzhu/gorm"
)

var db *gorm.DB

type User struct {
	gorm.Model
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	Email       string `json:"email" gorm:"unique;not null"`
	Password    string `json:"password" gorm:"not null"`
	PhoneNumber string `json:"phone" gorm:"not null"`
	Wallet      Wallet `json:"wallet" gorm:"foreignkey:UserID"`
}

type Wallet struct {
	gorm.Model
	UserID       uint
	Balance      float64 `json:"balance" gorm:"default:0"`
	Transactions []Transaction
}

type Transaction struct {
	gorm.Model
	WalletID  uint
	Amount    float64   `json:"amount"`
	Timestamp time.Time `json:"timestamp"`
	Method    string    `json:"method"`
	Type      string    `json:"type"`
}

type ExchangeRates struct {
	USD float64 `json:"usd"`
	EUR float64 `json:"eur"`
}
type PayLoad struct {
	FirstName   string  `json:"first_name" Usage:"required,alpha"`
	LastName    string  `json:"last_name" Usage:"required,alpha"`
	Amount      float64 `json:"amount"`
	TxRef       string  `json:"tx_ref"`
	Email       string  `json:"email"`
	Phone       string  `json:"phone"`
	Currency    string  `json:"currency"`
	CardNo      string  `json:"card_no"`
	Cvv         string  `json:"cvv"`
	Pin         string  `json:"pin"`
	ExpiryMonth string  `json:"expiry_month"`
	ExpiryYear  string  `json:"expiry_year"`
}
type ValidatePayload struct {
	Reference string `json:"transaction_reference"`
	Otp       string `json:"otp"`
	PublicKey string `json:"PBFPubKey"`
}

type Data struct {
	Base         string `json:"base"`
	To           string `json:"to"`
	From         string `json:"from"`
	Date         string `json:"date"`
	CurrencyCode string `json:"curreny_code"`
}

func init() {
	middeware.Connect()
	db = middeware.GetDB()
	db.AutoMigrate(&User{}, &Wallet{}, &Transaction{})

}
