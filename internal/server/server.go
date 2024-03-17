package server

import (
	"net/http"
	"os"
	"time"

	"gpsitty/internal/database"
	"gpsitty/internal/tcp"

	"github.com/gorilla/sessions"
	_ "github.com/joho/godotenv/autoload"
	"golang.org/x/oauth2"
)

type Server struct {
	DB      *database.Service
	Devices map[string]*tcp.Device
	Conf    *oauth2.Config
	Store   *sessions.CookieStore
}

func NewServer(conf *oauth2.Config, store *sessions.CookieStore, devices map[string]*tcp.Device) (*http.Server, error) {
	database, err := database.New()
	if err != nil {
		return nil, err
	}

	NewServer := &Server{
		DB:      database,
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
