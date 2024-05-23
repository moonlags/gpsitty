package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/mail"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

type AuthData struct {
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
}

type (
	securedHandler func(http.ResponseWriter, *http.Request, string)
	authHandler    func(http.ResponseWriter, *http.Request, AuthData)
)

func (s *Server) SecuredRoute(next securedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("token")
		if err != nil {
			HttpError(w, "no token cookie", http.StatusUnauthorized, err)
			return
		}

		token, err := jwt.Parse(cookie.Value, func(t *jwt.Token) (interface{}, error) {
			if t.Method != jwt.SigningMethodHS256 {
				return nil, fmt.Errorf("unexpected signing method %v", t.Method.Alg())
			}

			return []byte(os.Getenv("SIGNING_KEY")), nil
		})
		if err != nil || !token.Valid {
			HttpError(w, "invalid token", http.StatusUnauthorized, err)
			return
		}

		issuer, err := token.Claims.GetIssuer()
		if err != nil {
			HttpError(w, "token has invalid claims", http.StatusBadRequest, err)
			return
		}
		next(w, r, issuer)
	}
}

func (s *Server) AuthRoute(next authHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var data AuthData
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			HttpError(w, "failed to decode json", http.StatusBadRequest, err)
			return
		}

		if _, err := mail.ParseAddress(data.Email); err != nil || data.Password == "" {
			HttpError(w, "not all request values are present", http.StatusBadRequest, err)
			return
		}

		next(w, r, data)
	}
}
