package tcp

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"log"
	"net"
	"time"
)

func (s *Server) parsePacket(packet []byte, conn net.Conn) ([]byte, error) {
	log.Printf("got packet %v from %v\n", packet, conn.RemoteAddr().String())

	if !bytes.HasPrefix(packet, []byte{0x78, 0x78}) || !bytes.HasSuffix(packet, []byte{0x0d, 0x0a}) {
		return nil, errors.New("invalid packet format.")
	}

	packetLenght := int(packet[2])
	protocolNumber := packet[3]

	switch protocolNumber {
	default:
		return nil, errors.New("not supported protocol number.")
	case 1:
		if packetLenght != 13 {
			return nil, errors.New("invalid packet lenght.")
		}

		packetStruct := LoginPacket{hex.EncodeToString(packet[4:12])}

		return s.ProcessLogin(packetStruct, conn), nil
	case 8:
		if packetLenght != 1 {
			return nil, errors.New("invalid packet length.")
		}

		return []byte{0x78, 0x78, 1, 8, 0x0d, 0x0a}, nil
	case 0x10, 0x11:
		if packetLenght != 0x15 {
			return nil, errors.New("invalid packet length.")
		} else if _, ok := s.logged_connections[conn]; !ok {
			return nil, errors.New("device is not logged.")
		}

		var latitude, longitude uint32
		if err := binary.Read(bytes.NewReader(packet[11:15]), binary.BigEndian, &latitude); err != nil {
			return nil, err
		} else if err := binary.Read(bytes.NewReader(packet[15:19]), binary.BigEndian, &longitude); err != nil {
			return nil, err
		}

		packetStruct := PosititioningPacket{
			Latitude:  float32(latitude) / (30000.0 * 60.0),
			Longitude: float32(longitude) / (30000.0 * 60.0),
			Speed:     packet[19],
			Heading:   uint16(packet[21]) | (uint16(packet[20]&3) << 8),
			Timestamp: time.Now().Unix(),
		}

		return s.ProcessPositioning(packetStruct, protocolNumber, conn)
	}
}

type LoginPacket struct {
	IMEI string
}

func (s *Server) ProcessLogin(packet LoginPacket, conn net.Conn) []byte {
	log.Printf("device with imei %v logged in", packet.IMEI)

	s.device_connections[packet.IMEI] = conn
	s.logged_connections[conn] = packet.IMEI

	return []byte{0x78, 0x78, 1, 1, 0x0d, 0x0a}
}

type PosititioningPacket struct {
	Latitude  float32
	Longitude float32
	Speed     uint8
	Heading   uint16
	Timestamp int64
}

func (s *Server) ProcessPositioning(packet PosititioningPacket, protocolNumber uint8, conn net.Conn) ([]byte, error) {
	log.Printf("new positioning packet %v", packet)

	if err := s.db.InsertPosition(packet.Latitude, packet.Longitude, packet.Speed, packet.Heading, s.logged_connections[conn]); err != nil {
		return nil, err
	}

	timeNow := time.Now()

	return []byte{
		0x78, 0x78, 7, protocolNumber, byte(timeNow.Year() - 2000), byte(timeNow.Month()), byte(timeNow.Day()), byte(timeNow.Hour()),
		byte(timeNow.Minute()), byte(timeNow.Second()), 0x0d, 0x0a,
	}, nil
}
