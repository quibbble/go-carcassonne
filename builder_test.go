package go_carcassonne

import (
	"encoding/json"
	bg "github.com/quibbble/go-boardgame"
	"github.com/quibbble/go-boardgame/pkg/bgn"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Builder_BGN(t *testing.T) {
	builder := Builder{}
	teams := []string{TeamA, TeamB}
	carcassonne, err := builder.CreateWithBGN(&bg.BoardGameOptions{
		Teams: teams,
		MoreOptions: CarcassonneMoreOptions{
			Seed: 123,
		}})
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	err = carcassonne.Do(&bg.BoardGameAction{
		Team:       TeamA,
		ActionType: ActionPlaceTile,
		MoreDetails: PlaceTileActionDetails{
			X: 1,
			Y: 0,
		},
	})
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	err = carcassonne.Do(&bg.BoardGameAction{
		Team:       TeamA,
		ActionType: ActionPlaceToken,
		MoreDetails: PlaceTokenActionDetails{
			Pass: false,
			X:    1,
			Y:    0,
			Type: Thief,
			Side: SideLeft,
		},
	})
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	game := carcassonne.GetBGN()
	carcassonneLoaded, err := builder.Load(game)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	expected, _ := json.Marshal(carcassonne)
	actual, _ := json.Marshal(carcassonneLoaded)
	assert.Equal(t, expected, actual)
}

func Test_Builder_Notations(t *testing.T) {
	tags := map[string]string{
		"Game":  key,
		"Teams": "TeamA, TeamB",
		"Seed":  "123",
	}
	tests := []struct {
		name        string
		bgn         *bgn.Game
		shouldError bool
	}{
		{
			name:        "empty string should error",
			bgn:         &bgn.Game{},
			shouldError: true,
		},
		{
			name: "missing seed should error",
			bgn: &bgn.Game{
				Tags: map[string]string{
					"Game":  key,
					"Teams": "TeamA, TeamB",
				},
			},
			shouldError: true,
		},
		{
			name: "should create a new game",
			bgn: &bgn.Game{
				Tags: tags,
			},
			shouldError: false,
		},
		{
			name: "should create a new game and do actions",
			bgn: &bgn.Game{
				Tags: tags,
				Actions: []bgn.Action{
					{
						TeamIndex: 0,
						ActionKey: 'i',
						Details:   []string{"1", "0"},
					},
					{
						TeamIndex: 0,
						ActionKey: 'o',
						Details:   []string{"f", "1", "0", "f", "lb"},
					},
				},
			},
			shouldError: false,
		},
	}

	builder := Builder{}
	for _, test := range tests {
		_, err := builder.Load(test.bgn)
		assert.Equal(t, test.shouldError, err != nil, test.name)
	}
}
