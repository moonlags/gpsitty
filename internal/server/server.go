package server

import (
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"gpsitty/internal/database"

	"github.com/gorilla/sessions"
	_ "github.com/joho/godotenv/autoload"
	"golang.org/x/oauth2"
)

type Server struct {
	DB                database.Service
	DeviceConnections map[string]net.Conn
	Conf              *oauth2.Config
	Store             *sessions.CookieStore
}

func NewServer(conf *oauth2.Config, store *sessions.CookieStore, deviceConnections map[string]net.Conn) *http.Server {
	db, err := database.New()
	if err != nil {
		log.Fatal("failed to create database service:", err)
	}

	NewServer := &Server{
		DB:                db,
		DeviceConnections: deviceConnections,
		Conf:              conf,
		Store:             store,
	}

	// Declare Server config
	server := &http.Server{
		Addr:         ":" + os.Getenv("PORT"),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
