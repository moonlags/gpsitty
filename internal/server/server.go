package server

import (
	"net/http"
	"os"
	"time"

	"gpsitty/internal/database"
	"gpsitty/internal/tcp"
)

type Server struct {
	Queries *database.Queries
	Devices map[string]*tcp.Device
}

func NewServer(queries *database.Queries, devices map[string]*tcp.Device) (*http.Server, error) {
	NewServer := &Server{
		Queries: queries,
		Devices: devices,
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
