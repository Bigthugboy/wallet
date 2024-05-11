package models

import (
	"time"

	middeware "github.com/Bigthugboy/wallet/cmd/middleware"
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

type PaymentRequest struct {
	Email       string    `json:"email"`
	Amount      string    `json:"amount"`
	SubAccount  string    `json:"subaccount"`
	Currency    string    `json:"currency"`
	FirstName   string    `json:"first_name" Usage:"required,alpha"`
	LastName    string    `json:"last_name" Usage:"required,alpha"`
	DatePayed   time.Time `bson:"date_payed"`
	PhoneNumber string    `bson:"phone" Usage:"required"`
	// Payment     Payment   `json:"payment"`
}

type PaymentResponse struct {
	Status  bool         `json:"status"`
	Message string       `json:"message"`
	Data    ResponseData `json:"data"`
}

type Authorizations struct {
	AuthorizationCode string `json:"authorizationCode"`
}

type ResponseData struct {
	AuthorizationUrl string         `json:"authorization_url"`
	AccessCode       string         `json:"access_code"`
	Reference        string         `json:"reference"`
	Amount           string         `json:"amount"`
	Status           bool           `json:"status"`
	Authorization    Authorizations `json:"authorization"`
	StatusCode       string         `json:"status_code"`
}

type ValidateResponse struct {
	Status  bool         `json:"status"`
	Message string       `json:"message"`
	Data    ResponseData `json:"data"`
}

type Data struct {
	Base          string `json:"base"`
	To            string `json:"to"`
	From          string `json:"from"`
	Date          string `json:"date"`
	Currency_code string `json:"curreny_code"`
}

func init() {
	middeware.Connect()
	db = middeware.GetDB()
	db.AutoMigrate(&User{})
}
