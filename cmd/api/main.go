package main

import (
	"log"
	"net"

	"gpsitty/internal/auth"
	"gpsitty/internal/server"
	"gpsitty/internal/tcp"
)

func main() {
	device_connections := make(map[string]net.Conn)
	auth.New()

	server := server.NewServer(device_connections)
	serverTcp := tcp.NewServer(device_connections)

	go serverTcp.Listen()

	err := server.ListenAndServe()
	if err != nil {
		log.Fatal("failed to start server:", err)
	}
}
