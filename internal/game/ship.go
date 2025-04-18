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
	Type      ShipType
	Length    int
	Positions []Vector2
}

type ShipSpec struct {
	Length  int
	Offsets []Vector2
}

var shipSpecsH = map[ShipType]ShipSpec{
	Carrier: {
		Length: 5,
		Offsets: []Vector2{
			{0, 0},
			{1, 0},
			{2, 0},
			{2, 1},
			{2, -1},
		},
	},
	Cruiser: {
		Length: 4,
		Offsets: []Vector2{
			{0, 0},
			{1, 0},
			{2, 0},
			{3, 0},
		},
	},
	Battleship: {
		Length: 3,
		Offsets: []Vector2{
			{0, 0},
			{1, 0},
			{2, 0},
		},
	},
	Destroyer: {
		Length: 2,
		Offsets: []Vector2{
			{0, 0},
			{1, 0},
		},
	},
	Submarine: {
		Length: 1,
		Offsets: []Vector2{
			{0, 0},
		},
	},
}

var shipSpecsV = map[ShipType]ShipSpec{
	Carrier: {
		Length: 5,
		Offsets: []Vector2{
			{0, 0},
			{0, 1},
			{0, 2},
			{1, 2},
			{-1, 2},
		},
	},
	Cruiser: {
		Length: 4,
		Offsets: []Vector2{
			{0, 0},
			{0, 1},
			{0, 2},
			{0, 3},
		},
	},
	Battleship: {
		Length: 3,
		Offsets: []Vector2{
			{0, 0},
			{0, 1},
			{0, 2},
		},
	},
	Destroyer: {
		Length: 2,
		Offsets: []Vector2{
			{0, 0},
			{0, 1},
		},
	},
	Submarine: {
		Length: 1,
		Offsets: []Vector2{
			{0, 0},
		},
	},
}

func NewShip(shipType ShipType, x int, y int, direction string) (*Ship, error) {
	var shipSpec ShipSpec
	switch direction {
	case "H":
		shipSpec = shipSpecsH[shipType]
	case "V":
		shipSpec = shipSpecsV[shipType]
	default:
		return nil, fmt.Errorf("invalid direction: %s", direction)
	}

	length := shipSpec.Length
	positions := make([]Vector2, length)

	for i := range length {
		newX := x + shipSpec.Offsets[i].X
		newY := y + shipSpec.Offsets[i].Y
		if !isValidCoordinate(newX) || !isValidCoordinate(newY) {
			return nil, fmt.Errorf("invalid coordinates: (%d, %d)", newX, newY)
		}
		positions[i] = Vector2{
			X: newX,
			Y: newY,
		}
	}

	return &Ship{
		Type:      shipType,
		Length:    length,
		Positions: positions,
	}, nil
}

func isValidCoordinate(n int) bool {
	if n < 0 || n > 9 {
		return false
	}
	return true
}
