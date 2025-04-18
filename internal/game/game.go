package game

import (
	"fmt"
	"strconv"
)

type Fleet struct {
	Ships map[ShipType][]*Ship
}

type Player struct {
	id    int
	Name  string
	Fleet *Fleet
}

type Game struct {
	Players [2]*Player
}

func (g *Game) GetPlayer(connectionId int) *Player {
	if connectionId%2 == 0 {
		return g.Players[1]
	}
	return g.Players[0]
}

func NewFleet() *Fleet {
	return &Fleet{
		Ships: make(map[ShipType][]*Ship),
	}
}

// Player

func NewPlayer(id int, name string) *Player {
	fleet := NewFleet()

	return &Player{
		id:    id,
		Name:  name,
		Fleet: fleet,
	}
}

func (p *Player) AddShip(shipType string, x string, y string, s string) error {
	// Convert x and y to int
	xInt, err := strconv.Atoi(x)
	if err != nil {
		panic(err)
	}
	yInt, err := strconv.Atoi(y)
	if err != nil {
		panic(err)
	}

	ship, err := NewShip(ShipType(shipType), xInt, yInt, s)
	if err != nil {
		return err
	}
	p.Fleet.Ships[ShipType(shipType)] = append(p.Fleet.Ships[ShipType(shipType)], ship)

	return nil
}

func (p *Player) getNumber() int {
	number := 1
	if p.id%2 == 0 {
		number = 2
	}
	return number
}

func (p *Player) GetPlayerCode() string {
	return "P" + fmt.Sprintf("%d", p.getNumber())
}

// Game

func NewGame() *Game {
	// Create a new game with two players
	// Initialize the players

	return &Game{
		Players: [2]*Player{},
	}
}

// add player to game
func (g *Game) AddPlayer(player *Player) {
	if player.getNumber() == 1 {
		g.Players[0] = player
	} else {
		g.Players[1] = player
	}
}
