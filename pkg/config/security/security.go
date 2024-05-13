package security

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/Bigthugboy/wallet/pkg/models"
)

const (
	KeycloakBaseURL      = "http://localhost:8080/auth/realms/YOUR_REALM"
	KeycloakClientID     = "YOUR_CLIENT_ID"
	KeycloakClientSecret = "YOUR_CLIENT_SECRET"
)

func RegisterUser(user *models.User) error {
	data := url.Values{}
	data.Set("email", user.Email)
	data.Set("password", user.Password)
	data.Set("enabled", "true")

	req, err := http.NewRequest("POST", KeycloakBaseURL+"/users", strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(KeycloakClientID, KeycloakClientSecret)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("failed to register user: %s", resp.Status)
	}

	return nil
}

func LoginUser(login *models.KLoginPayload) (*models.kLoginRes, error) {
	formData := url.Values{
		"clientId":      {login.ClientId},
		"client_secret": {login.ClientSecret},
		"grant_type":    {login.ClientId},
		"username":      {login.ClientId},
		"password":      {login.ClientId},
	}
	encondedFormData := formData.Encode()
	req, err := http.NewRequest("POST", KeycloakBaseURL+"/users", strings.NewReader(encondedFormData))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return errors.New("error connectiong to keylock login")
	}

	loginRes := &KloginRes{}
	json.NewDecoder(resp.Body).Decode(loginRes)
	return loginRes, nil

}
