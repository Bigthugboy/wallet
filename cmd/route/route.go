package route

import (
	"context"
	"net/http"
	"strings"

	"github.com/Bigthugboy/wallet/pkg/config/security"
	"github.com/Bigthugboy/wallet/pkg/internals"
	"github.com/gorilla/mux"
)

var HandleRoutes = func(route *mux.Router, service internals.Service) {
	route.HandleFunc("/register", service.RegisterHandler).Methods("POST")
	route.HandleFunc("/login", service.LoginHandler).Methods("GET")
	route.HandleFunc("/payment", authenticate(service.MakePayment)).Methods("POST")
	route.HandleFunc("/validate-payment", authenticate(service.ValidatePayment)).Methods("POST")
	route.HandleFunc("/getAll", authenticate(service.TransactionHistory)).Methods("GET")
	route.HandleFunc("/transactions/{userID}/{transactionID}", authenticate(service.GetTransactionWithID)).Methods("GET")
	route.HandleFunc("/balance/{userID}", authenticate(service.CheckBalance)).Methods("GET")
	route.HandleFunc("/exchange-rate", authenticate(service.GetExchangeRate)).Methods("GET")
}

func authenticate(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		accessToken := extractAccessTokenFromRequest(r)
		if accessToken == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		userInfo, err := security.ExtractUserInfo(accessToken)
		if err != nil {
			http.Error(w, "Failed to extract user info", http.StatusInternalServerError)
			return
		}
		ctx := context.WithValue(r.Context(), "userID", userInfo.ID)

		r = r.WithContext(ctx)

		handler(w, r)
	}
}

func extractAccessTokenFromRequest(r *http.Request) string {
	accessToken := r.Header.Get("Authorization")
	if accessToken == "" {
		return ""
	}
	// Assuming the format is "Bearer <access_token>"
	parts := strings.Split(accessToken, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return ""
	}
	return parts[1]
}
