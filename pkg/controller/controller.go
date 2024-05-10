package controller

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/BigthugBoy/wallet/pkg/models"
	"github.com/Bigthugboy/wallet/internal"
	"github.com/Bigthugboy/wallet/pkg/config"
	"github.com/Bigthugboy/wallet/pkg/repo"
	"github.com/go-playground/validator"
	"github.com/jinzhu/gorm"
)

type Wallet struct {
	App *config.AppTools
	DB  repo.DBStore
}

func NewWallet(app *config.AppTools, db *gorm.DB) internal.MainStore {
	return &Wallet{
		App: app,
		DB:  repo.NewTravasDB(app, db),
	}
}

func (wa *Wallet) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var user models.User
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		log.Println("Error parsing form:", err)
		return
	}
	user.Password, _ = config.Encrypt(user.Password)
	// Validation
	if err := wa.App.Validate.Struct(&user); err != nil {
		if _, ok := err.(*validator.InvalidValidationError); !ok {
			http.Error(w, "Bad request", http.StatusBadRequest)
			log.Println("Validation error:", err)
			return
		}
	}
	track, err := wa.DB.InsertCustomer(user)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		log.Println("Error adding user to database:", err)
		return
	}
	switch track {
	case 1:
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	case 0:
		response := map[string]string{"message": "Registered Successfully"}

		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			log.Println("Error encoding JSON:", err)
			return
		}
	}
}
