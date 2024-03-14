package tcp

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"gpsitty/internal/database"
)

type Packet interface {
	Process(device *Device, server *Server) ([]byte, error)
}

func (d *Device) parsePacket(packet []byte) (Packet, error) {
	if !bytes.HasPrefix(packet, []byte{0x78, 0x78}) || !bytes.HasSuffix(packet, []byte{0x0d, 0x0a}) {
		return nil, errors.New("invalid packet format.")
	}
	packetLenght, protocolNumber := int(packet[2]), packet[3]
	switch protocolNumber {
	default:
		return nil, errors.New("not supported protocol number.")
	case 1:
		if packetLenght != 13 {
			return nil, errors.New("invalid packet lenght.")
		}
		packetStruct := &LoginPacket{hex.EncodeToString(packet[4:12])}
		return packetStruct, nil
	case 8:
		if packetLenght != 1 {
			return nil, errors.New("invalid packet length.")
		}
		packetStruct := &HeartBeatPacket{}
		return packetStruct, nil
	case 0x10, 0x11:
		if packetLenght != 0x15 {
			return nil, errors.New("invalid packet length.")
		}
		var latitude, longitude uint32
		if err := binary.Read(bytes.NewReader(packet[11:15]), binary.BigEndian, &latitude); err != nil {
			return nil, err
		} else if err := binary.Read(bytes.NewReader(packet[15:19]), binary.BigEndian, &longitude); err != nil {
			return nil, err
		}
		packetStruct := &PosititioningPacket{
			Latitude:  float32(latitude) / (30000.0 * 60.0),
			Longitude: float32(longitude) / (30000.0 * 60.0),
			Speed:     packet[19],
			Heading:   uint16(packet[21]) | (uint16(packet[20]&3) << 8),
			Timestamp: time.Now().Unix(),
		}
		return packetStruct, nil
	case 0x13:
		if packetLenght != 6 && packetLenght != 7 {
			return nil, errors.New("invalid packet length.")
		}

		packetStruct := &StatusPacket{
			BatteryPower: packet[4],
		}

		return packetStruct, nil
	}
}

type LoginPacket struct {
	IMEI string
}

func (p *LoginPacket) Process(device *Device, server *Server) ([]byte, error) {
	if _, ok := server.DeviceConnections[p.IMEI]; ok || device.IMEI != "" {
		return nil, errors.New("device already logged in.")
	}
	fmt.Printf("INFO: device %s logged in\n", p.IMEI)

	if err := server.DB.InsertDevice(database.Device{IMEI: device.IMEI}); err != nil {
		return nil, err
	}

	device.IMEI = p.IMEI
	server.DeviceConnections[p.IMEI] = device.Connection

	return []byte{0x78, 0x78, 1, 1, 0x0d, 0x0a}, nil
}

type PosititioningPacket struct {
	Latitude       float32
	Longitude      float32
	Speed          uint8
	Heading        uint16
	Timestamp      int64
	ProtocolNumber uint8
}

func (p *PosititioningPacket) Process(device *Device, server *Server) ([]byte, error) {
	if device.IMEI == "" {
		return nil, errors.New("device is not logged in.")
	}

	if err := server.DB.InsertPosition(p.Latitude, p.Longitude, p.Speed, p.Heading, device.IMEI); err != nil {
		return nil, err
	}

	now := time.Now()
	return []byte{0x78, 0x78, 7, p.ProtocolNumber, byte(now.Year() - 2000), byte(now.Month()), byte(now.Day()), byte(now.Hour()), byte(now.Minute()), byte(now.Second()), 0x0d, 0x0a}, nil
}

type HeartBeatPacket struct{}

func (p *HeartBeatPacket) Process(device *Device, server *Server) ([]byte, error) {
	if device.IMEI == "" {
		return nil, errors.New("device is not logged in.")
	}

	return []byte{0x78, 0x78, 1, 8, 0x0d, 0x0a}, nil
}

type StatusPacket struct {
	BatteryPower byte
}

func (p *StatusPacket) Process(device *Device, server *Server) ([]byte, error) {
	if device.IMEI == "" {
		return nil, errors.New("device is not logged in.")
	}

	if err := server.DB.UpdateBatteryPower(device.IMEI, p.BatteryPower); err != nil {
		return nil, err
	}

	var cooldown byte = 1
	if p.BatteryPower <= 15 {
		cooldown = 5
	}

	return []byte{0x78, 0x78, 2, 0x13, cooldown, 0x0d, 0x0a}, nil
}
