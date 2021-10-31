package go_carcassonne

import (
	"fmt"
	"github.com/quibbble/go-boardgame/pkg/bgerr"
	"strconv"
	"strings"
)

var (
	notationActionToInt    = map[string]int{ActionPlaceTile: 0, ActionPlaceToken: 1, ActionRotateTileRight: 2, ActionRotateTileLeft: 3}
	notationIntToAction    = map[string]string{"0": ActionPlaceTile, "1": ActionPlaceToken, "2": ActionRotateTileRight, "3": ActionRotateTileLeft}
	notationSideToInt      = map[string]int{SideTop: 0, SideRight: 1, SideBottom: 2, SideLeft: 3}
	notationIntToSide      = map[string]string{"0": SideTop, "1": SideRight, "2": SideBottom, "3": SideLeft}
	notationFarmSideToInt  = map[string]int{FarmSideTopA: 0, FarmSideTopB: 1, FarmSideRightA: 2, FarmSideRightB: 3, FarmSideBottomA: 4, FarmSideBottomB: 5, FarmSideLeftA: 6, FarmSideLeftB: 7}
	notationIntToFarmSide  = map[string]string{"0": FarmSideTopA, "1": FarmSideTopB, "2": FarmSideRightA, "3": FarmSideRightB, "4": FarmSideBottomA, "5": FarmSideBottomB, "6": FarmSideLeftA, "7": FarmSideLeftB}
	notationTokenTypeToInt = map[string]int{Farmer: 0, Knight: 1, Thief: 2, Monk: 3}
	notationIntToTokenType = map[string]string{"0": Farmer, "1": Knight, "2": Thief, "3": Monk}
)

func (p *PlaceTileActionDetails) encode() string {
	return fmt.Sprintf("%d,%d", p.X, p.Y)
}

func decodeNotationPlaceTileActionDetails(notation string) (*PlaceTileActionDetails, error) {
	split := strings.Split(notation, ",")
	if len(split) != 2 {
		return nil, loadFailure(fmt.Errorf("got %d but wanted %d fields in when decoding %s details", len(split), 2, ActionPlaceTile))
	}
	x, err := strconv.Atoi(split[0])
	if err != nil {
		return nil, loadFailure(err)
	}
	y, err := strconv.Atoi(split[1])
	if err != nil {
		return nil, loadFailure(err)
	}
	return &PlaceTileActionDetails{
		X: x,
		Y: y,
	}, nil
}

func (p *PlaceTokenActionDetails) encode() string {
	if p.Pass {
		return fmt.Sprintf("%d", boolToInt(p.Pass))
	} else if p.Type == Monk {
		return fmt.Sprintf("%d,%d,%d,%d", boolToInt(p.Pass), p.X, p.Y, notationTokenTypeToInt[p.Type])
	} else if p.Type == Farmer {
		return fmt.Sprintf("%d,%d,%d,%d,%d", boolToInt(p.Pass), p.X, p.Y, notationTokenTypeToInt[p.Type], notationFarmSideToInt[p.Side])
	}
	return fmt.Sprintf("%d,%d,%d,%d,%d", boolToInt(p.Pass), p.X, p.Y, notationTokenTypeToInt[p.Type], notationSideToInt[p.Side])
}

func decodeNotationPlaceTokenActionDetails(notation string) (*PlaceTokenActionDetails, error) {
	split := strings.Split(notation, ",")
	if len(split) < 1 || len(split) > 5 {
		return nil, loadFailure(fmt.Errorf("got %d but wanted %d to %d fields in when decoding %s details", len(split), 1, 5, ActionPlaceToken))
	}
	var pass bool
	if split[0] != "0" {
		pass = false
	} else if split[0] != "1" {
		pass = true
	} else {
		return nil, loadFailure(fmt.Errorf("got %s but wanted be 0 or 1 for for Pass when decoding %s details", split[0], ActionPlaceToken))
	}
	if pass {
		return &PlaceTokenActionDetails{Pass: pass}, nil
	}
	x, err := strconv.Atoi(split[1])
	if err != nil {
		return nil, loadFailure(err)
	}
	y, err := strconv.Atoi(split[2])
	if err != nil {
		return nil, loadFailure(err)
	}
	token := notationIntToTokenType[split[3]]
	if len(split) == 4 {
		return &PlaceTokenActionDetails{Pass: pass, X: x, Y: y, Type: token}, nil
	}
	side := notationIntToSide[split[4]]
	if token == Farmer {
		side = notationIntToFarmSide[split[4]]
	}
	return &PlaceTokenActionDetails{Pass: pass, X: x, Y: y, Type: token, Side: side}, nil
}

func loadFailure(err error) error {
	return &bgerr.Error{
		Err:    err,
		Status: bgerr.StatusGameLoadFailure,
	}
}
