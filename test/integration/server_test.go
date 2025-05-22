package server_test

import (
	"fmt"
	"io"
	"log"
	"net"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/pmouraguedes/battleship/internal/game"
	"github.com/pmouraguedes/battleship/internal/server"
)

var POSITIONS = [25]game.Vector2{
	// carrier
	{X: 1, Y: 1},
	{X: 2, Y: 1},
	{X: 3, Y: 1},
	{X: 3, Y: 2},
	{X: 3, Y: 0},
	// cruiser
	{X: 5, Y: 0},
	{X: 5, Y: 1},
	{X: 5, Y: 2},
	{X: 5, Y: 3},
	// battleship 1
	{X: 6, Y: 7},
	{X: 7, Y: 7},
	{X: 8, Y: 7},
	// battleship 2
	{X: 0, Y: 7},
	{X: 1, Y: 7},
	{X: 2, Y: 7},
	// destroyer 1
	{X: 4, Y: 5},
	{X: 5, Y: 5},
	// destroyer 2
	{X: 0, Y: 3},
	{X: 0, Y: 4},
	// destroyer 3
	{X: 7, Y: 0},
	{X: 8, Y: 0},
	// submarine 1
	{X: 9, Y: 2},
	// submarine 2
	{X: 9, Y: 4},
	// submarine 3
	{X: 9, Y: 9},
	// submarine 4
	{X: 0, Y: 9},
}

var (
	turnPlayerCode = "P1"
	errChan        = make(chan error, 2)
	wg             sync.WaitGroup
)

func TestServer(t *testing.T) {
	// Create a new server instance
	s := server.NewServer(":8000")

	// Start the server in a goroutine
	go s.Start()

	// Allow some time for the server to start
	time.Sleep(1 * time.Millisecond)

	// Create a connection to the server
	conn1 := startConnection(t, ":8000")
	defer conn1.Close()
	log.Printf("[test] conn1 created")

	// Create a second connection to the server
	conn2 := startConnection(t, ":8000")
	defer conn2.Close()
	log.Printf("[test] conn2 created")

	// Create a channel to synchronize the two clients
	wg.Add(1)
	go func() {
		defer wg.Done()
		doClientStuff(conn1, "P1", "Player1")
	}()
	time.Sleep(100 * time.Millisecond)

	wg.Add(1)
	go func() {
		defer wg.Done()
		doClientStuff(conn2, "P2", "Player2")
	}()

	wg.Wait()
	log.Printf("[test] Finished waiting for clients")
	close(errChan)

	for err := range errChan {
		if err != nil {
			t.Fatalf("Error in client: %v", err)
		}
	}
}

func startConnection(t *testing.T, address string) net.Conn {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		t.Fatalf("Failed to connect to server: %v", err)
	}
	return conn
}

func doClientStuff(conn net.Conn, clientCode string, clientName string) {
	// HELLO message
	sendHelloMessage(conn, clientCode, clientName)

	// SHIP messages
	sendFleetMessages(conn, clientCode)

	attacks := 0

	for {
		msg, err := readResponse(conn)
		if err == io.EOF {
			log.Printf("[client %s] connection closed by server", clientCode)
			break
		} else if err != nil {
			errChan <- fmt.Errorf("Error reading response: %v", err)
			return
		}

		log.Printf("[client %s] received: %s", clientCode, msg)
		if msg == fmt.Sprintf("TURN %s\n", clientCode) {
			log.Printf("[client %s] received TURN message", clientCode)

			// ATTACK messages
			sendAttackMessages(conn, clientCode, attacks)
			attacks += 3
		} else if gameOver := strings.HasPrefix(msg, "WIN"); gameOver {
			log.Printf("[client %s] received WIN message", clientCode)
		}
	}
}

func sendAttackMessages(conn net.Conn, clientCode string, attacks int) {
	for i := 0; i != game.TURN_MAX_ATTACKS; i++ {
		log.Printf("[client %s] sending attack message #%d", clientCode, attacks+i)

		position := POSITIONS[attacks+i]

		attackMessage := fmt.Sprintf("ATTACK %d %d\n", position.X, position.Y)
		sendClientMessage(conn, attackMessage)
		resp, err := readResponse(conn)
		if err != nil {
			errChan <- fmt.Errorf("Error reading response: %v", err)
			return
		}
		log.Printf("[client %s] received: %s", clientCode, resp)

		if gameOver := strings.HasPrefix(resp, "WIN"); gameOver {
			log.Printf("[client %s] received WIN message", clientCode)
			break
		}
		if !strings.HasPrefix(resp, "HIT") &&
			!strings.HasPrefix(resp, "MISS") &&
			!strings.HasPrefix(resp, "SUNK") {
			log.Printf("[client %s] received unexpected: %s", clientCode, resp)
			errChan <- fmt.Errorf("Expected HIT or MISS message, got: %s", resp)
			return
		}
	}
}

func sendHelloMessage(conn net.Conn, clientCode string, clientName string) {
	helloMessage := "HELLO " + clientName + "\n"
	sendClientMessage(conn, helloMessage)

	response, err := readResponse(conn)
	if err != nil {
		errChan <- fmt.Errorf("Error reading response: %v", err)
		return
	}
	if response != "WELCOME "+clientCode+" "+clientName+"\n" {
		errChan <- fmt.Errorf("Expected welcome message, got: %s", response)
		return
	}

	log.Printf("[client] %s received: %s", clientCode, response)
}

func sendFleetMessages(conn net.Conn, clientCode string) {
	// carrier
	carrierMessage := "SHIP CARRIER 1 1 H\n"
	sendClientMessage(conn, carrierMessage)
	response, err := readResponse(conn)
	if err != nil {
		errChan <- fmt.Errorf("Error reading response: %v", err)
		return
	}
	if response != "OK SHIP CARRIER\n" {
		// t.Fatalf("Expected OK message, got: %s", response)
		errChan <- fmt.Errorf("Expected OK message, got: %s", response)
		return
	}
	// cruiser
	cruiserMessage := "SHIP CRUISER 5 0 V"
	sendClientMessage(conn, cruiserMessage)
	response, err = readResponse(conn)
	if err != nil {
		errChan <- fmt.Errorf("Error reading response: %v", err)
		return
	}
	if response != "OK SHIP CRUISER\n" {
		errChan <- fmt.Errorf("Expected OK message, got: %s", response)
		return
	}
	// 2 battleships
	battleshipMessages := []string{
		"SHIP BATTLESHIP 6 7 H",
		"SHIP BATTLESHIP 0 7 H",
	}
	for _, msg := range battleshipMessages {
		sendClientMessage(conn, msg)
		response, err := readResponse(conn)
		if err != nil {
			errChan <- fmt.Errorf("Error reading response: %v", err)
			return
		}
		if response != "OK SHIP BATTLESHIP\n" {
			errChan <- fmt.Errorf("Expected OK message, got: %s", response)
			return
		}
	}
	// 3 destroyers
	destroyerMessages := []string{
		"SHIP DESTROYER 4 5 H",
		"SHIP DESTROYER 0 3 V",
		"SHIP DESTROYER 7 0 H",
	}
	for _, msg := range destroyerMessages {
		sendClientMessage(conn, msg)
		response, err := readResponse(conn)
		if err != nil {
			errChan <- fmt.Errorf("Error reading response: %v", err)
			return
		}
		if response != "OK SHIP DESTROYER\n" {
			errChan <- fmt.Errorf("Expected OK message, got: %s", response)
			return
		}
	}
	// 4 submarines
	submarineMessages := []string{
		"SHIP SUBMARINE 9 2 V",
		"SHIP SUBMARINE 9 4 V",
		"SHIP SUBMARINE 9 9 H",
		"SHIP SUBMARINE 0 9 H",
	}
	for _, msg := range submarineMessages {
		sendClientMessage(conn, msg)
		response, err := readResponse(conn)
		if err != nil {
			errChan <- fmt.Errorf("Error reading response: %v", err)
			return
		}

		if response != "OK SHIP SUBMARINE\n" {
			errChan <- fmt.Errorf("Expected OK message, got: %s", response)
			return
		}
	}

	log.Printf("[client %s] finished sending fleet messages", clientCode)
	sendClientMessage(conn, "READY\n")

	response, err = readResponse(conn)
	if err != nil {
		errChan <- fmt.Errorf("Error reading response: %v", err)
		return
	}

	if response != "START P1\n" {
		errChan <- fmt.Errorf("Expected START message, got: %s", response)
		return
	}
	log.Printf("[client %s] received START P1 message", clientCode)
}

func readResponse(conn net.Conn) (string, error) {
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		if err == io.EOF {
			log.Println("[test] connection closed by server")
		} else {
			log.Println("Error reading response:", err)
		}
		return "", err
	}

	return string(buf[:n]), nil
}

func sendClientMessage(conn net.Conn, message string) error {
	log.Printf("[client] sendClientMessage - sending: %s", message)
	_, err := conn.Write([]byte(message))
	if err != nil {
		log.Println("Error sending client message:", err)
		return err
	}
	return nil
}
