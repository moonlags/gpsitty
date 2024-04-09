package server

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"gpsitty/internal/database"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httprate"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger, middleware.Recoverer, httprate.LimitByIP(100, 1*time.Minute))
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
	}))

	r.Get("/auth/google/callback", s.AuthCallbackHandler)
	r.Get("/auth/google", s.AuthHandler)
	r.Get("/auth/logout", s.LogoutHandler)

	r.Get("/api/session", s.GetSession)
	r.Get("/api/link/{imei}", s.LinkDevice)

	return r
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

	var user struct {
		ID      string `json:"sub,omitempty"`
		Name    string `json:"name,omitempty"`
		Picture string `json:"picture,omitempty"`
		Email   string `json:"email,omitempty"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("failed to decode json: %v\n", err)
		return
	}

	if _, err := s.Queries.CreateUser(context.Background(), database.CreateUserParams{ID: user.ID, Name: user.Name, Email: user.Email, Avatar: user.Picture}); err != nil {
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

	user, err := s.Queries.GetUser(context.Background(), session.Values["id"].(string))
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
	session, err := s.Store.Get(r, "gpsitty")
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		log.Printf("failed to get session: %v\n", err)
		return
	}

	deviceImei := chi.URLParam(r, "imei")

	if err := s.Queries.LinkDevice(context.Background(), database.LinkDeviceParams{DeviceImei: deviceImei, Userid: session.Values["id"].(string)}); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("failed to link device: %v\n", err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// TODO: 0x14 sleep route; 0x48 restart; 0x48 shutdown; 0x61 light switch; 0x92 alarm on;0x93 alarm off
