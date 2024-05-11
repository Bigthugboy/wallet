package controller

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/Bigthugboy/wallet/pkg/config"

	"github.com/Bigthugboy/wallet/pkg/internals/repo"
	"github.com/Bigthugboy/wallet/pkg/models"
	"github.com/go-playground/validator"
)

type Wallet struct {
	App *config.AppTools
	DB  repo.DBStore
}

var secretKey = os.Getenv("SECRECT_KEY")

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
	track, err := wa.DB.InsertUser(user)
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

func (wa *Wallet) MakePayment(w http.ResponseWriter, r *http.Request) {
    paymentRequest := models.PaymentRequest{}

    // Read and parse the request body
    err := json.NewDecoder(r.Body).Decode(&paymentRequest)
    if err != nil {
        http.Error(w, "Bad request", http.StatusBadRequest)
        log.Println("Error decoding JSON:", err)
        return
    }

    requestData := map[string]interface{}{
        // Populate request payload according to Flutterwave's API documentation
    }

    requestDataBytes, err := json.Marshal(requestData)
    if err != nil {
        http.Error(w, "Internal server error", http.StatusInternalServerError)
        log.Println("Error encoding JSON:", err)
        return
    }
    client := &http.Client{}
    req, err := http.NewRequest("POST", "https://api.flutterwave.com/v3/payments", strings.NewReader(string(requestDataBytes)))
    if err != nil {
        http.Error(w, "Internal server error", http.StatusInternalServerError)
        log.Println("Error creating request:", err)
        return
    }
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", "Bearer" +secretKey)

    resp, err := client.Do(req)
    if err != nil {
        http.Error(w, "Internal server error", http.StatusInternalServerError)
        log.Println("Error making request:", err)
        return
    }
    defer resp.Body.Close()

    responseBody, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        http.Error(w, "Internal server error", http.StatusInternalServerError)
        log.Println("Error reading response body:", err)
        return
    }
    var paymentResponse models.PaymentResponse
    err = json.Unmarshal(responseBody, &paymentResponse)
    if err != nil {
        http.Error(w, "Internal server error", http.StatusInternalServerError)
        log.Println("Error decoding JSON response:", err)
        return
    }

    // Handle successful payment response
    // You may want to send the payment response back to the client or perform additional processing
    fmt.Println("Payment response:", paymentResponse)
    // Send the payment response back to the client
    json.NewEncoder(w).Encode(paymentResponse)
}
func (wa *Wallet) ValidatePayment(w http.ResponseWriter, r *http.Request)  {
	 
		reference := r.Query("reference")

		url := fmt.Sprintf("https://api.paystack.co/transaction/verify/%s", reference)
		client := &http.Client{}

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			_ = http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+secretKey)
		resp, _ := client.Do(req)

		validateResp := &models.ValidateResponse{}
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				http.Error(w,http.StatusInternalServerError,"error": "Failed to close response body")
				return
			}
		}(resp.Body)
		err = ctx.ShouldBindJSON(&validateResp)
		if err != nil {
			http.Error(w,http.StatusInternalServerError,"error": "invalid response")
			return
		}
		if resp.StatusCode != http.StatusOK {
			http.Error(w,resp.StatusCode,"error": "Transaction not successful")
			return
		}

		if validateResp.Data.Amount == " " {
			http.Error(http.StatusBadRequest, "error": "Failed to parse amount")
			return
		}
	}



