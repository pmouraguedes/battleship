package game

import (
	"fmt"
	"strconv"
)

type PlayerStatus int

const (
	WAITING_FOR_HELLO PlayerStatus = iota
	SETUP_FLEET
	// PLAYING
)

type Player struct {
	id        int
	name      string
	Fleet     *Fleet
	TurnCount int
	State     PlayerStatus
}

func newPlayer(id int, name string) *Player {
	fleet := newFleet()

	return &Player{
		id:        id,
		name:      name,
		Fleet:     fleet,
		TurnCount: 1,
	}
}

func (p *Player) ReceiveAttack(x string, y string) (bool, *ShipType) {
	xInt, err := strconv.Atoi(x)
	if err != nil {
		panic(err)
	}
	yInt, err := strconv.Atoi(y)
	if err != nil {
		panic(err)
	}
	position := Vector2{xInt, yInt}
	return p.Fleet.receiveAttack(position)
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

func (p *Player) GetPlayerCode() string {
	return "P" + fmt.Sprintf("%d", p.getNumber())
}

func (p *Player) getNumber() int {
	number := 1
	if p.id%2 == 0 {
		number = 2
	}
	return number
}
