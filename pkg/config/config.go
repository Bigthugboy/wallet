package config

import (
	"encoding/json"
	"fmt"
	"io"

	"log"
	"net/http"
	"os"
	"regexp"

	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
)

type AppTools struct {
	ErrorLogger *log.Logger
	InfoLogger  *log.Logger
	Validate    *validator.Validate
}

func NewAppTools() *AppTools {
	return &AppTools{
		log.New(os.Stdout, "[ Error ]", log.LstdFlags|log.Lshortfile),
		log.New(os.Stdout, "[ info ]", log.LstdFlags|log.Lshortfile),
		validator.New(),
	}
}

var rxEmail = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

func IsEmailValid(email string) bool {
	if len(email) < 3 || len(email) > 252 || !rxEmail.MatchString(email) {
		return false
	}

	return true
}

func ParseBody(r *http.Request, x interface{}) {
	if body, err := io.ReadAll(r.Body); err == nil {
		if err := json.Unmarshal([]byte(body), x); err != nil {
			return
		}
	}
}

func Encrypt(password string) (string, error) {
	if password == "" {
		return "", fmt.Errorf("no input value")
	} else {
		fromPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Sprint("cannot generate encrypted password"), err
		}
		hashedString := string(fromPassword)
		return hashedString, nil
	}

}

func Verify(password, hashedPassword string) (bool, error) {
	if password == "" || hashedPassword == "" {
		return false, nil
	}
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return false, fmt.Errorf("invalid string comparision : %v", err)
		}
		return false, err
	}
	return true, nil

}
