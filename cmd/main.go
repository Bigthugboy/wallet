package main

import (
	"log"
	"net/http"

	middeware "github.com/Bigthugboy/wallet/cmd/middleware"
	"github.com/Bigthugboy/wallet/cmd/route"
	"github.com/Bigthugboy/wallet/pkg/config"
	"github.com/Bigthugboy/wallet/pkg/controllers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

var app *config.AppTools

func main() {
	app.InfoLogger.Println("*---------- Connecting to the wallet database --------")

	err := godotenv.Load()
	if err != nil {
		app.ErrorLogger.Fatalf("cannot load up the env file : %v", err)
	}
	app.InfoLogger.Println("*---------- Starting Wallet Web Server -----------*")
	db := middeware.GetDB()

	srv := controllers.NewWallet(app, db)
	r := mux.NewRouter()
	route.HandleRoutes(r, srv)
	http.Handle("/", r)
	log.Fatal(http.ListenAndServe("Localhost:9090", r))
}
