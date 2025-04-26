package server

import (
	"log"
	"net"
	"sync"

	"github.com/pmouraguedes/battleship/internal/game"
)

type Server struct {
	address string
	gm      *GameManager
	mu      sync.Mutex
}

// func (s *Server) handleConnection(conn net.Conn, connectionId int) {
// 	defer conn.Close()
//
// 	// s.gm.addConnection(conn, connectionId)
// 	// s.initializeGame(connectionId)
// 	// s.gm.handle(conn, connectionId)
//
// 	buf := make([]byte, 1024)
// 	for {
// 		n, err := conn.Read(buf)
// 		if err != nil {
// 			log.Println(err)
// 			return
// 		}
// 		log.Printf("[%d] Received: %s", connectionId, string(buf[:n]))
//
// 		data, err := s.gm.handle(conn, connectionId)
// 		if err != nil {
// 			log.Println(err)
// 		}
//
// 		log.Printf("[%d] Sending: %s...", connectionId, data)
// 		_, err = conn.Write([]byte(data))
// 		if err != nil {
// 			log.Println(err)
// 			return
// 		}
// 	}
// }

func (s *Server) Start() {
	addr, err := net.ResolveTCPAddr("tcp", s.address)
	if err != nil {
		log.Fatal(err)
	}
	ln, err := net.ListenTCP("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close()
	log.Printf("[server] server started on %s", s.address)

	connectionId := 0

	for {
		log.Printf("[server] waiting for connection...")
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal(err)
		}
		connectionId++
		log.Printf("[server %d] new connection from %s", connectionId, conn.RemoteAddr())

		s.gm.addConnection(conn, connectionId)
		s.initializeGame(connectionId)

		log.Printf("[server %d] handling connection...", connectionId)
		go s.gm.handle(conn, connectionId)
		// go s.handleConnection(conn, connectionId)
	}
}

func NewServer(addressString string) Server {
	return Server{
		address: addressString,
		gm:      newGameManager(),
	}
}

func (s *Server) initializeGame(connectionId int) {
	// s.gm.mu.Lock()
	// defer s.gm.mu.Unlock()

	if connectionId%2 == 0 {
		if prevGameState, exists := s.gm.games[connectionId-1]; exists {
			log.Printf("[server %d] second player connected", connectionId)
			prevGameState.connections[1] = connectionId
			s.gm.games[connectionId] = prevGameState
		} else {
			log.Printf("[server %d] ERROR No game found for previous player", connectionId)
			return
		}
	} else {
		log.Printf("[server %d] first player connected, creating new game", connectionId)
		gameState := &GameState{
			game:        game.NewGame(),
			connections: [2]int{connectionId, -1},
			readyChan:   make(chan string, 1),
		}

		s.gm.games[connectionId] = gameState
	}
}
