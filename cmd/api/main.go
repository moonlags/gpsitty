package main

import (
	"log"

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
}

func main() {
	googleConf := auth.NewGoogleConfig("http://localhost:50731/auth/google/callback")
	store := auth.NewCookieStore(maxAge, isProd)

	devices := make(map[string]*tcp.Device)

	server, err := server.NewServer(googleConf, store, devices)
	if err != nil {
		log.Fatalf("failed to create http server: %v\n", err)
	}

	serverTcp, err := tcp.NewServer(devices)
	if err != nil {
		log.Fatalf("failed to create tcp server: %v\n", err)
	}

	go func() {
		log.Fatal(serverTcp.Listen())
	}()

	log.Fatal(server.ListenAndServe())
}
