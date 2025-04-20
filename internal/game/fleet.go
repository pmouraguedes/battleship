package game

import "fmt"

const (
	FLEET_UNIT_SIZE = 5*1 + 4*1 + 3*2 + 2*3 + 1*4
)

type Fleet struct {
	ships     map[ShipType][]*Ship
	positions map[Vector2]*Ship
	// TODO check usage of these public fields in handler.go
	Ready    bool
	UnitSize int
}

func newFleet() *Fleet {
	return &Fleet{
		ships:     make(map[ShipType][]*Ship),
		positions: make(map[Vector2]*Ship),
	}
}

func (f *Fleet) addShip(ship *Ship) error {
	f.ships[ship.shipType] = append(f.ships[ship.shipType], ship)
	for _, position := range ship.positions {
		if _, exists := f.positions[position]; exists {
			return fmt.Errorf("position %v already occupied", position)
		}
		f.positions[position] = ship

		// increment fleet unit size
		f.UnitSize++
	}
	return nil
}

func (f *Fleet) getShipAtPosition(position Vector2) (*Ship, bool) {
	ship, exists := f.positions[position]
	return ship, exists
}
