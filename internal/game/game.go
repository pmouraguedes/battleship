package game

type Game struct {
	players [2]*Player
}

func NewGame() *Game {
	return &Game{
		players: [2]*Player{},
	}
}

func (g *Game) GetPlayer(connectionId int) *Player {
	if connectionId%2 == 0 {
		return g.players[1]
	}
	return g.players[0]
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
