package go_carcassonne

// Token types
const (
	Farmer = "Farmer"
	Knight = "Knight"
	Thief  = "Thief"
	Monk   = "Monk"
)

var TokenTypes = []string{Farmer, Knight, Thief, Monk}

type token struct {
	X, Y int
	Team string
	Type string // Farmer, Knight, Thief, Monk
	Side string // normal side if Knight or Thief, farm side if Farmer, empty if Monk
}

func newToken(x, y int, team, typ, side string) *token {
	return &token{
		X:    x,
		Y:    y,
		Team: team,
		Type: typ,
		Side: side,
	}
}
