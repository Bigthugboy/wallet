package models

import (
	"fmt"
	"log"

	middeware "github.com/Bigthugboy/wallet/cmd/middlewares"
	"github.com/jinzhu/gorm"
)

var db *gorm.DB

type User struct {
	gorm.Model
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`

	Email       string `json:"email" gorm:"unique;not null"`
	Password    string `json:"password" gorm:"not null"`
	PhoneNumber string `json:"phone" gorm:"not null"`
	Wallet      Wallet `json:"wallet"`
}
type UserInfo struct {
	ID       int64
	Email    string
	Password string
}
type KLoginPayload struct {
	ClientID     string
	Username     string
	Password     string
	GrantType    string
	ClientSecret string
}

type LoginUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type KLoginRes struct {
	AccessToken string `json:"access_token"`
}

type Wallet struct {
	gorm.Model
	UserID       uint
	Balance      float64       `json:"balance" gorm:"-"`
	Transactions []Transaction `json:"transactions"`
	Currency     string        `json:"currency"`
}

type Transaction struct {
	gorm.Model
	UserID   uint
	WalletID uint
	Amount   float64 `json:"amount"`
	Method   string  `json:"method"`
	Type     string  `json:"type"`
	Currency string  `json:"currency"`
}

type PayLoad struct {
	FirstName   string  `json:"first_name" Usage:"required,alpha"`
	LastName    string  `json:"last_name" Usage:"required,alpha"`
	Amount      float64 `json:"amount"`
	TxRef       string  `json:"tx_ref"`
	Email       string  `json:"email"`
	Phone       string  `json:"phone"`
	Currency    string  `json:"currency"`
	CardNo      string  `json:"cardno"`
	Cvv         string  `json:"cvv"`
	Pin         string  `json:"pin"`
	ExpiryMonth string  `json:"expirymonth"`
	ExpiryYear  string  `json:"expiryyear"`
}
type ValidatePayload struct {
	Reference string `json:"tx_ref"`
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
	db.Model(&Wallet{}).Association("Transactions")

}

func dropTables(db *gorm.DB) {
	if err := db.DropTableIfExists(&User{}, &Wallet{}, &Transaction{}).Error; err != nil {
		log.Fatalf("Error dropping tables: %v", err)
	}
	fmt.Println("Tables dropped successfully")
}
