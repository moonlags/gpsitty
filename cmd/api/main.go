package main

import (
	"log"
	"net"
	"os"

	"gpsitty/internal/auth"
	"gpsitty/internal/server"
	"gpsitty/internal/tcp"

	"github.com/joho/godotenv"
)

const (
	maxAge = 3600 * 3
	isProd = false
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("failed to load .env file", err)
	}

	file, err := os.Create("/logs/logs.log")
	if err != nil {
		log.Fatalf("failed to create logs file: %v\n", err)
	}

	log.SetOutput(file)
}

func main() {
	deviceConnections := make(map[string]net.Conn)
	googleConf := auth.NewGoogleConfig("http://localhost:50731/auth/google/callback")
	store := auth.NewCookieStore(maxAge, isProd)

	server := server.NewServer(googleConf, store, deviceConnections)
	serverTcp := tcp.NewServer(deviceConnections)

	go serverTcp.Listen()

	err := server.ListenAndServe()
	if err != nil {
		log.Fatal("failed to start server:", err)
	}
}
