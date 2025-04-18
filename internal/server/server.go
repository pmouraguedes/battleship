package server

import (
	"log"
	"net"
)

type Server struct {
	address string
	gm      *GameManager
}

// func handleMessage(msg string, connectionId int) string {
// 	// if HELLO command
// 	// message: HELLO <player_name>
// 	if strings.HasPrefix(msg, "HELLO") {
// 		return handleHelloCommand(msg, connectionId)
// 	}
//
// 	// if SHIP command
// 	// message: SHIP <ship_type> <x> <y> <direction>
// 	if strings.HasPrefix(msg, "SHIP") {
// 		return handleShipCommand(msg, connectionId)
// 	}
//
// 	return "Invalid command\n"
// }

// func handleHelloCommand(msg string, connectionId int) string {
// 	// assert that the message is a HELLO command
// 	if !strings.HasPrefix(msg, "HELLO") {
// 		panic("Should be a HELLO command")
// 	}
//
// 	parts := strings.Fields(msg)
// 	if len(parts) != 2 {
// 		return "Invalid HELLO command\n"
// 	}
// 	playerName := parts[1]
// 	if len(playerName) < 3 || len(playerName) > 20 {
// 		return "Invalid player name\n"
// 	}
// 	playerCode := "P" + fmt.Sprintf("%d", connectionId)
//
// 	switch connectionId {
// 	case 1:
// 		return fmt.Sprintf("WELCOME %s %s\n", playerCode, playerName)
// 	case 2:
// 		opponentName := "Opponent"
// 		return fmt.Sprintf("WELCOME %s %s %s\n", playerCode, playerName, opponentName)
// 	default:
// 		panic("Only two players are allowed")
// 	}
// }

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
