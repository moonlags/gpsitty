package server

import (
	"net/http"
	"os"
	"time"

	"gpsitty/internal/database"
	"gpsitty/internal/tcp"

	"github.com/gorilla/sessions"
	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"
	"golang.org/x/oauth2"
)

type Server struct {
	Queries *database.Queries
	Devices map[string]*tcp.Device
	Conf    *oauth2.Config
	Store   *sessions.CookieStore
}

func NewServer(queries *database.Queries, conf *oauth2.Config, store *sessions.CookieStore, devices map[string]*tcp.Device) (*http.Server, error) {
	NewServer := &Server{
		Queries: queries,
		Devices: devices,
		Conf:    conf,
		Store:   store,
	}

	// Declare Server config
	server := &http.Server{
		Addr:         ":" + os.Getenv("PORT"),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server, nil
}
