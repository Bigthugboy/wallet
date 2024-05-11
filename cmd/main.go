package cmd

import (
	"log"
	"net/http"

	"github.com/Bigthugboy/wallet/cmd/route"
	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	route.HandleRoutes(r)
	http.Handle("/", r)
	log.Fatal(http.ListenAndServe("Localhost:9090", r))
}
