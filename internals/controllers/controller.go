package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"log"
	"net/http"

	"github.com/Bigthugboy/wallet/config"
	"github.com/Bigthugboy/wallet/internals"

	"github.com/Bigthugboy/wallet/internals/query"
	"github.com/Bigthugboy/wallet/internals/repo"
	"github.com/Bigthugboy/wallet/internals/security/jwtt"
	"github.com/anjolabassey/Rave-go/rave"
	"golang.org/x/crypto/bcrypt"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

var secretKey = "FLWSECK_TEST-7c8c2dcff4d2a9cb96fe3a34812e1e90-X"

var card = rave.Card{
	Rave: rave.Rave{
		Live:      false,
		PublicKey: "FLWPUBK_TEST-727132610f7bb0781b0343b0b0de55e7-X",
		SecretKey: secretKey,
	},
}

type Wallet struct {
	App *config.AppTools
	DB  repo.DBStore
}

func NewWallet(app *config.AppTools, db *gorm.DB) internals.Service {
	return &Wallet{
		App: app,
		DB:  query.NewWalletDB(app, db),
	}
}

func Encrypt(password string) (string, error) {
	if password == "" {
		return "", fmt.Errorf("no input value")
	} else {
		fromPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Sprint("cannot generate encrypted password"), err
		}
		hashedString := string(fromPassword)
		return hashedString, nil
	}
}

func (wa *Wallet) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var user internals.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		log.Println("Error decoding form:", err)
		return
	}
	user.Password, _ = Encrypt(user.Password)
	if err := wa.App.Validate.Struct(&user); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		log.Println("Validation error:", err)
		return
	}
	if err := wa.DB.CreateWallet(&user); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		log.Println("Error creating wallet for user:", err)
		return
	}
	_, err := wa.DB.InsertUser(user)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		log.Println("Error adding user to database:", err)
		return
	}
	response := map[string]string{"message": "Registered Successfully"}
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		log.Println("Error encoding JSON:", err)
		return
	}
	log.Println(response)
}

func (wa *Wallet) LoginHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("========================================================")
	w.Header().Set("Content-Type", "application/json")

	var payload internals.LoginUser
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		log.Println("error decoding payload:", err)
		return
	}
	// Search for user by email
	userID, _, err := wa.DB.SearchUserByEmail(payload.Email)
	if err != nil {
		http.Error(w, "unregistered user", http.StatusUnauthorized)
		log.Println("user not registered:", payload.Email)
		return
	}

	walletToken, refWalletToken, err := jwtt.Generate(payload.Email, userID)
	if err != nil {
		http.Error(w, "failed to generate JWT tokens", http.StatusInternalServerError)
		log.Println("error generating JWT tokens:", err)
		return
	}

	response := map[string]string{
		"wallet_token":     walletToken,
		"ref_wallet_token": refWalletToken,
	}
	fmt.Printf("Token created: %s\n", walletToken)
	fmt.Printf("Token created: %s\n", refWalletToken)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (wa *Wallet) MakePayment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var payload internals.PayLoad
	log.Println(payload)
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		log.Println("Error decoding payload:", err)
		log.Println(payload)
		return
	}
	details := rave.CardChargeData{
		Cardno:        payload.CardNo,
		Cvv:           payload.Cvv,
		Expirymonth:   payload.ExpiryMonth,
		Expiryyear:    payload.ExpiryYear,
		Pin:           payload.Pin,
		Amount:        payload.Amount,
		Currency:      "NGN",
		CustomerPhone: payload.Phone,
		Firstname:     payload.FirstName,
		Lastname:      payload.LastName,
		Email:         payload.Email,
		Txref:         payload.TxRef,
		RedirectUrl:   "https://localhost:9090/checkBalance",
	}
	// Charge the card
	err, resp := card.ChargeCard(details)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		log.Println("Error charging card:", err)
		return
	}
	transaction := internals.Wallet{
		UserID:   payload.UserID,
		Balance:  details.Amount,
		Currency: details.Currency,
		Amount:   details.Amount,
		Method:   "Card Payment",
	}
	_, err = wa.DB.SavePayment(transaction)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		log.Println("Error saving payment:", err)
		return
	}
	userID := payload.UserID
	err = wa.DB.UpdateWalletBalance(userID, payload.Amount)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		log.Println("Error updating wallet balance:", err)
		return
	}
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		log.Println("Error encoding JSON:", err)
		return
	}
}

func (wa *Wallet) ValidatePayment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var validatePayload internals.ValidatePayload
	if err := json.NewDecoder(r.Body).Decode(&validatePayload); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		log.Println("Error decoding validation payload:", err)
		return
	}
	payload := rave.CardValidateData{
		Otp:       validatePayload.Otp,
		Reference: validatePayload.Reference,
		PublicKey: secretKey,
	}
	// Validate the card
	err, resp := card.ValidateCard(payload)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		log.Println("Error validating card:", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		log.Println("Error encoding JSON:", err)
		return
	}
}

func (wa *Wallet) TransactionHistory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	userID := params["userID"]

	transactions, err := wa.DB.GetAllTransactions(userID)
	if err != nil {
		http.Error(w, "Error getting transaction from database:", http.StatusInternalServerError)
		log.Printf("fail to get documents from database %v", err)
		return
	}
	if len(transactions) == 0 {
		w.WriteHeader(http.StatusOK)
		log.Println("you have made any transaction yet")
		return
	}
	res, err := json.Marshal(transactions)
	if err != nil {
		http.Error(w, "Error encoding transactions to JSON", http.StatusInternalServerError)
		log.Printf("Failed to encode transactions to JSON: %v", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(res)

}

func (wa *Wallet) GetTransactionWithID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userID := params["userID"]
	transactionID := params["transactionID"]
	// Retrieve the transaction from the database
	transaction, err := wa.DB.GetTransactionWithID(userID, transactionID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			http.Error(w, "Transaction not found", http.StatusNotFound)
		} else {
			http.Error(w, "Error getting transaction", http.StatusInternalServerError)
		}
		log.Printf("Failed to get transaction: %v", err)
		return
	}
	response, err := json.Marshal(transaction)
	if err != nil {
		http.Error(w, "Error encoding transaction to JSON", http.StatusInternalServerError)
		log.Printf("Failed to encode transaction to JSON: %v", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)

}

func (wa *Wallet) CheckBalance(w http.ResponseWriter, r *http.Request) {
	log.Print("check if endpoint was hit")
	params := mux.Vars(r)
	userID := params["userID"]
	walletID := params["walletID"]

	balance, err := wa.DB.GetWalletBalance(userID, walletID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "Wallet not found", http.StatusNotFound)
		} else {
			http.Error(w, "Error getting wallet balance", http.StatusInternalServerError)
		}
		log.Printf("Failed to get wallet balance: %v", err)
		return
	}
	log.Printf("log balance %v", balance)
	fmt.Fprintf(w, "Your balance is: %.2f", balance)
}

func (wa *Wallet) GetExchangeRate(w http.ResponseWriter, r *http.Request) {
	url := "https://api.exchangeratesapi.net/v1/exchange-rates/currency_codes"

	req, _ := http.NewRequest("GET", url, nil)
	log.Println(req)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		http.Error(w, "error getting response", http.StatusInternalServerError)
		log.Println("error getting response ", err)
	}

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	fmt.Println(res)
	fmt.Println(string(body))
}
