package security

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

const (
	KeycloakBaseURL      = "http://localhost:8080/auth/realms/YOUR_REALM"
	KeycloakClientID     = "YOUR_CLIENT_ID"
	KeycloakClientSecret = "YOUR_CLIENT_SECRET"
)

func RegisterUser(user User) error {
	data := url.Values{}
	data.Set("username", user.Username)
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
