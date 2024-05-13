package controllers

import (
	"encoding/json"
	"fmt"
	"io"

	"log"
	"net/http"

	"github.com/Bigthugboy/wallet/pkg/config"
	"github.com/Bigthugboy/wallet/pkg/internals"
	"github.com/anjolabassey/Rave-go/rave"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"

	"github.com/Bigthugboy/wallet/pkg/internals/query"
	"github.com/Bigthugboy/wallet/pkg/internals/repo"
	"github.com/Bigthugboy/wallet/pkg/models"
)

var secretKey = "FLWSECK_TEST-7c8c2dcff4d2a9cb96fe3a34812e1e90-X"
var apiKey = "joNy4QC92c72ri4K"

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
func (wa *Wallet) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		log.Println("Error decoding form:", err)
		return
	}
	user.Password, _ = config.Encrypt(user.Password)
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
}

func (wa *Wallet) MakePayment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var payload models.PayLoad
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		log.Println("Error decoding payload:", err)
		return
	}
	// Prepare card charge details
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
	transaction := models.Transaction{
		UserID:   userID,
		Amount:   details.Amount,
		Type:     details.Chargetype,
		Currency: details.Currency,
		Method:   "Card Payment",
	}

	_, err = wa.DB.SavePayment(transaction)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		log.Println("Error saving payment:", err)
		return
	}
	err = wa.DB.UpdateWalletBalance(userID, payload.Amount)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		log.Println("Error updating wallet balance:", err)
		return
	}

	// Write the response
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		log.Println("Error encoding JSON:", err)
		return
	}
}

func (wa *Wallet) ValidatePayment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var validatePayload models.ValidatePayload
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
	//refactor after testing to use keyclock session to get userId
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
	// Extract user ID from the request context or session
	params := mux.Vars(r)
	userID := params["userID"]
	// Retrieve the user from the database
	user, err := wa.DB.GetUserByID(userID)
	if err != nil {
		http.Error(w, "Error getting user", http.StatusInternalServerError)
		log.Printf("Failed to get user: %v", err)
		return
	}
	balance := user.Wallet.Balance
	fmt.Fprintf(w, "Your balance is: %.2f", balance)

}
func (wa *Wallet) GetExchangeRate(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	apikey := params["apikey"]
	url := "https://api.exchangeratesapi.net/v1/exchange-rates/currency_codes "

	req, _ := http.NewRequest("GET", url, nil)

	res, err := http.DefaultClient.Do(req)
	req.Header.Set("Authorization", "Bearer "+apikey)
	defer func() {
		if res != nil && res.Body != nil {
			res.Body.Close()
		}
	}()
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		log.Println("Error making request:", err)
		return
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		log.Println("Error reading response body:", err)
		return
	}

	respondWithJSON(w, http.StatusOK, body)
}

func respondWithJSON(w http.ResponseWriter, status int, data []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(data)
}
