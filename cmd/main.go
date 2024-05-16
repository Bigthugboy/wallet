package main

import (
	"encoding/gob"
	"log"
	"net/http"

	"github.com/Bigthugboy/wallet/cmd/route"
	"github.com/Bigthugboy/wallet/config"
	"github.com/jinzhu/gorm"

	"github.com/Bigthugboy/wallet/internals"
	"github.com/Bigthugboy/wallet/internals/controllers"
	"github.com/gorilla/mux"
)

var app = config.NewAppTools()

func main() {
	gob.Register(internals.User{})
	gob.Register(internals.Wallet{})

	dsn := "root:damilola@tcp(127.0.0.1:3306)/wallet?charset=utf8mb4&parseTime=True&loc=Local"
	d, err := gorm.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}

	d.AutoMigrate(&internals.User{}, &internals.Wallet{})

	app.InfoLogger.Println("*---------- Connecting to the wallet database --------")

	app.InfoLogger.Println("*---------- Starting Wallet Web Server -----------*")
	app.InfoLogger.Println("*---------- Connected to Wallet Web Server -----------*")

	srv := controllers.NewWallet(app, d)

	r := mux.NewRouter()
	route.HandleRoutes(r, srv)
	http.Handle("/", r)
	log.Fatal(http.ListenAndServe("localhost:9090", r))
}
