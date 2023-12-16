package tcp

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"

	"gpsitty/internal/database"
)

// type DeviceData struct {
// 	IMEI                 string
// 	Posititions          []PosititioningPacket
// 	BatteryPower         uint8
// 	LastStatusPacketTime int64
// 	StatusCooldown       int32
// 	IsLoggedIn           bool
// 	InChargingState      bool
// 	Connection           *net.Conn
// }

type Server struct {
	port     int
	db       database.Service
	listener net.Listener
}

func NewServer() *Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatal("failed to create tcp listener", err)
	}

	server := &Server{
		port:     port,
		db:       database.New(),
		listener: listener,
	}

	return server
}

func (s *Server) Listen() {
	defer s.listener.Close()

	for {
		conn, err := s.listener.Accept()
		if err != nil {
			continue
		}

		go s.handleConnection(conn)
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()

	buffer := make([]byte, 512)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			break
		}

	}
}
