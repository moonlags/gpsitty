package auth

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	maxAge = 3600 * 24 * 7
	isProd = false
)

func SignJwt(id string) string {
	claims := &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(maxAge * time.Second)),
		Issuer:    id,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(os.Getenv("SIGNING_KEY")))
	if err != nil {
		log.Fatalf("failed to sing jwt: %v\n", err)
	}

	return signed
}

func SetAuthCookie(w http.ResponseWriter, value string) {
	cookie := &http.Cookie{
		Name:     "token",
		Value:    value,
		Path:     "/",
		MaxAge:   maxAge,
		HttpOnly: true,
		Secure:   isProd,
	}
	http.SetCookie(w, cookie)
}
