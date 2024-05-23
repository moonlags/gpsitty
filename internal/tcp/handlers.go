package tcp

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"time"

	"gpsitty/internal/database"
)

func (d *Device) ParsePacket(buffer []byte, server *Server) ([]byte, error) {
	if !bytes.HasPrefix(buffer, []byte{0x78, 0x78}) || !bytes.HasSuffix(buffer, []byte{0x0d, 0x0a}) {
		return nil, errors.New("invalid packet format")
	}

	packetLenght, protocolNumber := int(buffer[2]), buffer[3]

	switch protocolNumber {
	default:
		return nil, errors.New("not supported protocol number")
	case 1:
		if packetLenght != 13 {
			return nil, errors.New("invalid packet lenght")
		}

		packet := &LoginPacket{hex.EncodeToString(buffer[4:12])}
		return packet.Process(d, server)
	case 8:
		if packetLenght != 1 {
			return nil, errors.New("invalid packet length")
		}

		return []byte{0x78, 0x78, 1, 8, 0x0d, 0x0a}, nil
	case 0x10, 0x11:
		if packetLenght != 0x15 {
			return nil, errors.New("invalid packet length")
		}

		var latitude, longitude uint32
		if err := binary.Read(bytes.NewReader(buffer[11:15]), binary.BigEndian, &latitude); err != nil {
			return nil, err
		} else if err := binary.Read(bytes.NewReader(buffer[15:19]), binary.BigEndian, &longitude); err != nil {
			return nil, err
		}

		packet := &PosititioningPacket{
			Latitude:  float32(latitude) / (30000.0 * 60.0),
			Longitude: float32(longitude) / (30000.0 * 60.0),
			Speed:     buffer[19],
			Heading:   uint16(buffer[21]) | (uint16(buffer[20]&3) << 8),
		}

		return packet.Process(d, server.Queries)
	case 0x13:
		if packetLenght != 6 && packetLenght != 7 {
			return nil, errors.New("invalid packet length")
		}

		packet := &StatusPacket{
			BatteryPower: buffer[4],
		}

		return packet.Process(d, server.Queries)
	case 0x30:
		if packetLenght != 1 {
			return nil, errors.New("invalid packet length")
		}

		now := time.Now()
		return []byte{0x78, 0x78, 7, 0x30, byte(now.Year() - 2000), byte(now.Month()), byte(now.Day()), byte(now.Hour()), byte(now.Minute()), byte(now.Second()), 0x0d, 0x0a}, nil
	case 0x81, 0x83:
		if packetLenght != 1 {
			return nil, errors.New("invalid packet length")
		} else if d.IMEI == "" {
			return nil, errors.New("device is not logged in")
		}

		if err := server.Queries.UpdateCharging(context.Background(), database.UpdateChargingParams{Charging: false, Imei: d.IMEI}); err != nil {
			return nil, err
		}

		return []byte{0x78, 0x78, 1, protocolNumber, 0x0d, 0x0a}, nil
	case 0x82:
		if packetLenght != 1 {
			return nil, errors.New("invalid packet length")
		} else if d.IMEI == "" {
			return nil, errors.New("device is not logged in")
		}

		if err := server.Queries.UpdateCharging(context.Background(), database.UpdateChargingParams{Charging: true, Imei: d.IMEI}); err != nil {
			return nil, err
		}

		return []byte{0x78, 0x78, 1, 0x82, 0x0d, 0x0a}, nil
	case 26, 27:
		if d.IMEI == "" {
			return nil, errors.New("device is not logged in")
		}

		now := time.Now()
		return []byte{0x78, 0x78, 7, 0x30, hexToDec(now.Year() - 2000), hexToDec(int(now.Month())), hexToDec(now.Day()), hexToDec(now.Hour()), hexToDec(now.Minute()), hexToDec(now.Second()), 0x0d, 0x0a}, nil
	case 0x15:
		return []byte{0x78, 0x78, 1, 0x15, 0x0d, 0x0a}, nil
	}
}

func hexToDec(hex_byte int) byte {
	return byte(hex_byte%10 + hex_byte/10*16)
}

type LoginPacket struct {
	IMEI string
}

func (p *LoginPacket) Process(device *Device, server *Server) ([]byte, error) {
	if device.IMEI != "" {
		return nil, errors.New("device already logged in")
	}

	if err := server.Queries.InsertDevice(context.Background(), database.InsertDeviceParams{Imei: p.IMEI}); err != nil {
		return nil, err
	}

	device.IMEI = p.IMEI
	server.Devices[p.IMEI] = device

	return []byte{0x78, 0x78, 1, 1, 0x0d, 0x0a}, nil
}

type PosititioningPacket struct {
	Latitude       float32
	Longitude      float32
	Speed          uint8
	Heading        uint16
	ProtocolNumber uint8
}

func (p *PosititioningPacket) Process(device *Device, queries *database.Queries) ([]byte, error) {
	if device.IMEI == "" {
		return nil, errors.New("device is not logged in")
	}

	if err := queries.CreatePosition(context.Background(), database.CreatePositionParams{
		Latitude: float64(p.Latitude), Longitude: float64(p.Longitude), Speed: int64(p.Speed), Heading: int64(p.Heading), DeviceImei: device.IMEI,
	}); err != nil {
		return nil, err
	}

	now := time.Now()
	return []byte{0x78, 0x78, 7, p.ProtocolNumber, byte(now.Year() - 2000), byte(now.Month()), byte(now.Day()), byte(now.Hour()), byte(now.Minute()), byte(now.Second()), 0x0d, 0x0a}, nil
}

type StatusPacket struct {
	BatteryPower byte
}

func (p *StatusPacket) Process(device *Device, queries *database.Queries) ([]byte, error) {
	if device.IMEI == "" {
		return nil, errors.New("device is not logged in")
	}

	if err := queries.UpdateBatteryPower(context.Background(), database.UpdateBatteryPowerParams{Imei: device.IMEI, BatteryPower: int64(p.BatteryPower)}); err != nil {
		return nil, err
	}

	var cooldown byte = 1
	if p.BatteryPower <= 15 {
		cooldown = 5
	}

	return []byte{0x78, 0x78, 2, 0x13, cooldown, 0x0d, 0x0a}, nil
}
