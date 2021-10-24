package go_carcassonne

// Action types
const (
	PlaceTile   = "PlaceTile"
	PlaceToken  = "PlaceToken"
	RotateRight = "RotateRight"
	RotateLeft  = "RotateLeft"
	Reset       = "Reset"
)

// PlaceTileActionDetails is the action details for placing a tile
type PlaceTileActionDetails struct {
	// X and Y location where to place the tile
	X, Y int

	// Tile is the tile being placed
	// Not sent by user - set in server for documentation only
	Tile *tile
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

// CarcassonneSnapshotDetails are the details unique to carcassonne
type CarcassonneSnapshotDetails struct {
	PlayTile       *tile
	LastPlacedTile *tile
	Board          *board
	BoardTokens    []*token
	Tokens         map[string]int
	Scores         map[string]int
	TilesRemaining int
}

// StartTile the tile at 0,0 at the start of the game
var StartTile = newTile(City, Road, Farm, Road, NilStructure, false, false)

// Tiles are all the tiles in that will be placed
var Tiles = map[*tile]int{
	newTile(Farm, Farm, Farm, Farm, Cloister, false, false):     4,
	newTile(Farm, Farm, Road, Farm, Cloister, false, false):     2,
	newTile(City, City, City, City, NilStructure, true, true):   1,
	newTile(City, City, Farm, City, NilStructure, true, false):  3,
	newTile(City, City, Farm, City, NilStructure, true, true):   1,
	newTile(City, City, Road, City, NilStructure, true, false):  1,
	newTile(City, City, Road, City, NilStructure, true, true):   2,
	newTile(City, Farm, Farm, City, NilStructure, true, false):  3,
	newTile(City, Farm, Farm, City, NilStructure, true, true):   2,
	newTile(City, Road, Road, City, NilStructure, true, false):  3,
	newTile(City, Road, Road, City, NilStructure, true, true):   2,
	newTile(Farm, City, Farm, City, NilStructure, true, false):  1,
	newTile(Farm, City, Farm, City, NilStructure, true, true):   2,
	newTile(City, Farm, Farm, City, NilStructure, false, false): 2,
	newTile(City, Farm, City, Farm, NilStructure, false, false): 3,
	newTile(City, Farm, Farm, Farm, NilStructure, false, false): 5,
	newTile(City, Farm, Road, Road, NilStructure, false, false): 3,
	newTile(City, Road, Road, Farm, NilStructure, false, false): 3,
	newTile(City, Road, Road, Road, NilStructure, false, false): 3,
	newTile(City, Road, Farm, Road, NilStructure, false, false): 3,
	newTile(Road, Farm, Road, Farm, NilStructure, false, false): 8,
	newTile(Farm, Farm, Road, Road, NilStructure, false, false): 9,
	newTile(Farm, Road, Road, Road, NilStructure, false, false): 4,
	newTile(Road, Road, Road, Road, NilStructure, false, false): 1,
}
