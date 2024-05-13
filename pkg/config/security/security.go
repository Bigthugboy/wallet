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
	KeycloakBaseURL      = "http://localhost:8080/admin/realms/wallet"
	KeycloakClientID     = "wallet-auth"
	KeycloakClientSecret = "AsPw5x28Gph3xfopQL20UQywUtDfaX7r"
)

func RegisterUser(user *models.User) error {
	data := url.Values{}
	data.Set("email", user.Email)
	data.Set("password", user.Password)
	data.Set("client_id", KeycloakClientID)
	data.Add("grant_type", "password")
	data.Add("scope", "offline_access")
	data.Add("client_secret", KeycloakClientSecret)

	req, err := http.NewRequest("POST", "http://localhost:8080/realms/wallet/protocol/openid-connect/token", strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

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

func LoginUser(login *models.KLoginPayload) (*models.KLoginRes, error) {

	formData := url.Values{
		"client_id":     {KeycloakClientID},
		"client_secret": {KeycloakClientSecret},
		"grant_type":    {"password"},
		"username":      {login.Email},
		"password":      {login.Password},
	}
	encodedFormData := formData.Encode()

	req, err := http.NewRequest("POST", KeycloakBaseURL+"/protocol/openid-connect/token", strings.NewReader(encodedFormData))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("failed to login user")
	}

	loginRes := &models.KLoginRes{}
	err = json.NewDecoder(resp.Body).Decode(loginRes)
	if err != nil {
		return nil, err
	}

	return loginRes, nil
}

func ExtractUserInfo(accessToken string) (*models.UserInfo, error) {
	req, err := http.NewRequest("GET", KeycloakBaseURL+"/protocol/openid-connect/userinfo", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("failed to extract user info")
	}

	userInfo := &models.UserInfo{}
	err = json.NewDecoder(resp.Body).Decode(userInfo)
	if err != nil {
		return nil, err
	}

	return userInfo, nil
}
