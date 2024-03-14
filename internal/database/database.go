package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	_ "github.com/joho/godotenv/autoload"
)

type Service interface {
	Health() map[string]string
	GetUser(userID string) (*User, error)
	CreateUser(user User) error
	InsertPosition(latitude float32, longitude float32, speed uint8, heading uint16, imei string) error
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

func New() (Service, error) {
	connStr := fmt.Sprintf("%s:%s@(%s:%s)/%s", username, password, host, port, database)
	db, err := sqlx.Connect("mysql", connStr)
	if err != nil {
		return nil, err
	}
	s := &service{db: db}
	return s, nil
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

type User struct {
	ID      string   `json:"sub,omitempty" db:"id"`
	Name    string   `json:"name,omitempty"`
	Picture string   `json:"picture,omitempty" db:"picture"`
	Email   string   `json:"email,omitempty" db:"email"`
	Devices []string `db:"device_imei"`
}

func (s *service) GetUser(ID string) (*User, error) {
	query := `
		SELECT users.id, users.email, users.avatar, user_devices.device_imei
		FROM users
		INNER JOIN user_devices ON users.id = user_devices.userid
		WHERE users.id = $1;
	`

	var user User
	if err := s.db.Select(&user, query, ID); err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *service) CreateUser(user User) error {
	query := "INSERT INTO users (id,name,email,avatar) VALUES (:id,:name,:email,:avatar) ON DUPLICATE KEY UPDATE SET last_login_time = CURRENT_TIMESTAMP;"

	if _, err := s.db.NamedExec(query, map[string]interface{}{
		"id":     user.ID,
		"email":  user.Email,
		"avatar": user.Picture,
		"name":   user.Name,
	}); err != nil {
		return err
	}

	return nil
}

func (s *service) InsertPosition(latitude float32, longitude float32, speed uint8, heading uint16, imei string) error {
	query := `INSERT INTO positions (latitude,longitude,speed,heading,device_imei)
		VALUES (:latitude,:longitude,:speed,:heading,:device_imei);
	
	DELETE FROM positions
	WHERE id IN (
		SELECT id FROM (SELECT id FROM positions WHERE device_imei=:device_imei ORDER BY created_at ASC OFFSET 10) as subquery
	);`

	if _, err := s.db.NamedExec(query, map[string]interface{}{
		"latitude":    latitude,
		"longitude":   longitude,
		"speed":       speed,
		"heading":     heading,
		"device_imei": imei,
	}); err != nil {
		return err
	}

	return nil
}
