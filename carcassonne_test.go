package go_carcassonne

import (
	bg "github.com/quibbble/go-boardgame"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

const (
	TeamA = "TeamA"
	TeamB = "TeamB"
)

func Test_Carcassonne(t *testing.T) {
	carcassonne, err := NewCarcassonne(&bg.BoardGameOptions{
		Teams: []string{TeamA, TeamB},
		MoreOptions: CarcassonneMoreOptions{
			Seed: time.Now().UnixNano(),
		},
	})
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	assert.Equal(t, 1, len(carcassonne.state.board.board), "board missing start tile")

	// place all road tile to right of start tile
	carcassonne.state.playTile = newTile(Road, Road, Road, Road, NilStructure, false, false)
	carcassonne.state.turn = TeamA
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

	assert.Equal(t, 2, len(carcassonne.state.board.board), "board missing last placed tile")

	// place token in left B side of farmland
	err = carcassonne.Do(&bg.BoardGameAction{
		Team:       TeamA,
		ActionType: ActionPlaceToken,
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

	assert.Equal(t, 1, len(carcassonne.state.boardTokens), "missing placed token")
	assert.Equal(t, Farmer, carcassonne.state.boardTokens[0].Type, "incorrect token placed")
	assert.Equal(t, TeamB, carcassonne.state.turn, "incorrect team's turn")

	// place tile to left of start tile completing a road segment
	carcassonne.state.playTile = newTile(Road, Road, Farm, Road, NilStructure, false, false)
	err = carcassonne.Do(&bg.BoardGameAction{
		Team:       TeamB,
		ActionType: ActionPlaceTile,
		MoreDetails: PlaceTileActionDetails{
			X: -1,
			Y: 0,
		},
	})
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	assert.Equal(t, 3, len(carcassonne.state.board.board), "board missing last place tile")

	// claim the completed road by placing thief on right side of tile
	err = carcassonne.Do(&bg.BoardGameAction{
		Team:       TeamB,
		ActionType: ActionPlaceToken,
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

	assert.Equal(t, 1, len(carcassonne.state.boardTokens))
	assert.Equal(t, 3, carcassonne.state.scores[TeamB])
	assert.Equal(t, []string{TeamB}, carcassonne.state.board.board[0].Teams[SideLeft])
	assert.Equal(t, []string{TeamB}, carcassonne.state.board.board[0].Teams[SideRight])
	assert.Equal(t, []string{TeamB}, carcassonne.state.board.board[1].Teams[SideLeft])
	assert.Equal(t, []string{TeamB}, carcassonne.state.board.board[2].Teams[SideRight])
	assert.Equal(t, TeamA, carcassonne.state.turn)
}
