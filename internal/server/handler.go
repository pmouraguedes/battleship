package server

import (
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/pmouraguedes/battleship/internal/game"
)

type GameState struct {
	Game        *game.Game
	Connections [2]int
}

type GameManager struct {
	games map[int]*GameState
	mu    sync.RWMutex
}

func NewGameManager() *GameManager {
	return &GameManager{
		games: make(map[int]*GameState),
	}
}

// A new game is created as soon as the first client connects.
func (gm *GameManager) Handle(message string, connectionId int) string {
	gm.mu.Lock()
	defer gm.mu.Unlock()

	// Check if the game already exists
	if _, exists := gm.games[connectionId]; exists {
		log.Println("Game already exists for connectionId:", connectionId)
		return gm.handleMessage(message, connectionId)
	}

	// if connectionId is even, use the game of previous connectionId
	if connectionId%2 == 0 {
		if prevGameState, exists := gm.games[connectionId-1]; exists {
			log.Println("Second player first message, associating with already created game")
			prevGameState.Connections[1] = connectionId
			gm.games[connectionId] = prevGameState
			return gm.handleMessage(message, connectionId)
		} else {
			panic("No game found for previous player")
		}
	}

	log.Println("Creating new game for connectionId:", connectionId)
	// Create a new game state
	gameState := &GameState{
		Game:        game.NewGame(),
		Connections: [2]int{connectionId, -1},
	}

	gm.games[connectionId] = gameState

	return gm.handleMessage(message, connectionId)
}

func (gm *GameManager) handleMessage(msg string, connectionId int) string {
	// message: HELLO <player_name>
	if strings.HasPrefix(msg, "HELLO") {
		return gm.handleHelloCommand(msg, connectionId)
	}

	// message: SHIP <ship_type> <x> <y> <direction>
	if strings.HasPrefix(msg, "SHIP") {
		return gm.handleShipCommand(msg, connectionId)
	}

	return "Invalid command\n"
}

func (gm *GameManager) handleHelloCommand(msg string, connectionId int) string {
	// assert that the message is a HELLO command
	if !strings.HasPrefix(msg, "HELLO") {
		panic("Should be a HELLO command")
	}

	parts := strings.Fields(msg)
	if len(parts) != 2 {
		return "Invalid HELLO command\n"
	}
	playerName := parts[1]
	if len(playerName) < 1 || len(playerName) > 20 {
		return "Invalid player name\n"
	}

	player := game.NewPlayer(connectionId, playerName)
	gm.games[connectionId].Game.AddPlayer(player)

	switch connectionId {
	case 1:
		return fmt.Sprintf("WELCOME %s %s\n", player.GetPlayerCode(), playerName)
	case 2:
		opponentName := "Opponent"
		return fmt.Sprintf("WELCOME %s %s %s\n", player.GetPlayerCode(), playerName, opponentName)
	default:
		panic("Only two players are allowed")
	}
}

func (gm *GameManager) handleShipCommand(msg string, connectionId int) string {
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
	if !isValidDirection(parts[4]) {
		return "Invalid direction\n"
	}

	// Get player
	player := gm.games[connectionId].Game.GetPlayer(connectionId)
	if player == nil {
		return "ERROR hello command not received yet"
	}

	err := player.AddShip(shipType, x, y, parts[4])
	if err != nil {
		return "ERROR Invalid placement"
	}

	return fmt.Sprintf("OK SHIP %s\n", shipType)
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
