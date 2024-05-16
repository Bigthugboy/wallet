package route

import (
	middewares "github.com/Bigthugboy/wallet/cmd/middlewares"
	"github.com/Bigthugboy/wallet/internals"
	"github.com/gorilla/mux"
)

var HandleRoutes = func(route *mux.Router, service internals.Service) {
	route.HandleFunc("/register", service.RegisterHandler).Methods("POST")
	route.HandleFunc("/login", service.LoginHandler).Methods("POST")
	route.HandleFunc("/payment", middewares.Authenticate(service.MakePayment)).Methods("POST")
	route.HandleFunc("/validate-payment", middewares.Authenticate(service.ValidatePayment)).Methods("POST")
	route.HandleFunc("/transactions/{userID}/{transactionID}", middewares.Authenticate(service.GetTransactionWithID)).Methods("GET")
	route.HandleFunc("/balance/{userID}/{walletID}", middewares.Authenticate(service.CheckBalance)).Methods("GET")
	route.HandleFunc("/exchange-rate", service.GetExchangeRate).Methods("GET")
}
