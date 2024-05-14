package route

import (
	"net/http"

	"github.com/Bigthugboy/wallet/pkg/config/security/jwtt"
	"github.com/Bigthugboy/wallet/pkg/internals"
	"github.com/gorilla/mux"
)

var HandleRoutes = func(route *mux.Router, service internals.Service) {
	route.HandleFunc("/register", service.RegisterHandler).Methods("POST")
	route.HandleFunc("/login", service.LoginHandler).Methods("POST")
	route.HandleFunc("/payment", authenticate(service.MakePayment)).Methods("POST")
	route.HandleFunc("/validate-payment", authenticate(service.ValidatePayment)).Methods("POST")
	route.HandleFunc("/getAll", authenticate(service.TransactionHistory)).Methods("GET")
	route.HandleFunc("/transactions/{userID}/{transactionID}", authenticate(service.GetTransactionWithID)).Methods("GET")
	route.HandleFunc("/balance/{userID}", authenticate(service.CheckBalance)).Methods("GET")
	route.HandleFunc("/exchange-rate", authenticate(service.GetExchangeRate)).Methods("GET")
}

func authenticate(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			http.Error(w, "Authorization header is required", http.StatusUnauthorized)
			return
		}
		_, err := jwtt.Parse(tokenString)
		if err != nil {
			http.Error(w, "Invalid token: "+err.Error(), http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	}
}
