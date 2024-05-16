package internals

import "net/http"

type Service interface {
	RegisterHandler(w http.ResponseWriter, r *http.Request)
	MakePayment(w http.ResponseWriter, r *http.Request)
	LoginHandler(w http.ResponseWriter, r *http.Request)
	ValidatePayment(w http.ResponseWriter, r *http.Request)
	TransactionHistory(w http.ResponseWriter, r *http.Request)
	GetTransactionWithID(w http.ResponseWriter, r *http.Request)
	CheckBalance(w http.ResponseWriter, r *http.Request)
	GetExchangeRate(w http.ResponseWriter, r *http.Request)
}
