package game

type GameStatus int

const (
	// WAITING_FOR_HELLO GameStatus = iota
	// SETUP_FLEET
	PLAYING GameStatus = iota
)

const (
	TURN_MAX_ATTACKS = 3
)

type Game struct {
	State     GameStatus
	players   [2]*Player
	TurnCount int
}

func NewGame() *Game {
	return &Game{
		players:   [2]*Player{},
		TurnCount: 1,
	}
}

func (g *Game) GetPlayer(connectionId int) *Player {
	if connectionId%2 == 0 {
		return g.players[1]
	}
	return g.players[0]
}

func (g *Game) GetOtherPlayer(connectionId int) *Player {
	return g.GetPlayer(connectionId ^ 1)
}

func (g *Game) IsReady() bool {
	if g.players[0] == nil || g.players[1] == nil {
		return false
	}
	return g.players[0].Fleet.Ready && g.players[1].Fleet.Ready
}

func (g *Game) AddPlayer(connectionId int, playerName string) *Player {
	player := newPlayer(connectionId, playerName)

	if player.getNumber() == 1 {
		g.players[0] = player
	} else {
		g.players[1] = player
	}

	return player
}

func (g *Game) IsPlayersTurn(player *Player) bool {
	if g.TurnCount%2 == 0 {
		return player.getNumber() == 2
	}
	return player.getNumber() == 1
}
