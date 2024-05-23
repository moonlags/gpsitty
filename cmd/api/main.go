package main

import (
	"database/sql"
	"log"

	"gpsitty/internal/database"
	"gpsitty/internal/server"
	"gpsitty/internal/tcp"

	"github.com/joho/godotenv"

	_ "github.com/mattn/go-sqlite3"
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
	devices := make(map[string]*tcp.Device)
	queries := NewDB()

	server, err := server.NewServer(queries, devices)
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
	conn, err := sql.Open("sqlite3", "gpsitty.db")
	if err != nil {
		log.Fatalf("failed to open db: %v\n", err)
	}
	return database.New(conn)
}
