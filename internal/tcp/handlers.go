package tcp

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"net"
	"time"
)

type Packet interface {
	Process(device *Device, server *Server) ([]byte, error)
}

func (d *Device) parsePacket(device_connections map[string]net.Conn, packet []byte) (Packet, error) {
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

	// if err :=

	now := time.Now()
}

func (s *Server) ProcessPositioning(packet PosititioningPacket, protocolNumber uint8, conn net.Conn) ([]byte, error) {
	fmt.Printf("new positioning packet %v", packet)

	if err := s.db.InsertPosition(packet.Latitude, packet.Longitude, packet.Speed, packet.Heading, s.logged_connections[conn]); err != nil {
		return nil, err
	}

	timeNow := time.Now()

	return []byte{
		0x78, 0x78, 7, protocolNumber, byte(timeNow.Year() - 2000), byte(timeNow.Month()), byte(timeNow.Day()), byte(timeNow.Hour()),
		byte(timeNow.Minute()), byte(timeNow.Second()), 0x0d, 0x0a,
	}, nil
}
