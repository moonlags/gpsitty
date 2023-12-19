package server

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/markbates/goth/gothic"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
	}))

	r.Get("/", s.HelloWorldHandler)
	r.Get("/health", s.healthHandler)

	r.Get("/auth/{provider}/callback", s.authCallbackHandler)
	r.Get("/logout/{provider}", s.logoutHandler)
	r.Get("/auth/{provider}", s.authHandler)

	r.Get("/session", s.getSession)

	return r
}

func (s *Server) HelloWorldHandler(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]string)
	resp["message"] = "Hello World"

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("error handling JSON marshal. Err: %v", err)
	}

	_, _ = w.Write(jsonResp)
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	jsonResp, _ := json.Marshal(s.db.Health())
	_, _ = w.Write(jsonResp)
}

func (s *Server) authCallbackHandler(w http.ResponseWriter, r *http.Request) {
	provider := chi.URLParam(r, "provider")
	r = r.WithContext(context.WithValue(context.Background(), "provider", provider))

	user, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		http.Error(w, "failed to complete user auth", http.StatusBadRequest)
		return
	}

	if err := gothic.StoreInSession("userid", user.UserID, r, w); err != nil {
		http.Error(w, "failed to store user id", http.StatusUnauthorized)
		return
	}

	if err := s.db.CreateUser(user.UserID, user.Email, user.AvatarURL); err != nil {
		log.Println("failed to create user, but thats maybe ok:", err)
	}

	http.Redirect(w, r, "http://localhost:5173/", http.StatusFound)
}

func (s *Server) logoutHandler(w http.ResponseWriter, r *http.Request) {
	provider := chi.URLParam(r, "provider")
	r = r.WithContext(context.WithValue(context.Background(), "provider", provider))

	if err := gothic.Logout(w, r); err != nil {
		http.Error(w, "failed to logout", http.StatusUnauthorized)
		return
	}

	http.Redirect(w, r, "http://localhost:5173", http.StatusTemporaryRedirect)
}

func (s *Server) authHandler(w http.ResponseWriter, r *http.Request) {
	provider := chi.URLParam(r, "provider")
	r = r.WithContext(context.WithValue(context.Background(), "provider", provider))

	user, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		gothic.BeginAuthHandler(w, r)
		return
	}

	if err := gothic.StoreInSession("userID", user.UserID, r, w); err != nil {
		http.Error(w, "can not store user id in session", http.StatusUnauthorized)
		return
	}

	if err := s.db.CreateUser(user.UserID, user.Email, user.AvatarURL); err != nil {
		log.Println("failed to create user, but thats maybe ok:", err)
	}

	http.Redirect(w, r, "http://localhost:5173", http.StatusFound)
}

func (s *Server) getSession(w http.ResponseWriter, r *http.Request) {
	userid, err := gothic.GetFromSession("userid", r)
	if err != nil {
		http.Error(w, "failed to get user id from session", http.StatusUnauthorized)
		return
	}

	user, err := s.db.GetUserWithDevices(userid)
	if err != nil {
		log.Println("failed to get user with devices:", err)
		http.Error(w, "failed to get user with devices", http.StatusUnauthorized)
	}

	jsonResp, err := json.Marshal(user)
	if err != nil {
		log.Fatal("failed to marshal json", err)
	}

	if _, err := w.Write(jsonResp); err != nil {
		log.Fatal("failed to write to response", err)
	}
}

func (s *Server) linkDevice(w http.ResponseWriter, r *http.Request) {
	// todo
}
