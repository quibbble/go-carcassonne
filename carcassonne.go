package go_carcassonne

import (
	"fmt"
	"github.com/mitchellh/mapstructure"
	bg "github.com/quibbble/go-boardgame"
	"github.com/quibbble/go-boardgame/pkg/bgerr"
	"github.com/quibbble/go-boardgame/pkg/bgn"
	"math/rand"
	"strings"
)

const (
	minTeams = 2
	maxTeams = 5
)

type Carcassonne struct {
	state   *state
	actions []*bg.BoardGameAction
	seed    int64
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
	}
	return &Carcassonne{
		state:   newState(options.Teams, rand.New(rand.NewSource(options.Seed))),
		actions: make([]*bg.BoardGameAction, 0),
		seed:    options.Seed,
	}, nil
}

func (c *Carcassonne) Do(action *bg.BoardGameAction) error {
	switch action.ActionType {
	case ActionRotateTileRight:
		if err := c.state.RotateTileRight(action.Team); err != nil {
			return err
		}
		if len(c.actions) > 0 && c.actions[len(c.actions)-1].ActionType == ActionRotateTileLeft {
			// last action was RotateTileLeft so RotateTileRight undoes RotateTileLeft
			c.actions = c.actions[:len(c.actions)-1]
		} else if len(c.actions) > 2 &&
			c.actions[len(c.actions)-1].ActionType == ActionRotateTileRight &&
			c.actions[len(c.actions)-2].ActionType == ActionRotateTileRight &&
			c.actions[len(c.actions)-3].ActionType == ActionRotateTileRight {
			// last action three actions were RotateTileRight so fourth RotateTileRight undoes past three
			c.actions = c.actions[:len(c.actions)-3]
		} else {
			c.actions = append(c.actions, action)
		}
	case ActionRotateTileLeft:
		if err := c.state.RotateTileLeft(action.Team); err != nil {
			return err
		}
		if len(c.actions) > 0 && c.actions[len(c.actions)-1].ActionType == ActionRotateTileRight {
			// last action was RotateTileRight so RotateTileLeft undoes RotateTileRight
			c.actions = c.actions[:len(c.actions)-1]
		} else if len(c.actions) > 2 &&
			c.actions[len(c.actions)-1].ActionType == ActionRotateTileLeft &&
			c.actions[len(c.actions)-2].ActionType == ActionRotateTileLeft &&
			c.actions[len(c.actions)-3].ActionType == ActionRotateTileLeft {
			// last action three actions were RotateTileLeft so fourth RotateTileLeft undoes past three
			c.actions = c.actions[:len(c.actions)-3]
		} else {
			c.actions = append(c.actions, action)
		}
	case ActionPlaceTile:
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
	}, nil
}

func (c *Carcassonne) GetBGN() *bgn.Game {
	tags := map[string]string{
		"Game":  key,
		"Teams": strings.Join(c.state.teams, ", "),
		"Seed":  fmt.Sprintf("%d", c.seed),
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
			bgnAction.Details = details.encode()
		case ActionPlaceToken:
			var details PlaceTokenActionDetails
			_ = mapstructure.Decode(action.MoreDetails, &details)
			bgnAction.Details = details.encode()
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
