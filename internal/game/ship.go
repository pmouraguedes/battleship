package game

import "fmt"

type ShipType string

const (
	Carrier    ShipType = "CARRIER"
	Cruiser    ShipType = "CRUISER"
	Battleship ShipType = "BATTLESHIP"
	Destroyer  ShipType = "DESTROYER"
	Submarine  ShipType = "SUBMARINE"
)

type Vector2 struct {
	X, Y int
}

// Ship
type Ship struct {
	shipType  ShipType
	length    int
	positions []Vector2
}

type ShipSpec struct {
	length  int
	offsets []Vector2
}

var shipSpecsH = map[ShipType]ShipSpec{
	Carrier: {
		length: 5,
		offsets: []Vector2{
			{0, 0},
			{1, 0},
			{2, 0},
			{2, 1},
			{2, -1},
		},
	},
	Cruiser: {
		length: 4,
		offsets: []Vector2{
			{0, 0},
			{1, 0},
			{2, 0},
			{3, 0},
		},
	},
	Battleship: {
		length: 3,
		offsets: []Vector2{
			{0, 0},
			{1, 0},
			{2, 0},
		},
	},
	Destroyer: {
		length: 2,
		offsets: []Vector2{
			{0, 0},
			{1, 0},
		},
	},
	Submarine: {
		length: 1,
		offsets: []Vector2{
			{0, 0},
		},
	},
}

var shipSpecsV = map[ShipType]ShipSpec{
	Carrier: {
		length: 5,
		offsets: []Vector2{
			{0, 0},
			{0, 1},
			{0, 2},
			{1, 2},
			{-1, 2},
		},
	},
	Cruiser: {
		length: 4,
		offsets: []Vector2{
			{0, 0},
			{0, 1},
			{0, 2},
			{0, 3},
		},
	},
	Battleship: {
		length: 3,
		offsets: []Vector2{
			{0, 0},
			{0, 1},
			{0, 2},
		},
	},
	Destroyer: {
		length: 2,
		offsets: []Vector2{
			{0, 0},
			{0, 1},
		},
	},
	Submarine: {
		length: 1,
		offsets: []Vector2{
			{0, 0},
		},
	},
}

func newShip(shipType ShipType, x int, y int, direction string) (*Ship, error) {
	var shipSpec ShipSpec
	switch direction {
	case "H":
		shipSpec = shipSpecsH[shipType]
	case "V":
		shipSpec = shipSpecsV[shipType]
	default:
		return nil, fmt.Errorf("invalid direction: %s", direction)
	}

	length := shipSpec.length
	positions := make([]Vector2, length)

	for i := range length {
		newX := x + shipSpec.offsets[i].X
		newY := y + shipSpec.offsets[i].Y
		if !isValidCoordinate(newX) || !isValidCoordinate(newY) {
			return nil, fmt.Errorf("invalid coordinates: (%d, %d)", newX, newY)
		}
		positions[i] = Vector2{
			X: newX,
			Y: newY,
		}
	}

	return &Ship{
		shipType:  shipType,
		length:    length,
		positions: positions,
	}, nil
}

func isValidCoordinate(n int) bool {
	if n < 0 || n > 9 {
		return false
	}
	return true
}
