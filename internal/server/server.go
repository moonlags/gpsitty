package server

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"

	"gpsitty/internal/database"

	_ "github.com/joho/godotenv/autoload"
	"github.com/markbates/goth"
)

type Server struct {
	port    int
	db      database.Service
	users   map[string]goth.User
	devices map[string]net.Conn
}

func NewServer() *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	NewServer := &Server{
		port:    port,
		db:      database.New(),
		users:   make(map[string]goth.User),
		devices: make(map[string]net.Conn),
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
