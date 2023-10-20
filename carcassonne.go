package go_carcassonne

import (
	"fmt"
	"math/rand"
	"strings"

	"github.com/mitchellh/mapstructure"
	bg "github.com/quibbble/go-boardgame"
	"github.com/quibbble/go-boardgame/pkg/bgerr"
	"github.com/quibbble/go-boardgame/pkg/bgn"
)

const (
	minTeams = 2
	maxTeams = 5
)

type Carcassonne struct {
	state   *state
	actions []*bg.BoardGameAction
	options *CarcassonneMoreOptions
}

func NewCarcassonne(options *bg.BoardGameOptions) (*Carcassonne, error) {
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
	} else if duplicates(options.Teams) {
		return nil, &bgerr.Error{
			Err:    fmt.Errorf("duplicate teams found"),
			Status: bgerr.StatusInvalidOption,
		}
	}
	var details CarcassonneMoreOptions
	if err := mapstructure.Decode(options.MoreOptions, &details); err != nil {
		return nil, &bgerr.Error{
			Err:    err,
			Status: bgerr.StatusInvalidOption,
		}
	}
	return &Carcassonne{
		state:   newState(options.Teams, rand.New(rand.NewSource(details.Seed))),
		actions: make([]*bg.BoardGameAction, 0),
		options: &details,
	}, nil
}

func (c *Carcassonne) Do(action *bg.BoardGameAction) error {
	if len(c.state.winners) > 0 {
		return &bgerr.Error{
			Err:    fmt.Errorf("game already over"),
			Status: bgerr.StatusGameOver,
		}
	}
	switch action.ActionType {
	case ActionRotateTileRight:
		if err := c.state.RotateTileRight(action.Team); err != nil {
			return err
		}
	case ActionRotateTileLeft:
		if err := c.state.RotateTileLeft(action.Team); err != nil {
			return err
		}
	case ActionPlaceTile:
		var details PlaceTileActionDetails
		if err := mapstructure.Decode(action.MoreDetails, &details); err != nil {
			return &bgerr.Error{
				Err:    err,
				Status: bgerr.StatusInvalidActionDetails,
			}
		}
		tile := newTile(details.Tile.Top, details.Tile.Right, details.Tile.Bottom, details.Tile.Left, details.Tile.Center, details.Tile.ConnectedCitySides, details.Tile.Banner)
		if err := c.state.PlaceTile(action.Team, tile, details.X, details.Y); err != nil {
			return err
		}
		c.actions = append(c.actions, action)
	case ActionPlaceToken:
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
		c.actions = append(c.actions, action)
	case bg.ActionSetWinners:
		var details bg.SetWinnersActionDetails
		if err := mapstructure.Decode(action.MoreDetails, &details); err != nil {
			return &bgerr.Error{
				Err:    err,
				Status: bgerr.StatusInvalidActionDetails,
			}
		}
		if err := c.state.SetWinners(details.Winners); err != nil {
			return err
		}
		c.actions = append(c.actions, action)
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
	details := CarcassonneSnapshotData{
		LastPlacedTiles: c.state.lastPlacedTiles,
		Board:           c.state.board.board,
		BoardTokens:     c.state.boardTokens,
		Tokens:          c.state.tokens,
		Scores:          c.state.scores,
		TilesRemaining:  len(c.state.deck.tiles),
	}
	if len(team) == 1 {
		details.PlayTile = c.state.playTiles[team[0]]
	}
	var targets []*bg.BoardGameAction
	if len(c.state.winners) == 0 && (len(team) == 0 || (len(team) == 1 && team[0] == c.state.turn)) {
		targets = c.state.targets()
	}
	return &bg.BoardGameSnapshot{
		Turn:     c.state.turn,
		Teams:    c.state.teams,
		Winners:  c.state.winners,
		MoreData: details,
		Targets:  targets,
		Actions:  c.actions,
		Message:  c.state.message(),
	}, nil
}

func (c *Carcassonne) GetBGN() *bgn.Game {
	tags := map[string]string{
		"Game":  key,
		"Teams": strings.Join(c.state.teams, ", "),
		"Seed":  fmt.Sprintf("%d", c.options.Seed),
	}
	actions := make([]bgn.Action, 0)
	for _, action := range c.actions {
		bgnAction := bgn.Action{
			TeamIndex: indexOf(c.state.teams, action.Team),
			ActionKey: rune(actionToNotation[action.ActionType][0]),
		}
		switch action.ActionType {
		case ActionPlaceTile:
			var details PlaceTileActionDetails
			_ = mapstructure.Decode(action.MoreDetails, &details)
			bgnAction.Details = details.encodeBGN()
		case ActionPlaceToken:
			var details PlaceTokenActionDetails
			_ = mapstructure.Decode(action.MoreDetails, &details)
			bgnAction.Details = details.encodeBGN()
		case bg.ActionSetWinners:
			var details bg.SetWinnersActionDetails
			_ = mapstructure.Decode(action.MoreDetails, &details)
			bgnAction.Details, _ = details.EncodeBGN(c.state.teams)
		}
		actions = append(actions, bgnAction)
	}
	return &bgn.Game{
		Tags:    tags,
		Actions: actions,
	}
}
