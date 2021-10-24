package go_carcassonne

import (
	bg "github.com/quibbble/go-boardgame"
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	TeamA = "TeamA"
	TeamB = "TeamB"
)

func Test_Carcassonne(t *testing.T) {
	Carcassonne := NewCarcassonne(bg.BoardGameOptions{Teams: []string{TeamA, TeamB}})

	assert.Equal(t, len(Carcassonne.state.board.board), 1, "board missing start tile")

	// place all road tile to right of start tile
	Carcassonne.state.playTile = NewTile(Road, Road, Road, Road, NilStructure, false, false)
	Carcassonne.state.turn = TeamA
	err := Carcassonne.Do(bg.BoardGameAction{
		Team:       TeamA,
		ActionType: PlaceTile,
		MoreDetails: PlaceTileActionDetails{
			X: 1,
			Y: 0,
		},
	})
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	assert.Equal(t, len(Carcassonne.state.board.board), 2, "board missing last placed tile")

	// place token in left B side of farmland
	err = Carcassonne.Do(bg.BoardGameAction{
		Team:       TeamA,
		ActionType: PlaceToken,
		MoreDetails: PlaceTokenActionDetails{
			Pass: false,
			X:    1,
			Y:    0,
			Type: Farmer,
			Side: FarmSideLeftB,
		},
	})
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	assert.Equal(t, len(Carcassonne.state.boardTokens), 1, "missing placed token")
	assert.Equal(t, Carcassonne.state.boardTokens[0].Type, Farmer, "incorrect token placed")
	assert.Equal(t, Carcassonne.state.turn, TeamB, "incorrect team's turn")

	// place tile to left of start tile completing a road segment
	Carcassonne.state.playTile = NewTile(Road, Road, Farm, Road, NilStructure, false, false)
	err = Carcassonne.Do(bg.BoardGameAction{
		Team:       TeamB,
		ActionType: PlaceTile,
		MoreDetails: PlaceTileActionDetails{
			X: -1,
			Y: 0,
		},
	})
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	assert.Equal(t, len(Carcassonne.state.board.board), 3, "board missing last place tile")

	// claim the completed road by placing thief on right side of tile
	err = Carcassonne.Do(bg.BoardGameAction{
		Team:       TeamB,
		ActionType: PlaceToken,
		MoreDetails: PlaceTokenActionDetails{
			Pass: false,
			X:    -1,
			Y:    0,
			Type: Thief,
			Side: SideRight,
		},
	})
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	assert.Equal(t, len(Carcassonne.state.boardTokens), 1)
	assert.Equal(t, Carcassonne.state.scores[TeamB], 3)
	assert.Equal(t, Carcassonne.state.board.board[0].Teams[SideLeft], []string{TeamB})
	assert.Equal(t, Carcassonne.state.board.board[0].Teams[SideRight], []string{TeamB})
	assert.Equal(t, Carcassonne.state.board.board[1].Teams[SideLeft], []string{TeamB})
	assert.Equal(t, Carcassonne.state.board.board[2].Teams[SideRight], []string{TeamB})
	assert.Equal(t, Carcassonne.state.turn, TeamA)
}
