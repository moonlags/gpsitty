package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/joho/godotenv/autoload"
	_ "github.com/mattn/go-sqlite3"
)

type Service interface {
	Health() map[string]string
	GetUserWithDevices(userID string) UserWithDevices
}

type service struct {
	db *sql.DB
}

func New() Service {
	db, err := sql.Open("sqlite3", "./internal/database/gpsitty.db?cache=shared&mode=rwc&_journal_mode=WAL&busy_timeout=10000&_fk=1")
	if err != nil {
		log.Fatal(err)
	}
	s := &service{db: db}
	return s
}

func (s *service) Health() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	err := s.db.PingContext(ctx)
	if err != nil {
		log.Fatalf(fmt.Sprintf("db down: %v", err))
	}

	return map[string]string{
		"message": "It's healthy",
	}
}

type UserWithDevices struct {
	UserID    string
	Email     string
	AvatarURL string
	Devices   []struct {
		IMEI             string
		BatteryPower     int
		LastStatusPacket time.Time
		StatusCooldown   int
		Charging         bool
	}
}

func (s *service) GetUserWithDevices(userID string) UserWithDevices {
}
