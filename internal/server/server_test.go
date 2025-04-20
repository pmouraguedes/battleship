package server

import (
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
	"testing"
	"time"
)

var (
	turnPlayerCode = "P1"
	dualChan       = make(chan int, 2)
	singleChan     = make(chan int, 1)
	// finishChan     = make(chan int, 1)
	errChan = make(chan error, 2)
	wg      sync.WaitGroup
)

func TestServer(t *testing.T) {
	// Create a new server instance
	s := NewServer(":8000")

	// Start the server in a goroutine
	go s.Start()

	// Allow some time for the server to start
	time.Sleep(1 * time.Millisecond)

	// Create a connection to the server
	conn1 := startConnection(t, ":8000")
	defer conn1.Close()

	// Create a second connection to the server
	conn2 := startConnection(t, ":8000")
	defer conn2.Close()

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
	log.Printf("Finished waiting for clients")
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

	<-dualChan
	log.Printf("Client %s received dualChan signal", clientCode)

	// sleep for a while to allow the server to process the messages
	// time.Sleep(100 * time.Millisecond)

	// ATTACK messages
	log.Printf("Client %s sending attack messages", clientCode)
	sendAttackMessages(conn, clientCode)
}

func sendAttackMessages(conn net.Conn, clientCode string) {
	i := 0
	j := 0
	// n := 0

	if turnPlayerCode != clientCode {
		<-singleChan
	}

	for {
		for range 1 {
			if turnPlayerCode != clientCode {
				log.Printf("Client %s waiting for turn", clientCode)
				<-singleChan
			}

			attackMessage := fmt.Sprintf("ATTACK %d %d", i, j)

			sendMessage(conn, attackMessage)
			resp := readResponse(conn)
			// log.Printf("----Client %s received: %s", clientCode, resp)
			if !strings.HasPrefix(resp, "HIT") &&
				!strings.HasPrefix(resp, "MISS") &&
				!strings.HasPrefix(resp, "SUNK") {
				// log.Printf("-----Client %s received: %s", clientCode, resp)
				switchTurnPlayerCode()
				singleChan <- 1
				errChan <- fmt.Errorf("Expected HIT or MISS message, got: %s", resp)
				return
			}
			j++
			if j > 9 {
				j = 0
				i++
			}
			if i > 9 {
				// log.Printf("i > 9")
				break
			}
		}

		switchTurnPlayerCode()

		log.Printf("Client %s sending singleChan signal", clientCode)
		singleChan <- 1

		if i > 9 {
			// log.Printf("i > 9")
			break
		}
	}
	log.Printf("Client %s finished sending attack messages", clientCode)

	errChan <- nil
}

func sendHelloMessage(conn net.Conn, clientCode string, clientName string) {
	helloMessage := "HELLO " + clientName
	sendMessage(conn, helloMessage)

	response := readResponse(conn)
	if response != "WELCOME "+clientCode+" "+clientName+"\n" {
		errChan <- fmt.Errorf("Expected welcome message, got: %s", response)
		return
	}

	log.Printf("Client %s received: %s", clientName, response)
}

func sendFleetMessages(conn net.Conn, clientCode string) {
	// carrier
	carrierMessage := "SHIP CARRIER 1 1 H"
	sendMessage(conn, carrierMessage)
	response := readResponse(conn)
	if response != "OK SHIP CARRIER\n" {
		// t.Fatalf("Expected OK message, got: %s", response)
		errChan <- fmt.Errorf("Expected OK message, got: %s", response)
		return
	}
	// cruiser
	cruiserMessage := "SHIP CRUISER 5 0 V"
	sendMessage(conn, cruiserMessage)
	response = readResponse(conn)
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
		sendMessage(conn, msg)
		response = readResponse(conn)
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
		sendMessage(conn, msg)
		response = readResponse(conn)
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
		sendMessage(conn, msg)
		response = readResponse(conn)
		if response != "OK SHIP SUBMARINE\n" {
			errChan <- fmt.Errorf("Expected OK message, got: %s", response)
			return
		}
	}

	log.Printf("Client %s finished sending fleet messages", clientCode)
	sendMessage(conn, "READY")

	response = readResponse(conn)
	if response != "START P1\n" {
		errChan <- fmt.Errorf("Expected START message, got: %s", response)
		return
	}

	dualChan <- 1
}

func sendMessage(conn net.Conn, message string) {
	_, err := conn.Write([]byte(message))
	if err != nil {
		log.Println("Error sending message:", err)
		return
	}
}

func readResponse(conn net.Conn) string {
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		log.Println("Error reading response:", err)
		return ""
	}
	return string(buf[:n])
}

func switchTurnPlayerCode() {
	if turnPlayerCode == "P1" {
		turnPlayerCode = "P2"
	} else {
		turnPlayerCode = "P1"
	}
}
