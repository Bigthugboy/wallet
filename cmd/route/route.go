package route

import (
	controller "github.com/Bigthugboy/wallet/pkg/controllers"
	"github.com/gorilla/mux"
)

var HandleRoutes = func(route *mux.Router) {
	route.HandleFunc("/register", controller.RegisterHandler).Methods("POST")
	route.HandleFunc("/register", controller.MakePayment).Methods("POST")
	route.HandleFunc("/register", controller.ValidatePayment).Methods("POST")
	route.HandleFunc("/register", controller.TransactionHistory).Methods("GET")
	route.HandleFunc("/register", controller.GetTransactionWithID).Methods("GET")
	route.HandleFunc("/register", controller.CheckBalance).Methods("GET")
	route.HandleFunc("/register", controller.GetExchangeRate).Methods("GET")

}
