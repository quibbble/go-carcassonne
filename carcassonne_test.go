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
	carcassonne, err := NewCarcassonne(bg.BoardGameOptions{Teams: []string{TeamA, TeamB}})
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	assert.Equal(t, len(carcassonne.state.board.board), 1, "board missing start tile")

	// place all road tile to right of start tile
	carcassonne.state.playTile = newTile(Road, Road, Road, Road, NilStructure, false, false)
	carcassonne.state.turn = TeamA
	err = carcassonne.Do(bg.BoardGameAction{
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

	assert.Equal(t, len(carcassonne.state.board.board), 2, "board missing last placed tile")

	// place token in left B side of farmland
	err = carcassonne.Do(bg.BoardGameAction{
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

	assert.Equal(t, len(carcassonne.state.boardTokens), 1, "missing placed token")
	assert.Equal(t, carcassonne.state.boardTokens[0].Type, Farmer, "incorrect token placed")
	assert.Equal(t, carcassonne.state.turn, TeamB, "incorrect team's turn")

	// place tile to left of start tile completing a road segment
	carcassonne.state.playTile = newTile(Road, Road, Farm, Road, NilStructure, false, false)
	err = carcassonne.Do(bg.BoardGameAction{
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

	assert.Equal(t, len(carcassonne.state.board.board), 3, "board missing last place tile")

	// claim the completed road by placing thief on right side of tile
	err = carcassonne.Do(bg.BoardGameAction{
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

	assert.Equal(t, len(carcassonne.state.boardTokens), 1)
	assert.Equal(t, carcassonne.state.scores[TeamB], 3)
	assert.Equal(t, carcassonne.state.board.board[0].Teams[SideLeft], []string{TeamB})
	assert.Equal(t, carcassonne.state.board.board[0].Teams[SideRight], []string{TeamB})
	assert.Equal(t, carcassonne.state.board.board[1].Teams[SideLeft], []string{TeamB})
	assert.Equal(t, carcassonne.state.board.board[2].Teams[SideRight], []string{TeamB})
	assert.Equal(t, carcassonne.state.turn, TeamA)
}
