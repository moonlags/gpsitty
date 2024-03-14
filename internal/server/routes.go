package server

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"gpsitty/internal/database"
)

func (s *Server) RegisterRoutes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", s.HealthHandler)

	mux.HandleFunc("/auth/google/callback", s.AuthCallbackHandler)
	mux.HandleFunc("/logout", s.LogoutHandler)
	mux.HandleFunc("/auth/google", s.AuthHandler)

	mux.HandleFunc("/session", s.GetSession)

	handler := CorsMiddleware(mux)

	return handler
}

func CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (s *Server) HealthHandler(w http.ResponseWriter, r *http.Request) {
	jsonResp, _ := json.Marshal(s.DB.Health())
	_, _ = w.Write(jsonResp)
}

func (s *Server) AuthCallbackHandler(w http.ResponseWriter, r *http.Request) {
	token, err := s.Conf.Exchange(context.Background(), r.FormValue("code"))
	if err != nil {
		log.Fatalf("failed to exchange: %v\n", err)
	}

	if !token.Valid() {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	client := s.Conf.Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("failed to get response from google: %v\n", err)
		return
	}
	defer resp.Body.Close()

	var user database.User
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("failed to decode json: %v\n", err)
		return
	}

	if err := s.DB.CreateUser(user); err != nil {
		log.Fatal("failed to create user:", err)
	}

	session, err := s.Store.Get(r, "gpsitty")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("failed to get session: %v\n", err)
		return
	}

	session.Values["id"] = user.ID
	if err := session.Save(r, w); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("failed to save session: %v\n", err)
		return
	}

	http.Redirect(w, r, "http://localhost:5173/", http.StatusFound)
}

func (s *Server) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:   "gpsitty",
		MaxAge: -1,
	})

	http.Redirect(w, r, "http://localhost:5173", http.StatusTemporaryRedirect)
}

func (s *Server) AuthHandler(w http.ResponseWriter, r *http.Request) {
	url := s.Conf.AuthCodeURL("state")

	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (s *Server) GetSession(w http.ResponseWriter, r *http.Request) {
	session, err := s.Store.Get(r, "gpsitty")
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		log.Printf("failed to get session: %v\n", err)
		return
	}

	user, err := s.DB.GetUser(session.Values["id"].(string))
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		log.Printf("failed to get user: %v\n", err)
		return
	}

	jsonResp, err := json.Marshal(user)
	if err != nil {
		log.Fatalf("failed to marshal json: %v\n", err)
	}

	if _, err := w.Write(jsonResp); err != nil {
		log.Fatalf("failed to write to response: %v\n", err)
	}
}

func (s *Server) LinkDevice(w http.ResponseWriter, r *http.Request) {
	panic("not implemented") // todo
}
