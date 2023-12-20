package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"
)

type Service interface {
	Health() map[string]string
	GetUserWithDevices(userID string) (*UserWithDevices, error)
	CreateUser(userID string, email string, avatar string) error
}

type service struct {
	db *sqlx.DB
}

var (
	database = os.Getenv("DB_DATABASE")
	password = os.Getenv("DB_PASSWORD")
	username = os.Getenv("DB_USERNAME")
	port     = os.Getenv("DB_PORT")
	host     = os.Getenv("DB_HOST")
)

func New() Service {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", username, password, host, port, database)
	db, err := sqlx.Connect("postgres", connStr)
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
	UserID    string   `db:"id"`
	Email     string   `db:"email"`
	AvatarURL string   `db:"avatar"`
	Devices   []string `db:"device_imei"`
}

func (s *service) GetUserWithDevices(userID string) (*UserWithDevices, error) {
	query := `
		SELECT users.id, users.email, users.avatar, user_devices.device_imei
		FROM users
		INNER JOIN user_devices ON users.id = user_devices.userid
		WHERE users.id = $1
	`

	var user UserWithDevices
	if err := s.db.Select(&user, query, userID); err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *service) CreateUser(userID string, email string, avatar string) error {
	query := "INSERT INTO users (id,email,avatar) VALUES (:id,:email,:avatar) ON CONFLICT (id) DO UPDATE SET last_login_time = CURRENT_TIMESTAMP"

	if _, err := s.db.NamedExec(query, map[string]interface{}{
		"id":     userID,
		"email":  email,
		"avatar": avatar,
	}); err != nil {
		return err
	}

	return nil
}
