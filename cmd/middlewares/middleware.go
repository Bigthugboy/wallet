package middewares

import (
	"log"
	"net/http"
	"strings"

	"github.com/Bigthugboy/wallet/internals/security/jwtt"

	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func Authenticate(next http.HandlerFunc) http.HandlerFunc {
	log.Print("log authentication process")
	return func(w http.ResponseWriter, r *http.Request) {
		log.Print("check authorization")
		authHeader := r.Header.Get("Authorization")
		log.Println(authHeader)
		if authHeader == "" {
			http.Error(w, "Authorization header is required", http.StatusUnauthorized)
			return
		}
		log.Print("split token from bearer")

		parts := strings.Split(authHeader, " ")
		log.Print(parts)
		if len(parts) != 2 || parts[0] != "Bearer" {
			log.Println(parts[0])
			log.Println(parts[1])

			log.Print(len(parts[0]))
			http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
			return
		}

		// Extract the JWT token
		log.Println(" print token to be passed")
		tokenString := parts[1]

		_, err := jwtt.Parse(tokenString)
		log.Print("token passed is " + tokenString)
		if err != nil {
			http.Error(w, "Invalid token: "+err.Error(), http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	}
}
