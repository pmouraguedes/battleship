package game

import (
	"fmt"
	"log"
)

const (
	FLEET_UNIT_SIZE = 5*1 + 4*1 + 3*2 + 2*3 + 1*4
)

type Fleet struct {
	ships              map[ShipType][]*Ship
	positions          map[Vector2]*Ship
	remainingShipUnits int
	Ready              bool
	UnitSize           int
}

func newFleet() *Fleet {
	return &Fleet{
		ships:              make(map[ShipType][]*Ship),
		positions:          make(map[Vector2]*Ship),
		remainingShipUnits: FLEET_UNIT_SIZE,
	}
}

func (f *Fleet) addShip(ship *Ship) error {
	f.ships[ship.shipType] = append(f.ships[ship.shipType], ship)
	for _, position := range ship.positions {
		if _, exists := f.positions[position]; exists {
			return fmt.Errorf("position %v already occupied", position)
		}
		f.positions[position] = ship

		f.UnitSize++
	}
	return nil
}

func (f *Fleet) getShipAtPosition(position Vector2) (*Ship, bool) {
	ship, exists := f.positions[position]
	return ship, exists
}

func (f *Fleet) receiveAttack(position Vector2) (bool, *ShipType) {
	ship, exists := f.getShipAtPosition(position)
	if !exists {
		return false, nil
	}

	ship.receiveAttack()
	f.remainingShipUnits--

	log.Printf("[fleet] remaining ship units: %d", f.remainingShipUnits)

	if ship.isSunk() {
		return true, &ship.shipType
	} else {
		return true, nil
	}
}

func (f *Fleet) allShipsSunk() bool {
	return f.remainingShipUnits == 0
}
