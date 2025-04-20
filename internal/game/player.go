package game

import (
	"fmt"
	"strconv"
)

type Player struct {
	id   int
	name string
	// TODO check usage of these public fields in handler.go
	Fleet *Fleet
}

func NewPlayer(id int, name string) *Player {
	fleet := newFleet()

	return &Player{
		id:    id,
		name:  name,
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

	ship, err := newShip(ShipType(shipType), xInt, yInt, s)
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
