package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"gpsitty/internal/auth"
	"gpsitty/internal/database"
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
	queries := NewDB()

	server, err := server.NewServer(queries, googleConf, store, devices)
	if err != nil {
		log.Fatalf("failed to create http server: %v\n", err)
	}

	serverTcp, err := tcp.NewServer(queries, devices)
	if err != nil {
		log.Fatalf("failed to create tcp server: %v\n", err)
	}

	go func() {
		log.Fatal(serverTcp.Listen())
	}()

	log.Fatal(server.ListenAndServe())
}

func NewDB() *database.Queries {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", os.Getenv("DB_USERNAME"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"))

	conn, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("failed to open db: %v\n", err)
	}
	return database.New(conn)
}
