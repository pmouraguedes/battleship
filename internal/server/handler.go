package server

import (
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/pmouraguedes/battleship/internal/game"
)

type GameState struct {
	game        *game.Game
	connections [2]int
	readyChan   chan int
}

type GameManager struct {
	games map[int]*GameState
	mu    sync.RWMutex
	// conns map[int]net.Conn
}

func newGameManager() *GameManager {
	return &GameManager{
		games: make(map[int]*GameState),
		// conns: make(map[int]net.Conn),
	}
}

func (gm *GameManager) getGame(connectionId int) *game.Game {
	return gm.games[connectionId].game
}

// func (gm *GameManager) AddConnection(conn net.Conn, connectionId int) {
// 	gm.mu.Lock()
// 	defer gm.mu.Unlock()
// gm.conns[connectionId] = conn
// }

// A new game is created as soon as the first client connects.
func (gm *GameManager) Handle(msg string, connectionId int) (string, error) {
	// message: HELLO <player_name>
	if strings.HasPrefix(msg, "HELLO") {
		return gm.handleHelloCommand(msg, connectionId)
	}

	// message: SHIP <ship_type> <x> <y> <direction>
	if strings.HasPrefix(msg, "SHIP") {
		return gm.handleShipCommand(msg, connectionId)
	}

	// message: READY
	if strings.HasPrefix(msg, "READY") {
		// set the player as ready and write to the channel
		err := gm.handleReadyCommand(msg, connectionId)
		if err != nil {
			return "ERROR\n", err
		}

		if gm.games[connectionId].game.IsReady() {
			return "START P1\n", nil
		}
		// else wait for the other player to be ready
		for range 2 {
			id := <-gm.games[connectionId].readyChan
			log.Println("Player", id, "is ready")
		}

		return "START P1\n", nil
	}

	return "ERROR Invalid command\n", fmt.Errorf("invalid command")
}

func (gm *GameManager) handleHelloCommand(msg string, connectionId int) (string, error) {
	// assert that the message is a HELLO command
	if !strings.HasPrefix(msg, "HELLO") {
		panic("Should be a HELLO command")
	}

	parts := strings.Fields(msg)
	if len(parts) != 2 {
		return "ERROR invalid HELLO command\n", fmt.Errorf("invalid HELLO command")
	}

	playerName := parts[1]
	if len(playerName) < 1 || len(playerName) > 20 {
		return "ERROR invalid player name\n", fmt.Errorf("invalid player name")
	}

	gm.mu.Lock()
	defer gm.mu.Unlock()

	// With HELLO, the game should not exist for this connectionId
	if _, exists := gm.games[connectionId]; exists {
		log.Println("Game already exists for connectionId:", connectionId)
		return "", fmt.Errorf("game already exists for connectionId: %d", connectionId)
	}

	// if connectionId is even, use the game of previous connectionId
	if connectionId%2 == 0 {
		if prevGameState, exists := gm.games[connectionId-1]; exists {
			log.Println("second player first message, associating with already created game")
			prevGameState.connections[1] = connectionId
			gm.games[connectionId] = prevGameState
		} else {
			return "ERROR No game found for previous player\n", fmt.Errorf("no game found for previous player")
		}
	} else {
		log.Println("Creating new game for connectionId:", connectionId)
		// Create a new game state
		gameState := &GameState{
			game:        game.NewGame(),
			connections: [2]int{connectionId, -1},
			readyChan:   make(chan int, 2),
		}

		gm.games[connectionId] = gameState
	}

	player := gm.getGame(connectionId).AddPlayer(connectionId, playerName)

	switch connectionId {
	case 1:
		return fmt.Sprintf("WELCOME %s %s\n", player.GetPlayerCode(), playerName), nil
	case 2:
		opponentName := "Opponent"
		return fmt.Sprintf("WELCOME %s %s %s\n", player.GetPlayerCode(), playerName, opponentName), nil
	default:
		panic("Only two players are allowed")
	}
}

func (gm *GameManager) handleShipCommand(msg string, connectionId int) (string, error) {
	// assert that the message is a SHIP command
	if !strings.HasPrefix(msg, "SHIP") {
		panic("Should be a SHIP command")
	}

	parts := strings.Fields(msg)
	if len(parts) != 5 {
		return "ERROR Invalid SHIP command\n", fmt.Errorf("invalid SHIP command")
	}
	shipType := parts[1]
	x := parts[2]
	y := parts[3]
	if !isValidShipType(shipType) {
		return "ERROR Invalid ship type\n", fmt.Errorf("invalid ship type")
	}
	if !isValidDirection(parts[4]) {
		return "ERROR Invalid direction\n", fmt.Errorf("invalid direction")
	}

	// Get player
	player := gm.games[connectionId].game.GetPlayer(connectionId)
	if player == nil {
		return "ERROR hello command not received yet\n", fmt.Errorf("hello command not received yet")
	}

	err := player.AddShip(shipType, x, y, parts[4])
	if err != nil {
		return "ERROR Invalid placement\n", fmt.Errorf("invalid placement")
	}

	return fmt.Sprintf("OK SHIP %s\n", shipType), nil
}

func (gm *GameManager) handleReadyCommand(_ string, connectionId int) error {
	log.Println("Handling ready command for connectionId:", connectionId)

	player := gm.games[connectionId].game.GetPlayer(connectionId)
	if player == nil {
		log.Println("Player not found")
		return fmt.Errorf("player not found")
	}
	if player.Fleet.Ready {
		log.Println("Player already ready")
		return fmt.Errorf("player already ready")
	}
	if player.Fleet.UnitSize < game.FLEET_UNIT_SIZE {
		log.Println("Player fleet not full")
		return fmt.Errorf("player fleet not full")
	}

	player.Fleet.Ready = true

	// write to the channel
	gm.games[connectionId].readyChan <- connectionId

	return nil
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

func isValidShipType(shipType string) bool {
	shipTypes := [5]string{"CARRIER", "BATTLESHIP", "CRUISER", "DESTROYER", "SUBMARINE"}
	for _, st := range shipTypes {
		if shipType == st {
			return true
		}
	}
	return false
}
