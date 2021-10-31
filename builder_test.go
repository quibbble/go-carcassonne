package go_carcassonne

import (
	"encoding/json"
	bg "github.com/quibbble/go-boardgame"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Builder_Notation(t *testing.T) {
	builder := Builder{}
	teams := []string{TeamA, TeamB}
	carcassonne, err := builder.CreateAdvanced(&bg.BoardGameOptions{Teams: teams, Seed: 123})
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

	notation := carcassonne.GetNotation()
	carcassonneLoaded, err := builder.Load(teams, notation)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	expected, _ := json.Marshal(carcassonne)
	actual, _ := json.Marshal(carcassonneLoaded)
	assert.Equal(t, expected, actual)
}

func Test_Builder_Notations(t *testing.T) {
	tests := []struct {
		name        string
		notation    string
		shouldError bool
	}{
		{
			name:        "empty string should error",
			notation:    "",
			shouldError: true,
		},
		{
			name:        "empty values should error",
			notation:    "::",
			shouldError: true,
		},
		{
			name:        "only number of teams should error",
			notation:    "2::",
			shouldError: true,
		},
		{
			name:        "only seed should error",
			notation:    ":123:",
			shouldError: true,
		},
		{
			name:        "number of teams string should error",
			notation:    "A:123:",
			shouldError: true,
		},
		{
			name:        "seed string should error",
			notation:    "2:A:",
			shouldError: true,
		},
		{
			name:        "game with no actions should succeed",
			notation:    "2:123:",
			shouldError: false,
		},
		{
			name:        "game with incorrect action should error",
			notation:    "2:123:A;",
			shouldError: true,
		},
		{
			name:        "game with invalid player index should error",
			notation:    "2:123:9,2;",
			shouldError: true,
		},
		{
			name:        "game with invalid action number should error",
			notation:    "2:123:0,5;",
			shouldError: true,
		},
		{
			name:        "game with rotate tile right action should succeed",
			notation:    "2:123:0,2;",
			shouldError: false,
		},
	}

	builder := Builder{}
	teams := []string{TeamA, TeamB}
	for _, test := range tests {
		_, err := builder.Load(teams, test.notation)
		assert.Equal(t, test.shouldError, err != nil, test.name)
	}
}
