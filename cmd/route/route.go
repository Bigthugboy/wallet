package route

import (
	"github.com/Bigthugboy/wallet/pkg/internals"
	"github.com/gorilla/mux"
)

var HandleRoutes = func(route *mux.Router, service internals.Service) {
	route.HandleFunc("/register", service.RegisterHandler()).Methods("POST")
	route.HandleFunc("/register", service.MakePayment()).Methods("POST")
	route.HandleFunc("/register", service.ValidatePayment()).Methods("POST")
	route.HandleFunc("/register", service.TransactionHistory()).Methods("GET")
	route.HandleFunc("/transactions/{userID}/{transactionID}", service.GetTransactionWithID()).Methods("GET")
	route.HandleFunc("/balance/{userID}", service.CheckBalance()).Methods("GET")
	route.HandleFunc("/register", service.GetExchangeRate()).Methods("GET")

}
