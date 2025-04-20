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
func (gm *GameManager) handle(msg string, connectionId int) (string, error) {
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
		return gm.handleReadyCommand(msg, connectionId)
	}

	// message: ATTACK <x> <y>
	if strings.HasPrefix(msg, "ATTACK") {
		return gm.handleAttackCommand(msg, connectionId)
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
		return "ERROR Game already exists\n", fmt.Errorf("game already exists for connectionId %d", connectionId)
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
			readyChan:   make(chan int),
		}

		gm.games[connectionId] = gameState
	}

	player := gm.getGame(connectionId).AddPlayer(connectionId, playerName)

	switch connectionId {
	case 1:
		return fmt.Sprintf("WELCOME %s %s\n", player.GetPlayerCode(), playerName), nil
	case 2:
		return fmt.Sprintf("WELCOME %s %s\n", player.GetPlayerCode(), playerName), nil
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
	player := gm.getGame(connectionId).GetPlayer(connectionId)
	if player == nil {
		return "ERROR hello command not received yet\n", fmt.Errorf("hello command not received yet")
	}

	err := player.AddShip(shipType, x, y, parts[4])
	if err != nil {
		return "ERROR Invalid placement\n", fmt.Errorf("invalid placement")
	}

	return fmt.Sprintf("OK SHIP %s\n", shipType), nil
}

func (gm *GameManager) handleReadyCommand(_ string, connectionId int) (string, error) {
	log.Println("Handling ready command for connectionId:", connectionId)

	player := gm.games[connectionId].game.GetPlayer(connectionId)
	if player == nil {
		log.Println("Player not found")
		return "ERROR player not found\n", fmt.Errorf("player not found")
	}
	if player.Fleet.Ready {
		log.Println("Player already ready")
		return "ERROR player already ready\n", fmt.Errorf("player already ready")
	}
	if player.Fleet.UnitSize < game.FLEET_UNIT_SIZE {
		log.Println("Player fleet not full")
		return "ERROR player fleet not full\n", fmt.Errorf("player fleet not full")
	}

	player.Fleet.Ready = true

	if gm.games[connectionId].game.IsReady() {
		// write to the channel to notify the other player
		gm.games[connectionId].readyChan <- connectionId
		return "START P1\n", nil
	}

	// else wait for the other player to be ready
	id := <-gm.games[connectionId].readyChan
	log.Println("Player", id, "is ready")

	return "START P1\n", nil
}

func (gm *GameManager) handleAttackCommand(msg string, connectionId int) (string, error) {
	// assert that the message is a ATTACK command
	if !strings.HasPrefix(msg, "ATTACK") {
		panic("Should be a ATTACK command")
	}

	parts := strings.Fields(msg)
	if len(parts) != 3 {
		return "ERROR Invalid ATTACK command\n", fmt.Errorf("invalid ATTACK command")
	}

	thisGame := gm.getGame(connectionId)

	player := thisGame.GetPlayer(connectionId)
	if player == nil {
		return "ERROR hello command not received yet\n", fmt.Errorf("hello command not received yet")
	}

	opponent := thisGame.GetOtherPlayer(connectionId)
	if opponent == nil {
		return "ERROR opponent not found\n", fmt.Errorf("opponent not found")
	}

	if !thisGame.IsPlayersTurn(player) {
		return "ERROR not your turn\n", fmt.Errorf("not your turn")
	}

	log.Printf("Game turn count: %d", thisGame.TurnCount)
	log.Printf("Player turn count: %d", player.TurnCount)
	if player.TurnCount >= game.TURN_MAX_ATTACKS {
		player.TurnCount = 1
		gm.getGame(connectionId).TurnCount++
	} else {
		player.TurnCount++
	}

	x := parts[1]
	y := parts[2]
	hit, sunkShipType := opponent.ReceiveAttack(x, y)

	// TODO check if the game is over

	if hit {
		if sunkShipType != nil {
			// sunk
			return fmt.Sprintf("SUNK %s %s %s\n", x, y, *sunkShipType), nil
		} else {
			// hit but not sunk
			return fmt.Sprintf("HIT %s %s\n", x, y), nil
		}
	}

	return fmt.Sprintf("MISS %s %s\n", x, y), nil
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
