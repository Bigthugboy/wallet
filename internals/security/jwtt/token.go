package jwtt

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Bigthugboy/wallet/internals"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/sessions"
)

type WalletCliams struct {
	jwt.RegisteredClaims
	Email string
	ID    int64
}

var secretKey = "404E635266556A586E3272357538782F413F4428472B4B6250645367566B5970"

func Generate(email string, id int64) (string, string, error) {
	wCliams := WalletCliams{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "walletAdmin",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: &jwt.NumericDate{Time: time.Now().Add(24 * time.Hour)},
		},
		Email: email,
		ID:    id,
	}
	refWalletCliams := jwt.RegisteredClaims{
		Issuer:    "walletAdmin",
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: &jwt.NumericDate{Time: time.Now().Add(48 * time.Hour)},
	}

	// Generate JWT tokens
	walletToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, wCliams).SignedString([]byte(secretKey))
	if err != nil {
		return "", "", err
	}
	refWalletToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refWalletCliams).SignedString([]byte(secretKey))
	if err != nil {
		return "", "", err
	}

	return walletToken, refWalletToken, nil
}

func Parse(tokenString string) (*WalletCliams, error) {
	log.Print("+++++++++++++++++++++++++++++++++++++")
	token, err := jwt.ParseWithClaims(tokenString, &WalletCliams{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	log.Print("----------------------")
	if err != nil {
		return nil, fmt.Errorf("error parsing token: %v", err)
	}
	log.Print("vvvvvvvvvvvvvvvvvvvvvvvvvvv")
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	log.Print("ccccccccccccccccccccccccccccccccccc")
	claims, ok := token.Claims.(*WalletCliams)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}

// StoreSession stores session data in the cookie.
func StoreSession(w http.ResponseWriter, r *http.Request, id int64, email, password string) error {
	userInfo := &internals.UserInfo{
		ID:       id,
		Email:    email,
		Password: password,
	}
	session, err := Sessions(r).Store().Get(r, "session")
	if err != nil {
		return fmt.Errorf("error retrieving session: %v", err)
	}
	session.Values["info"] = userInfo
	if err := session.Save(r, w); err != nil {
		return fmt.Errorf("error saving session data: %v", err)
	}

	return nil
}

func Sessions(r *http.Request) *sessions.Session {
	store := sessions.NewCookieStore([]byte("wallet"))
	session, _ := store.Get(r, "session")
	return session
}
