package main

import (
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
)

func handleShipCommand(msg string, connectionId int) string {
	// assert that the message is a SHIP command
	if !strings.HasPrefix(msg, "SHIP") {
		panic("Should be a SHIP command")
	}

	parts := strings.Fields(msg)
	if len(parts) != 5 {
		return "Invalid SHIP command\n"
	}
	shipType := parts[1]
	x := parts[2]
	y := parts[3]
	if !isValidShipType(shipType) {
		return "Invalid ship type\n"
	}
	if !isValidCoordinate(x) || !isValidCoordinate(y) {
		return "Invalid coordinates\n"
	}
	if !isValidDirection(parts[4]) {
		return "Invalid direction\n"
	}

	return fmt.Sprintf("OK SHIP %s\n", shipType)
}

func handleMessage(msg string, connectionId int) string {
	// if HELLO command
	// message: HELLO <player_name>
	if strings.HasPrefix(msg, "HELLO") {
		return handleHelloCommand(msg, connectionId)
	}

	// if SHIP command
	// message: SHIP <ship_type> <x> <y> <direction>
	if strings.HasPrefix(msg, "SHIP") {
		return handleShipCommand(msg, connectionId)
	}

	return "Invalid command\n"
}

func handleHelloCommand(msg string, connectionId int) string {
	// assert that the message is a HELLO command
	if !strings.HasPrefix(msg, "HELLO") {
		panic("Should be a HELLO command")
	}

	parts := strings.Fields(msg)
	if len(parts) != 2 {
		return "Invalid HELLO command\n"
	}
	playerName := parts[1]
	if len(playerName) < 3 || len(playerName) > 20 {
		return "Invalid player name\n"
	}
	playerCode := "P" + fmt.Sprintf("%d", connectionId)

	switch connectionId {
	case 1:
		return fmt.Sprintf("WELCOME %s %s\n", playerCode, playerName)
	case 2:
		opponentName := "Opponent"
		return fmt.Sprintf("WELCOME %s %s %s\n", playerCode, playerName, opponentName)
	default:
		panic("Only two players are allowed")
	}
}

func isValidDirection(s string) bool {
	if len(s) != 1 {
		return false
	}
	if s[0] != 'H' && s[0] != 'V' {
		return false
	}
	return true
}

func isValidCoordinate(x string) bool {
	if len(x) != 1 {
		return false
	}
	if x[0] < '0' || x[0] > '9' {
		return false
	}
	return true
}

func isValidShipType(shipType string) bool {
	shipTypes := [5]string{"CARRIER", "BATTLESHIP", "CRUISER", "DESTROYER", "SUBMARINE"}
	for _, st := range shipTypes {
		if shipType == st {
			return true
		}
	}
	return false
}

func handleConnection(conn net.Conn, connectionId int) {
	defer conn.Close()
	buf := make([]byte, 1024)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			log.Println(err)
			return
		}
		fmt.Printf("[%d] Received: %s", connectionId, string(buf[:n]))

		data := handleMessage(string(buf[:n]), connectionId)

		_, err = conn.Write([]byte(data))
		if err != nil {
			log.Println(err)
			return
		}
	}
}

func main() {
	addr, err := net.ResolveTCPAddr("tcp", ":8000")
	if err != nil {
		log.Fatal(err)
	}
	ln, err := net.ListenTCP("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close()
	fmt.Println("Listening on port 8000")

	connectionId := 0
	var wg sync.WaitGroup
	wg.Add(2)

	// Only accept the first two connections
	for connectionId < 2 {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal(err)
		}
		connectionId++
		go handleConnection(conn, connectionId)
	}

	wg.Wait()
	fmt.Println("Server shutting down")
}
