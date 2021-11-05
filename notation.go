package go_carcassonne

import (
	"fmt"
	"github.com/quibbble/go-boardgame/pkg/bgerr"
	"strconv"
)

var (
	actionToNotation   = map[string]string{ActionPlaceTile: "i", ActionPlaceToken: "o", ActionRotateTileRight: "r", ActionRotateTileLeft: "l"}
	notationToAction   = reverseMap(actionToNotation)
	sideToNotation     = map[string]string{SideTop: "t", SideRight: "r", SideBottom: "b", SideLeft: "l"}
	notationToSide     = reverseMap(sideToNotation)
	farmSideToNotation = map[string]string{FarmSideTopA: "ta", FarmSideTopB: "tb", FarmSideRightA: "ra", FarmSideRightB: "rb", FarmSideBottomA: "ba", FarmSideBottomB: "bb", FarmSideLeftA: "la", FarmSideLeftB: "lb"}
	notationToFarmSide = reverseMap(farmSideToNotation)
	tokenToNotation    = map[string]string{Farmer: "f", Knight: "k", Thief: "t", Monk: "m"}
	notationToToken    = reverseMap(tokenToNotation)
	boolToNotation     = map[bool]string{true: "t", false: "f"}
	notationToBool     = map[string]bool{"t": true, "f": false}
)

func (p *PlaceTileActionDetails) encode() []string {
	return []string{fmt.Sprintf("%d", p.X), fmt.Sprintf("%d", p.Y)}
}

func decodePlaceTileActionDetails(notation []string) (*PlaceTileActionDetails, error) {
	if len(notation) != 2 {
		return nil, loadFailure(fmt.Errorf("got %d but wanted %d fields in when decoding %s details", len(notation), 2, ActionPlaceTile))
	}
	x, err := strconv.Atoi(notation[0])
	if err != nil {
		return nil, loadFailure(err)
	}
	y, err := strconv.Atoi(notation[1])
	if err != nil {
		return nil, loadFailure(err)
	}
	return &PlaceTileActionDetails{
		X: x,
		Y: y,
	}, nil
}

func (p *PlaceTokenActionDetails) encode() []string {
	if p.Pass {
		return []string{fmt.Sprintf("%s", boolToNotation[p.Pass])}
	} else if p.Type == Monk {
		return []string{fmt.Sprintf("%s", boolToNotation[p.Pass]), fmt.Sprintf("%d", p.X), fmt.Sprintf("%d", p.Y), tokenToNotation[p.Type]}
	} else if p.Type == Farmer {
		return []string{fmt.Sprintf("%s", boolToNotation[p.Pass]), fmt.Sprintf("%d", p.X), fmt.Sprintf("%d", p.Y), tokenToNotation[p.Type], farmSideToNotation[p.Side]}
	}
	return []string{fmt.Sprintf("%s", boolToNotation[p.Pass]), fmt.Sprintf("%d", p.X), fmt.Sprintf("%d", p.Y), tokenToNotation[p.Type], sideToNotation[p.Side]}
}

func decodePlaceTokenActionDetails(notation []string) (*PlaceTokenActionDetails, error) {
	if len(notation) < 1 || len(notation) > 5 {
		return nil, loadFailure(fmt.Errorf("got %d but wanted %d to %d fields in when decoding %s details", len(notation), 1, 5, ActionPlaceToken))
	}
	pass, ok := notationToBool[notation[0]]
	if !ok {
		return nil, loadFailure(fmt.Errorf("got %s but wanted be 0 or 1 for for Pass when decoding %s details", notation[0], ActionPlaceToken))
	}
	if pass {
		return &PlaceTokenActionDetails{Pass: pass}, nil
	}
	x, err := strconv.Atoi(notation[1])
	if err != nil {
		return nil, loadFailure(err)
	}
	y, err := strconv.Atoi(notation[2])
	if err != nil {
		return nil, loadFailure(err)
	}
	token := notationToToken[notation[3]]
	if len(notation) == 4 {
		return &PlaceTokenActionDetails{Pass: pass, X: x, Y: y, Type: token}, nil
	}
	side := notationToSide[notation[4]]
	if token == Farmer {
		side = notationToFarmSide[notation[4]]
	}
	return &PlaceTokenActionDetails{Pass: pass, X: x, Y: y, Type: token, Side: side}, nil
}

func loadFailure(err error) error {
	return &bgerr.Error{
		Err:    err,
		Status: bgerr.StatusGameLoadFailure,
	}
}
