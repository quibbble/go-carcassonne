package go_carcassonne

import (
	"fmt"
	bg "github.com/quibbble/go-boardgame"
	"github.com/quibbble/go-boardgame/pkg/bgerr"
)

const (
	key      = "carcassonne"
	minTeams = 2
	maxTeams = 5
)

type Builder struct{}

func (b *Builder) Create(options bg.BoardGameOptions) (bg.BoardGame, error) {
	if len(options.Teams) < minTeams {
		return nil, &bgerr.Error{
			Err:    fmt.Errorf("at least %d teams required to create a game of %s", minTeams, key),
			Status: bgerr.StatusTooFewTeams,
		}
	} else if len(options.Teams) > maxTeams {
		return nil, &bgerr.Error{
			Err:    fmt.Errorf("at most %d teams allowed to create a game of %s", maxTeams, key),
			Status: bgerr.StatusTooManyTeams,
		}
	}
	return NewCarcassonne(options), nil
}

func (b *Builder) Key() string {
	return key
}
