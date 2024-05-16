package security

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strings"

	"github.com/Bigthugboy/wallet/internals"
)

const (
	KeycloakBaseURL      = "http://localhost:8080/admin/"
	KeycloakClientID     = "wallet-test"
	KeycloakClientSecret = "ZARfeL8Krkk6mAZEXWkrZuzTE6hUeB4"
)

type AuthService interface {
	Login(payload *internals.KLoginPayload) (*internals.KLoginRes, error)
	ExtractUserInfo(accessToken string) (*internals.UserInfo, error)
}

type Client struct {
	httpClient      *http.Client
	keycloakBaseURL string
}

func NewClient(httpClient *http.Client, keycloakBaseURL string) *Client {
	return &Client{
		httpClient:      httpClient,
		keycloakBaseURL: keycloakBaseURL,
	}
}

func (c *Client) Login(payload *internals.KLoginPayload) (*internals.KLoginRes, error) {
	formData := url.Values{
		"client_id":     {payload.ClientID},
		"client_secret": {payload.ClientSecret},
		"grant_type":    {"password"},
		"username":      {payload.Username},
		"password":      {payload.Password},
	}
	encodedFormData := formData.Encode()

	req, err := http.NewRequest("POST", "http://localhost:8080/realms/Test/protocol/openid-connect/token", strings.NewReader(encodedFormData))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("failed to login user")
	}

	loginRes := &internals.KLoginRes{}
	err = json.NewDecoder(resp.Body).Decode(loginRes)
	if err != nil {
		return nil, err
	}

	return loginRes, nil
}

func (c *Client) ExtractUserInfo(accessToken string) (*internals.UserInfo, error) {
	req, err := http.NewRequest("GET", c.keycloakBaseURL+"/protocol/openid-connect/userinfo", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("failed to extract user info")
	}

	userInfo := &internals.UserInfo{}
	err = json.NewDecoder(resp.Body).Decode(userInfo)
	if err != nil {
		return nil, err
	}

	return userInfo, nil
}
