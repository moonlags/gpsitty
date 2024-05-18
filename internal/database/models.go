// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0

package database

import (
	"time"
)

type Device struct {
	Imei             string
	BatteryPower     int16
	Charging         bool
	LastStatusPacket time.Time
}

type Position struct {
	ID         int32
	Latitude   float64
	Longitude  float64
	Speed      int16
	Heading    int16
	DeviceImei string
	CreatedAt  time.Time
}

type User struct {
	ID            string
	Name          string
	Email         string
	Avatar        string
	LastLoginTime time.Time
}

type UserDevice struct {
	Userid     string
	DeviceImei string
}