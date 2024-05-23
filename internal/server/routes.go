package server

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"gpsitty/internal/auth"
	"gpsitty/internal/database"

	"github.com/alexedwards/argon2id"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httprate"
	"github.com/oklog/ulid/v2"
)

var clientHost string = os.Getenv("CLIENT_HOST")

type ContextValue struct{}

func (s *Server) RegisterRoutes() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger, middleware.Recoverer, httprate.LimitByIP(64, 2*time.Minute))
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
	}))

	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", s.AuthRoute(s.RegisterHandler))
		r.Post("/login", s.AuthRoute(s.LoginHandler))
		r.Get("/logout", s.LogoutHandler)
	})

	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/session", s.SecuredRoute(s.GetSession))
		r.Get("/link/{imei}", s.SecuredRoute(s.LinkDevice))
		r.Get("/devices", s.SecuredRoute(s.GetDevices))
	})

	return r
}

func (s *Server) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:   "token",
		MaxAge: -1,
		Path:   "/",
	})

	http.Redirect(w, r, clientHost, http.StatusOK)
}

func (s *Server) RegisterHandler(w http.ResponseWriter, r *http.Request, data AuthData) {
	hash, err := argon2id.CreateHash(data.Password, &argon2id.Params{Parallelism: 1, SaltLength: 16, KeyLength: 16, Iterations: 4})
	if err != nil {
		HttpError(w, "failed to hash a password", http.StatusInternalServerError, err)
		return
	}

	id := ulid.Make().String()
	if err := s.Queries.CreateUser(context.Background(), database.CreateUserParams{
		ID:       id,
		Email:    data.Email,
		Password: hash,
	}); err != nil {
		HttpError(w, "email already exists", http.StatusBadRequest, err)
		return
	}

	signed := auth.SignJwt(id)
	auth.SetAuthCookie(w, signed)
}

func (s *Server) LoginHandler(w http.ResponseWriter, r *http.Request, data AuthData) {
	user, err := s.Queries.GetUserByEmail(context.Background(), data.Email)
	if err != nil {
		HttpError(w, "wrong email or password", http.StatusUnauthorized, err)
		return
	}

	if match, err := argon2id.ComparePasswordAndHash(data.Password, user.Password); !match || err != nil {
		HttpError(w, "wrong email or password", http.StatusUnauthorized, err)
		return
	}

	signed := auth.SignJwt(user.ID)
	auth.SetAuthCookie(w, signed)
}

func (s *Server) GetSession(w http.ResponseWriter, r *http.Request, id string) {
	user, err := s.Queries.GetUser(context.Background(), id)
	if err != nil {
		HttpError(w, "user not found", http.StatusBadRequest, err)
		return
	}
	user.Password = ""

	if err := json.NewEncoder(w).Encode(user); err != nil {
		log.Fatalf("failed to encode user: %v\n", err)
	}
}

func (s *Server) LinkDevice(w http.ResponseWriter, r *http.Request, id string) {
	deviceImei := chi.URLParam(r, "imei")

	if err := s.Queries.LinkDevice(context.Background(), database.LinkDeviceParams{DeviceImei: deviceImei, Userid: id}); err != nil {
		HttpError(w, "device not found", http.StatusBadRequest, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *Server) GetDevices(w http.ResponseWriter, r *http.Request, id string) {
	devices, err := s.Queries.GetDevices(context.Background(), id)
	if err != nil {
		HttpError(w, "failed to get devices", http.StatusBadRequest, err)
		return
	}

	if err := json.NewEncoder(w).Encode(devices); err != nil {
		log.Fatalf("failed to encode devices: %v\n", err)
	}
}

// TODO: 0x14 sleep route; 0x48 restart; 0x48 shutdown; 0x61 light switch; 0x92 alarm on;0x93 alarm off
