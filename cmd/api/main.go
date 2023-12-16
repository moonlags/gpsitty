package main

import (
	"gpsitty/internal/auth"
	"gpsitty/internal/server"
	"gpsitty/internal/tcp"
)

func main() {
	auth.New()

	server := server.NewServer()
	serverTcp := tcp.NewServer()

	go serverTcp.Listen()

	err := server.ListenAndServe()
	if err != nil {
		panic("cannot start server")
	}
}
