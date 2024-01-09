package tcp

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"

	"gpsitty/internal/database"

	"go.uber.org/zap"
)

type Server struct {
	port               int
	db                 database.Service
	listener           net.Listener
	device_connections map[string]net.Conn
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
	}
	return server
}

type Device struct {
	Connection net.Conn
	IMEI       string
	Logger     *zap.SugaredLogger
}

func (s *Server) Listen() {
	defer s.listener.Close()
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatal("failed to create zap logger:", err)
	}
	defer logger.Sync()
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			continue
		}
		device := Device{
			Connection: conn,
			Logger:     logger.Sugar(),
		}
		go device.handleConnection(s.device_connections)
	}
}

func (d *Device) handleConnection(device_connections map[string]net.Conn) {
	defer d.Connection.Close()
	d.Logger.Infof("new connection from %s\n", d.Connection.RemoteAddr().String())
	buffer := make([]byte, 512)
	for {
		n, err := d.Connection.Read(buffer)
		if err != nil {
			d.Logger.Errorf("failed to read from %s: %s\n", d.Connection.RemoteAddr().String(), err.Error())
			break
		}
		d.Logger.Infof("%v from %s\n", buffer[:n], d.Connection.RemoteAddr().String())
		packet, err := parsePacket(d, device_connections, buffer[:n])
		if err != nil {
			d.Logger.Errorf("failed to parse packet from %s: %s\n", d.Connection.RemoteAddr().String(), err.Error())
			continue
		}
		response, err := packet.Process(d, device_connections)
		if err != nil {
			d.Logger.Errorf("failed to process packet from %s: %s\n", d.Connection.RemoteAddr().String(), err.Error())
			continue
		}
		if _, err := d.Connection.Write(response); err != nil {
			d.Logger.Errorf("failed to write to %s: %s\n", d.Connection.RemoteAddr().String(), err.Error())
			break
		}
	}
	delete(device_connections, d.IMEI)
}
