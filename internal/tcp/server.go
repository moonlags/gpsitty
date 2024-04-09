package tcp

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"os"

	"gpsitty/internal/database"
)

type Server struct {
	Queries *database.Queries
	Devices map[string]*Device
}

func NewServer(devices map[string]*Device) (*Server, error) {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", os.Getenv("DB_USERNAME"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"))

	conn, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	queries := database.New(conn)

	server := &Server{
		Queries: queries,
		Devices: devices,
	}

	return server, nil
}

type Device struct {
	Connection net.Conn
	IMEI       string
}

func (s *Server) Listen() error {
	listener, err := net.Listen("tcp", ":"+os.Getenv("TCP_PORT"))
	if err != nil {
		return err
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("failed to accept a connection: %v\n", err)
			continue
		}

		device := Device{
			Connection: conn,
		}

		go device.handleConnection(s)
	}
}

func (d *Device) handleConnection(server *Server) {
	defer d.Connection.Close()

	buffer := make([]byte, 256)
	for {
		n, err := d.Connection.Read(buffer)
		if err != nil {
			log.Printf("failed to read from connection: %v\n", err)
			break
		}

		log.Printf("packet from %v: %v\n", d.Connection.RemoteAddr(), buffer[:n])

		response, err := d.ParsePacket(buffer[:n], server)
		if err != nil {
			log.Printf("failed to parse packet: %v\n", err)
			continue
		}

		if _, err := d.Connection.Write(response); err != nil {
			log.Printf("failed to write to %v: %v\n", d.Connection.RemoteAddr(), err)
			break
		}
	}
	delete(server.Devices, d.IMEI)
}
