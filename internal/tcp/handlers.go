package tcp

import (
	"bytes"
	"encoding/hex"
	"errors"
	"log"
	"net"
)

type PosititioningPacket struct {
	Latitude  float32
	Longitude float32
	Speed     uint8
	Heading   uint16
	Timestamp int64
}

func (s *Server) parsePacket(packet []byte, remoteAddr net.Addr) ([]byte, error) {
	log.Printf("got packet %v from %v\n", packet, remoteAddr)

	if !bytes.HasPrefix(packet, []byte{0x78, 0x78}) || !bytes.HasSuffix(packet, []byte{0x0d, 0x0a}) {
		return nil, errors.New("invalid packet format")
	}

	packetLenght := int(packet[2])
	protocolNumber := packet[3]

	switch protocolNumber {
	default:
		return nil, errors.New("not supported protocol number")
	case 1:
		if packetLenght != 13 {
			return nil, errors.New("invalid packet lenght")
		}

		packetStruct := LoginPacket{hex.EncodeToString(packet[4:12])}

		return s.ProcessLogin(packetStruct), nil
	}
}

type LoginPacket struct {
	IMEI string
}

func (s *Server) ProcessLogin(packet LoginPacket) []byte {
	log.Printf("device with imei %v logged in", packet.IMEI)

	// todo process properly

	return []byte{0x78, 0x78, 1, 1, 0x0d, 0x0a}
}
