package go_carcassonne

import (
	"fmt"
	"strconv"

	bg "github.com/quibbble/go-boardgame"
	"github.com/quibbble/go-boardgame/pkg/bgerr"
)

var (
	actionToNotation = map[string]string{ActionPlaceTile: "i", ActionPlaceToken: "o", ActionRotateTileRight: "r", ActionRotateTileLeft: "l", bg.ActionSetWinners: "w"}
	notationToAction = reverseMap(actionToNotation)

	sideToNotation = map[string]string{SideTop: "t", SideRight: "r", SideBottom: "b", SideLeft: "l"}
	notationToSide = reverseMap(sideToNotation)

	farmSideToNotation = map[string]string{FarmSideTopA: "ta", FarmSideTopB: "tb", FarmSideRightA: "ra", FarmSideRightB: "rb", FarmSideBottomA: "ba", FarmSideBottomB: "bb", FarmSideLeftA: "la", FarmSideLeftB: "lb"}
	notationToFarmSide = reverseMap(farmSideToNotation)

	structureToNotation = map[string]string{Road: "r", Farm: "f", City: "c", Cloister: "m", NilStructure: "n"}
	notationToStructure = reverseMap(structureToNotation)

	tokenToNotation = map[string]string{Farmer: "f", Knight: "k", Thief: "t", Monk: "m"}
	notationToToken = reverseMap(tokenToNotation)

	boolToNotation = map[bool]string{true: "t", false: "f"}
	notationToBool = map[string]bool{"t": true, "f": false}
)

func (p *PlaceTileActionDetails) encodeBGN() []string {
	return []string{
		strconv.Itoa(p.X), strconv.Itoa(p.Y),
		structureToNotation[p.Tile.Top], structureToNotation[p.Tile.Right], structureToNotation[p.Tile.Bottom], structureToNotation[p.Tile.Left],
		structureToNotation[p.Tile.Center],
		boolToNotation[p.Tile.ConnectedCitySides], boolToNotation[p.Tile.Banner],
	}
}

func decodePlaceTileActionDetailsBGN(notation []string) (*PlaceTileActionDetails, error) {
	if len(notation) != 9 {
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
	top, ok := notationToStructure[notation[2]]
	if !ok {
		return nil, loadFailure(fmt.Errorf("failed to get top of tile"))
	}
	right, ok := notationToStructure[notation[3]]
	if !ok {
		return nil, loadFailure(fmt.Errorf("failed to get right of tile"))
	}
	bottom, ok := notationToStructure[notation[4]]
	if !ok {
		return nil, loadFailure(fmt.Errorf("failed to get bottom of tile"))
	}
	left, ok := notationToStructure[notation[5]]
	if !ok {
		return nil, loadFailure(fmt.Errorf("failed to get left of tile"))
	}
	center, ok := notationToStructure[notation[6]]
	if !ok {
		return nil, loadFailure(fmt.Errorf("failed to get center of tile"))
	}
	connectedCitySides, ok := notationToBool[notation[7]]
	if !ok {
		return nil, loadFailure(fmt.Errorf("failed to get connectedCitySides"))
	}
	banner, ok := notationToBool[notation[8]]
	if !ok {
		return nil, loadFailure(fmt.Errorf("failed to get banner"))
	}
	return &PlaceTileActionDetails{
		X: x,
		Y: y,
		Tile: TileActionDetails{
			top, right, bottom, left, center, connectedCitySides, banner,
		},
	}, nil
}

func (p *PlaceTokenActionDetails) encodeBGN() []string {
	if p.Pass {
		return []string{boolToNotation[p.Pass]}
	} else if p.Type == Monk {
		return []string{boolToNotation[p.Pass], strconv.Itoa(p.X), strconv.Itoa(p.Y), tokenToNotation[p.Type]}
	} else if p.Type == Farmer {
		return []string{boolToNotation[p.Pass], strconv.Itoa(p.X), strconv.Itoa(p.Y), tokenToNotation[p.Type], farmSideToNotation[p.Side]}
	}
	return []string{boolToNotation[p.Pass], strconv.Itoa(p.X), strconv.Itoa(p.Y), tokenToNotation[p.Type], sideToNotation[p.Side]}
}

func decodePlaceTokenActionDetailsBGN(notation []string) (*PlaceTokenActionDetails, error) {
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
		Status: bgerr.StatusBGNDecodingFailure,
	}
}
