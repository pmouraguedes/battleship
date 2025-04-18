package server

import (
	"log"
	"net"
)

type Server struct {
	address string
	gm      *GameManager
}

func (s Server) handleConnection(conn net.Conn, connectionId int) {
	defer conn.Close()
	buf := make([]byte, 1024)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			log.Println(err)
			return
		}
		log.Printf("[%d] Received: %s", connectionId, string(buf[:n]))

		data := s.gm.Handle(string(buf[:n]), connectionId)

		_, err = conn.Write([]byte(data))
		if err != nil {
			log.Println(err)
			return
		}
	}
}

func (s Server) Start() {
	addr, err := net.ResolveTCPAddr("tcp", s.address)
	if err != nil {
		log.Fatal(err)
	}
	ln, err := net.ListenTCP("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close()
	log.Printf("Server started on %s", s.address)

	connectionId := 0

	// Only accept the first two connections
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal(err)
		}
		connectionId++
		go s.handleConnection(conn, connectionId)
	}
}

func NewServer(addressString string) Server {
	return Server{
		address: addressString,
		gm:      NewGameManager(),
	}
}
