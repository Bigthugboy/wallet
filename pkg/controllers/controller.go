package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/Bigthugboy/wallet/pkg/config"
	"github.com/Bigthugboy/wallet/pkg/internals"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"

	"github.com/Bigthugboy/wallet/pkg/internals/repo"
	"github.com/Bigthugboy/wallet/pkg/models"
	"github.com/go-playground/validator"
)

var secretKey = os.Getenv("SECRECT_KEY")
var apiKey = os.Getenv("API_KEY")

type Wallet struct {
	App *config.AppTools
	DB  repo.DBStore
}

func NewWallet(app *config.AppTools, db *gorm.DB) internals.Service {
	return &Wallet{
		App: app,
		DB:  repo.NewWalletDB(app, db),
	}
}

func (wa *Wallet) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var user models.User
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		log.Println("Error parsing form:", err)
		return
	}
	user.Password, _ = config.Encrypt(user.Password)
	// Validation
	if err := wa.App.Validate.Struct(&user); err != nil {
		if _, ok := err.(*validator.InvalidValidationError); !ok {
			http.Error(w, "Bad request", http.StatusBadRequest)
			log.Println("Validation error:", err)
			return
		}
	}
	// Create wallet for the user
	if err := wa.DB.CreateWallet(&user); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	track, err := wa.DB.InsertUser(user)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		log.Println("Error adding user to database:", err)
		return
	}
	switch track {
	case 1:
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	case 0:
		response := map[string]string{"message": "Registered Successfully"}

		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			log.Println("Error encoding JSON:", err)
			return
		}
	}

}

func (wa *Wallet) MakePayment(w http.ResponseWriter, r *http.Request) {
	paymentRequest := models.PaymentRequest{}
	err := json.NewDecoder(r.Body).Decode(&paymentRequest)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		log.Println("Error decoding JSON:", err)
		return
	}
	// requestData := map[string]interface{}{
	// 	// Populate request payload according to Flutterwave's API documentation
	// }
	requestDataBytes, err := json.Marshal(paymentRequest)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		log.Println("Error encoding JSON:", err)
		return
	}
	client := &http.Client{}
	req, err := http.NewRequest("POST", "https://api.flutterwave.com/v3/payments", strings.NewReader(string(requestDataBytes)))
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		log.Println("Error creating request:", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer"+secretKey)

	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		log.Println("Error making request:", err)
		return
	}
	defer resp.Body.Close()

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		log.Println("Error reading response body:", err)
		return
	}
	var paymentResponse models.PaymentResponse
	err = json.Unmarshal(responseBody, &paymentResponse)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		log.Println("Error decoding JSON response:", err)
		return
	}

	fmt.Println("Payment response:", paymentResponse)
	json.NewEncoder(w).Encode(paymentResponse)
}

func (wa *Wallet) ValidatePayment(w http.ResponseWriter, r *http.Request) {
	reference := r.URL.Query().Get("reference")
	url := fmt.Sprintf("https://api.flutterwave.com/v3/transactions/%s/verify", reference)
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+secretKey)

	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		http.Error(w, "Transaction not successful", resp.StatusCode)
		return
	}
	var validateResp models.ValidateResponse
	err = json.NewDecoder(resp.Body).Decode(&validateResp)
	if err != nil {
		http.Error(w, "Invalid response", http.StatusInternalServerError)
		return
	}

	if validateResp.Data.Amount == "" {
		http.Error(w, "Failed to parse amount", http.StatusBadRequest)
		return
	}

	fmt.Println("Validation response:", validateResp)
	json.NewEncoder(w).Encode(validateResp)
}

func (wa *Wallet) TransactionHistory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	//refactor after testing to use keyclock session to get userId
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		log.Println("Error parsing form:", err)
		return
	}
	var user models.User
	transactions, err := wa.DB.GetAllTransactions(user.ID)
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
	var exchangeData models.Data
	if err := json.NewDecoder(r.Body).Decode(&exchangeData); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		log.Println("Error decoding JSON:", err)
		return
	}
	requestDataBytes, err := json.Marshal(exchangeData)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		log.Println("Error encoding JSON:", err)
		return
	}
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://api.exchangeratesapi.net/v1/exchange-rates/latest", bytes.NewReader(requestDataBytes))
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		log.Println("Error creating request:", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		log.Println("Error making request:", err)
		return
	}
	defer resp.Body.Close()
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		log.Println("Error reading response body:", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	w.Write(responseBody)
}
