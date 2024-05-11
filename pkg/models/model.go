package models

import (
	"time"

	middeware "github.com/Bigthugboy/wallet/cmd/middleware"
	"github.com/jinzhu/gorm"
)

var db *gorm.DB

type User struct {
	gorm.Model
	FirstName   string `gorm: ""json:"firstName"`
	LastName    string `json:"lastName"`
	Email       string `json:"email" Usuage:"required`
	Password    string `json:"password" Usuage:"required`
	PhoneNumber string `json:"phone" Usage:"required"`
}

type Payment struct {
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
	Email       string    `json:"email"`
	Amount      string    `json:"amount"`
	SubAccount  string    `json:"subaccount"`
	Currency    string    `json:"currency"`
	FirstName   string    `json:"first_name" Usage:"required,alpha"`
	LastName    string    `json:"last_name" Usage:"required,alpha"`
	DatePayed   time.Time `bson:"date_payed"`
	PhoneNumber string    `bson:"phone" Usage:"required"`
	Payment     Payment   `json:"payment"`
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

func init() {
	middeware.Connect()
	db = middeware.GetDB()
	db.AutoMigrate(&User{})
}
