package go_carcassonne

// Action types
const (
	ActionPlaceTile       = "PlaceTile"
	ActionPlaceToken      = "PlaceToken"
	ActionRotateTileRight = "RotateTileRight"
	ActionRotateTileLeft  = "RotateTileLeft"
)

// CarcassonneMoreOptions are the additional options for creating a game of Carcassonne
type CarcassonneMoreOptions struct {
	Seed int64
}

// PlaceTileActionDetails is the action details for placing a tile
type PlaceTileActionDetails struct {
	// X and Y location where to place the tile
	X, Y int

	// Tile is the tile being placed
	Tile TileActionDetails
}

type TileActionDetails struct {
	Top, Right, Bottom, Left, Center string
	ConnectedCitySides, Banner       bool
}

// PlaceTokenActionDetails is the action details for placing a token
type PlaceTokenActionDetails struct {
	// Pass set to pass placing a token
	Pass bool

	// X and Y location where to place the token
	X, Y int

	// Type is the type of token to place
	Type string

	// Side is the side to place
	Side string
}

// CarcassonneSnapshotData is the game data unique to Carcassonne
type CarcassonneSnapshotData struct {
	PlayTile       *tile
	LastPlacedTile *tile
	Board          []*tile
	BoardTokens    []*token
	Tokens         map[string]int
	Scores         map[string]int
	TilesRemaining int
}

// startTile the tile at 0,0 at the start of the game
var startTile = newTile(City, Road, Farm, Road, NilStructure, false, false)

type tileAmounts struct {
	tile   *tile
	amount int
}

// tiles are all the tiles that will be placed
var tiles = []*tileAmounts{
	{tile: newTile(Farm, Farm, Farm, Farm, Cloister, false, false), amount: 4},
	{tile: newTile(Farm, Farm, Road, Farm, Cloister, false, false), amount: 2},
	{tile: newTile(City, City, City, City, NilStructure, true, true), amount: 1},
	{tile: newTile(City, City, Farm, City, NilStructure, true, false), amount: 3},
	{tile: newTile(City, City, Farm, City, NilStructure, true, true), amount: 1},
	{tile: newTile(City, City, Road, City, NilStructure, true, false), amount: 1},
	{tile: newTile(City, City, Road, City, NilStructure, true, true), amount: 2},
	{tile: newTile(City, Farm, Farm, City, NilStructure, true, false), amount: 3},
	{tile: newTile(City, Farm, Farm, City, NilStructure, true, true), amount: 2},
	{tile: newTile(City, Road, Road, City, NilStructure, true, false), amount: 3},
	{tile: newTile(City, Road, Road, City, NilStructure, true, true), amount: 2},
	{tile: newTile(Farm, City, Farm, City, NilStructure, true, false), amount: 1},
	{tile: newTile(Farm, City, Farm, City, NilStructure, true, true), amount: 2},
	{tile: newTile(City, Farm, Farm, City, NilStructure, false, false), amount: 2},
	{tile: newTile(City, Farm, City, Farm, NilStructure, false, false), amount: 3},
	{tile: newTile(City, Farm, Farm, Farm, NilStructure, false, false), amount: 5},
	{tile: newTile(City, Farm, Road, Road, NilStructure, false, false), amount: 3},
	{tile: newTile(City, Road, Road, Farm, NilStructure, false, false), amount: 3},
	{tile: newTile(City, Road, Road, Road, NilStructure, false, false), amount: 3},
	{tile: newTile(City, Road, Farm, Road, NilStructure, false, false), amount: 3},
	{tile: newTile(Road, Farm, Road, Farm, NilStructure, false, false), amount: 8},
	{tile: newTile(Farm, Farm, Road, Road, NilStructure, false, false), amount: 9},
	{tile: newTile(Farm, Road, Road, Road, NilStructure, false, false), amount: 4},
	{tile: newTile(Road, Road, Road, Road, NilStructure, false, false), amount: 1},
}
