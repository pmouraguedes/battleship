package game

import (
	"fmt"
	"strconv"
)

type Fleet struct {
	Ships     map[ShipType][]*Ship
	Positions map[Vector2]*Ship
	Ready     bool
	UnitSize  int
}

type Player struct {
	id    int
	Name  string
	Fleet *Fleet
}

type Game struct {
	Players [2]*Player
}

const (
	FLEET_UNIT_SIZE = 5*1 + 4*1 + 3*2 + 2*3 + 1*4
)

// Game

func (g *Game) GetPlayer(connectionId int) *Player {
	if connectionId%2 == 0 {
		return g.Players[1]
	}
	return g.Players[0]
}

func (g *Game) IsReady() bool {
	if g.Players[0] == nil || g.Players[1] == nil {
		return false
	}
	return g.Players[0].Fleet.Ready && g.Players[1].Fleet.Ready
}

// Fleet

func NewFleet() *Fleet {
	return &Fleet{
		Ships:     make(map[ShipType][]*Ship),
		Positions: make(map[Vector2]*Ship),
	}
}

func (f *Fleet) addShip(ship *Ship) error {
	f.Ships[ship.Type] = append(f.Ships[ship.Type], ship)
	for _, position := range ship.Positions {
		if _, exists := f.Positions[position]; exists {
			return fmt.Errorf("position %v already occupied", position)
		}
		f.Positions[position] = ship

		// increment fleet unit size
		f.UnitSize++
	}
	return nil
}

func (f *Fleet) GetShipAtPosition(position Vector2) (*Ship, bool) {
	ship, exists := f.Positions[position]
	return ship, exists
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
	err = p.Fleet.addShip(ship)
	return err
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

func (g *Game) AddPlayer(player *Player) {
	if player.getNumber() == 1 {
		g.Players[0] = player
	} else {
		g.Players[1] = player
	}
}
