package go_carcassonne

import (
	"fmt"
	"strconv"
	"strings"

	bg "github.com/quibbble/go-boardgame"
	"github.com/quibbble/go-boardgame/pkg/bgn"
)

const key = "Carcassonne"

type Builder struct{}

func (b *Builder) Create(options *bg.BoardGameOptions) (bg.BoardGame, error) {
	return NewCarcassonne(options)
}

func (b *Builder) CreateWithBGN(options *bg.BoardGameOptions) (bg.BoardGameWithBGN, error) {
	return NewCarcassonne(options)
}

func (b *Builder) Load(game *bgn.Game) (bg.BoardGameWithBGN, error) {
	if game.Tags["Game"] != key {
		return nil, loadFailure(fmt.Errorf("game tag does not match game key"))
	}
	teamsStr, ok := game.Tags["Teams"]
	if !ok {
		return nil, loadFailure(fmt.Errorf("missing teams tag"))
	}
	teams := strings.Split(teamsStr, ", ")
	seedStr, ok := game.Tags["Seed"]
	if !ok {
		return nil, loadFailure(fmt.Errorf("missing seed tag"))
	}
	seed, err := strconv.Atoi(seedStr)
	if err != nil {
		return nil, loadFailure(err)
	}
	g, err := b.CreateWithBGN(&bg.BoardGameOptions{
		Teams: teams,
		MoreOptions: CarcassonneMoreOptions{
			Seed: int64(seed),
		},
	})
	if err != nil {
		return nil, err
	}
	for _, action := range game.Actions {
		if action.TeamIndex >= len(teams) {
			return nil, loadFailure(fmt.Errorf("team index %d out of range", action.TeamIndex))
		}
		team := teams[action.TeamIndex]
		actionType := notationToAction[string(action.ActionKey)]
		if actionType == "" {
			return nil, loadFailure(fmt.Errorf("invalid action key %s", string(action.ActionKey)))
		}
		var details interface{}
		switch actionType {
		case ActionPlaceTile:
			result, err := decodePlaceTileActionDetailsBGN(action.Details)
			if err != nil {
				return nil, err
			}
			details = result
		case ActionPlaceToken:
			result, err := decodePlaceTokenActionDetailsBGN(action.Details)
			if err != nil {
				return nil, err
			}
			details = result
		case bg.ActionSetWinners:
			result, err := bg.DecodeSetWinnersActionDetailsBGN(action.Details, teams)
			if err != nil {
				return nil, err
			}
			details = result
		}
		if err := g.Do(&bg.BoardGameAction{
			Team:        team,
			ActionType:  actionType,
			MoreDetails: details,
		}); err != nil {
			return nil, err
		}
	}
	return g, nil
}

func (b *Builder) Key() string {
	return key
}
