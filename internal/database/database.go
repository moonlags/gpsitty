package database

import (
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Service struct {
	DB *sqlx.DB
}

func New() (*Service, error) {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", os.Getenv("DB_USERNAME"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"))

	database, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		return nil, err
	}
	s := &Service{DB: database}
	return s, nil
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
	if err := s.DB.Select(&user, query, ID); err != nil {
		return nil, err
	}

	query = `
		SELECT device_imei
		FROM user_devices
		WHERE userid = $1;
	`
	if err := s.DB.Select(user.Devices, query, ID); err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *Service) CreateUser(user User) error {
	query := "INSERT INTO users (id,name,email,avatar) VALUES (:id,:name,:email,:avatar) ON CONFLICT(id) DO UPDATE SET last_login_time = CURRENT_TIMESTAMP;"

	_, err := s.DB.NamedExec(query, map[string]interface{}{
		"id":     user.ID,
		"email":  user.Email,
		"avatar": user.Picture,
		"name":   user.Name,
	})
	return err
}

func (s *Service) InsertPosition(latitude float32, longitude float32, speed uint8, heading uint16, imei string) error {
	query := `INSERT INTO positions (latitude,longitude,speed,heading,device_imei)
		VALUES (:latitude,:longitude,:speed,:heading,:device_imei);
	
	DELETE FROM positions
	WHERE id IN (
		SELECT id FROM FROM positions WHERE device_imei=:device_imei ORDER BY created_at DESC OFFSET 10)
	);`

	_, err := s.DB.NamedExec(query, map[string]interface{}{
		"latitude":    latitude,
		"longitude":   longitude,
		"speed":       speed,
		"heading":     heading,
		"device_imei": imei,
	})
	return err
}

func (s *Service) UpdateBatteryPower(imei string, batteryPower byte) error {
	query := `UPDATE SET battery_power = $1 WHERE imei = $2;`

	_, err := s.DB.Exec(query, batteryPower, imei)
	return err
}

type Device struct {
	IMEI             string `json:"imei,omitempty" db:"imei"`
	BatteryPower     byte   `json:"battery_power,omitempty" db:"battery_power"`
	IsCharging       bool   `json:"is_charging,omitempty" db:"charging"`
	LastStatusPacket int64  `json:"last_status_packet,omitempty" db:"last_status_packet"`
}

func (s *Service) InsertDevice(device Device) error {
	query := `INSERT INTO devices (imei, battery_power, charging) VALUES (:imei,:battery_power,:charging);`

	_, err := s.DB.NamedExec(query, device)
	return err
}

func (s *Service) LinkDevice(imei string, userID string) error {
	query := `INSERT INTO user_devices (userid, device_imei) VALUES ($1, $2);`

	_, err := s.DB.Exec(query, userID, imei)
	return err
}

func (s *Service) Charging(charging bool, imei string) error {
	query := `UPDATE devices SET charging = $1 WHERE imei = $2;`

	_, err := s.DB.Exec(query, charging, imei)
	return err
}
