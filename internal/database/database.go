package database

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type Service struct {
	db *sqlx.DB
}

func New() (*Service, error) {
	db, err := sqlx.Connect("sqlite3", "gpsitty.db")
	if err != nil {
		return nil, err
	}
	s := &Service{db: db}
	return s, nil
}

func (s *Service) Health() map[string]string {
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
	Name    string   `json:"name,omitempty" db:"name"`
	Picture string   `json:"picture,omitempty" db:"picture"`
	Email   string   `json:"email,omitempty" db:"email"`
	Devices []string `json:"devices,omitempty"`
}

func (s *Service) GetUser(ID string) (*User, error) {
	query := `
		SELECT id, name, email, avatar
		FROM users
		WHERE id = $1;
	`

	var user User
	if err := s.db.Select(&user, query, ID); err != nil {
		return nil, err
	}

	query = `
		SELECT device_imei
		FROM user_devices
		WHERE userid = $1;
	`
	if err := s.db.Select(user.Devices, query, ID); err != nil {
		return nil, err
	}

	fmt.Printf("get user query: %v\n", user)

	return &user, nil
}

func (s *Service) CreateUser(user User) error {
	query := "INSERT INTO users (id,name,email,avatar) VALUES (:id,:name,:email,:avatar) ON CONFLICT(id) DO UPDATE SET last_login_time = CURRENT_TIMESTAMP;"

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

func (s *Service) InsertPosition(latitude float32, longitude float32, speed uint8, heading uint16, imei string) error {
	query := `INSERT INTO positions (latitude,longitude,speed,heading,device_imei)
		VALUES (:latitude,:longitude,:speed,:heading,:device_imei);
	
	DELETE FROM positions
	WHERE id IN (
		SELECT id FROM (SELECT id FROM positions WHERE device_imei=:device_imei ORDER BY created_at DESC LIMIT -1 OFFSET 10)
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

func (s *Service) UpdateBatteryPower(imei string, batteryPower byte) error {
	query := `UPDATE SET battery_power = $1 WHERE imei = $2;`

	if _, err := s.db.Query(query, batteryPower, imei); err != nil {
		return err
	}

	return nil
}

type Device struct {
	Connection       net.Conn `json:"connection,omitempty"`
	IMEI             string   `json:"imei,omitempty" db:"imei"`
	BatteryPower     byte     `json:"battery_power,omitempty" db:"battery_power"`
	IsCharging       byte     `json:"is_charging,omitempty" db:"charging"`
	LastStatusPacket int64    `json:"last_status_packet,omitempty" db:"last_status_packet"`
}

func (s *Service) InsertDevice(device Device) error {
	query := `INSERT INTO devices (imei, battery_power, charging) VALUES (:imei,:battery_power,:charging);`

	if _, err := s.db.NamedExec(query, device); err != nil {
		return err
	}
	return nil
}

//! get device
