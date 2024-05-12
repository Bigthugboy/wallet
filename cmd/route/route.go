package route

import (
	"github.com/Bigthugboy/wallet/pkg/internals"
	"github.com/gorilla/mux"
)

var HandleRoutes = func(route *mux.Router, service internals.Service) {
	route.HandleFunc("/register", service.RegisterHandler).Methods("POST")
	route.HandleFunc("/payment", service.MakePayment).Methods("POST")
	route.HandleFunc("/validate-payment", service.ValidatePayment).Methods("POST")
	route.HandleFunc("/getAll", service.TransactionHistory).Methods("GET")
	route.HandleFunc("/transactions/{userID}/{transactionID}", service.GetTransactionWithID).Methods("GET")
	route.HandleFunc("/balance/{userID}", service.CheckBalance).Methods("GET")
	route.HandleFunc("/", service.GetExchangeRate).Methods("GET")

}
