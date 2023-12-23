package server

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"

	"gpsitty/internal/database"

	_ "github.com/joho/godotenv/autoload"
)

type Server struct {
	port               int
	db                 database.Service
	device_connections map[string]net.Conn
}

func NewServer(device_connections map[string]net.Conn) *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))

	db, err := database.New()
	if err != nil {
		log.Fatal("failed to create database service:", err)
	}

	NewServer := &Server{
		port:               port,
		db:                 db,
		device_connections: device_connections,
	}

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
