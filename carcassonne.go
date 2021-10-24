package go_carcassonne

import (
	"fmt"
	"github.com/mitchellh/mapstructure"
	bg "github.com/quibbble/go-boardgame"
	"github.com/quibbble/go-boardgame/pkg/bgerr"
)

const (
	minTeams = 2
	maxTeams = 5
)

type Carcassonne struct {
	state   *state
	actions []*bg.BoardGameAction
}

func NewCarcassonne(options bg.BoardGameOptions) (*Carcassonne, error) {
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
	return &Carcassonne{
		state:   newState(options.Teams),
		actions: make([]*bg.BoardGameAction, 0),
	}, nil
}

func (c *Carcassonne) Do(action bg.BoardGameAction) error {
	switch action.ActionType {
	case RotateTileRight:
		if err := c.state.RotateTileRight(action.Team); err != nil {
			return err
		}
	case RotateTileLeft:
		if err := c.state.RotateTileLeft(action.Team); err != nil {
			return err
		}
	case PlaceTile:
		var details PlaceTileActionDetails
		if err := mapstructure.Decode(action.MoreDetails, &details); err != nil {
			return &bgerr.Error{
				Err:    err,
				Status: bgerr.StatusInvalidActionDetails,
			}
		}
		details.Tile = c.state.playTile
		if err := c.state.PlaceTile(action.Team, details.X, details.Y); err != nil {
			return err
		}
		action.MoreDetails = details
		c.actions = append(c.actions, &action)
	case PlaceToken:
		var details PlaceTokenActionDetails
		if err := mapstructure.Decode(action.MoreDetails, &details); err != nil {
			return &bgerr.Error{
				Err:    err,
				Status: bgerr.StatusInvalidActionDetails,
			}
		}
		if err := c.state.PlaceToken(action.Team, details.Pass, details.X, details.Y, details.Type, details.Side); err != nil {
			return err
		}
		c.actions = append(c.actions, &action)
	case Reset:
		c.state = newState(c.state.teams)
		c.actions = make([]*bg.BoardGameAction, 0)
	default:
		return &bgerr.Error{
			Err:    fmt.Errorf("cannot process action type %s", action.ActionType),
			Status: bgerr.StatusUnknownActionType,
		}
	}
	return nil
}

func (c *Carcassonne) GetSnapshot(team ...string) (*bg.BoardGameSnapshot, error) {
	if len(team) > 1 {
		return nil, &bgerr.Error{
			Err:    fmt.Errorf("get snapshot requires zero or one team"),
			Status: bgerr.StatusTooManyTeams,
		}
	}
	details := CarcassonneSnapshotDetails{
		PlayTile:       c.state.playTile,
		LastPlacedTile: c.state.lastPlacedTile,
		Board:          c.state.board.board,
		BoardTokens:    c.state.boardTokens,
		Tokens:         c.state.tokens,
		Scores:         c.state.scores,
		TilesRemaining: len(c.state.deck.tiles),
	}
	if len(team) == 1 && c.state.turn != team[0] {
		details.PlayTile = nil
		details.LastPlacedTile = nil
	}
	return &bg.BoardGameSnapshot{
		Turn:     c.state.turn,
		Teams:    c.state.teams,
		Winners:  c.state.winners,
		MoreData: details,
		Actions:  c.actions,
	}, nil
}
