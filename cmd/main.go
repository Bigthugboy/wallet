package main

import (
	"encoding/gob"
	"log"
	"net/http"

	middeware "github.com/Bigthugboy/wallet/cmd/middlewares"
	"github.com/Bigthugboy/wallet/cmd/route"
	"github.com/Bigthugboy/wallet/pkg/config"

	"github.com/Bigthugboy/wallet/pkg/controllers"
	"github.com/Bigthugboy/wallet/pkg/models"
	"github.com/gorilla/mux"
)

var app = config.NewAppTools()

func main() {
	gob.Register(models.User{})
	gob.Register(models.Wallet{})
	gob.Register(models.Transaction{})

	app.InfoLogger.Println("*---------- Connecting to the wallet database --------")
	app.InfoLogger.Println("*---------- Starting Wallet Web Server -----------*")
	app.InfoLogger.Println("*---------- Connected to Wallet Web Server -----------*")
	db := middeware.GetDB()

	srv := controllers.NewWallet(app, db)

	r := mux.NewRouter()
	route.HandleRoutes(r, srv)
	http.Handle("/", r)
	log.Fatal(http.ListenAndServe("localhost:9090", r))
}
