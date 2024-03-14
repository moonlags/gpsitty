package tcp

import (
	"fmt"
	"log"
	"net"
	"os"

	"gpsitty/internal/database"
)

type Server struct {
	DB                database.Service
	DeviceConnections map[string]net.Conn
}

func NewServer(deviceConnections map[string]net.Conn) *Server {
	db, err := database.New()
	if err != nil {
		log.Fatal("FATAL: failed to create database service:", err)
	}

	server := &Server{
		DB:                db,
		DeviceConnections: deviceConnections,
	}

	return server
}

type Device struct {
	Connection net.Conn
	IMEI       string
}

func (s *Server) Listen() {
	listener, err := net.Listen("tcp", ":"+os.Getenv("TCP_PORT"))
	if err != nil {
		log.Fatal("FATAL: failed to create tcp listener", err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
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
	fmt.Printf("INFO: new connection from %s\n", d.Connection.RemoteAddr().String())

	buffer := make([]byte, 512)
	for {
		n, err := d.Connection.Read(buffer)
		if err != nil {
			fmt.Printf("ERROR: failed to read from %s: %s\n", d.Connection.RemoteAddr().String(), err.Error())
			break
		}

		fmt.Printf("INFO: %v from %s\n", buffer[:n], d.Connection.RemoteAddr().String())

		packet, err := d.parsePacket(server.DeviceConnections, buffer[:n])
		if err != nil {
			fmt.Printf("ERROR: failed to parse packet from %s: %s\n", d.Connection.RemoteAddr().String(), err.Error())
			continue
		}

		response, err := packet.Process(d, server)
		if err != nil {
			fmt.Printf("ERROR: failed to process packet from %s: %s\n", d.Connection.RemoteAddr().String(), err.Error())
			continue

		}
		if _, err := d.Connection.Write(response); err != nil {
			fmt.Printf("ERROR: failed to write to %s: %s\n", d.Connection.RemoteAddr().String(), err.Error())
			break
		}
	}
	delete(server.DeviceConnections, d.IMEI)
}
