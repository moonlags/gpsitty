package tcp

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"

	"gpsitty/internal/database"
)

type Server struct {
	port               int
	db                 database.Service
	listener           net.Listener
	device_connections map[string]net.Conn
	logged_connections map[net.Conn]string
}

func NewServer(device_connections map[string]net.Conn) *Server {
	port, _ := strconv.Atoi(os.Getenv("TCP_PORT"))

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatal("failed to create tcp listener", err)
	}

	db, err := database.New()
	if err != nil {
		log.Fatal("failed to create database service:", err)
	}

	server := &Server{
		port:               port,
		db:                 db,
		listener:           listener,
		device_connections: device_connections,
		logged_connections: make(map[net.Conn]string),
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

		response, err := s.parsePacket(buffer[:n], conn)
		if err != nil {
			log.Println("failed to parse pakcet ERROR:", err)
			continue
		}

		if _, err := conn.Write(response); err != nil {
			break
		}
	}
}
