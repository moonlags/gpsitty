package server

import (
	"database/sql"
	"fmt"
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

func NewServer(conf *oauth2.Config, store *sessions.CookieStore, devices map[string]*tcp.Device) (*http.Server, error) {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", os.Getenv("DB_USERNAME"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"))

	conn, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	queries := database.New(conn)

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
