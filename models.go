package go_carcassonne

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
var StartTile = NewTile(City, Road, Farm, Road, NilStructure, false, false)

// Tiles are all the tiles in that will be placed
var Tiles = map[*tile]int{
	NewTile(Farm, Farm, Farm, Farm, Cloister, false, false):     4,
	NewTile(Farm, Farm, Road, Farm, Cloister, false, false):     2,
	NewTile(City, City, City, City, NilStructure, true, true):   1,
	NewTile(City, City, Farm, City, NilStructure, true, false):  3,
	NewTile(City, City, Farm, City, NilStructure, true, true):   1,
	NewTile(City, City, Road, City, NilStructure, true, false):  1,
	NewTile(City, City, Road, City, NilStructure, true, true):   2,
	NewTile(City, Farm, Farm, City, NilStructure, true, false):  3,
	NewTile(City, Farm, Farm, City, NilStructure, true, true):   2,
	NewTile(City, Road, Road, City, NilStructure, true, false):  3,
	NewTile(City, Road, Road, City, NilStructure, true, true):   2,
	NewTile(Farm, City, Farm, City, NilStructure, true, false):  1,
	NewTile(Farm, City, Farm, City, NilStructure, true, true):   2,
	NewTile(City, Farm, Farm, City, NilStructure, false, false): 2,
	NewTile(City, Farm, City, Farm, NilStructure, false, false): 3,
	NewTile(City, Farm, Farm, Farm, NilStructure, false, false): 5,
	NewTile(City, Farm, Road, Road, NilStructure, false, false): 3,
	NewTile(City, Road, Road, Farm, NilStructure, false, false): 3,
	NewTile(City, Road, Road, Road, NilStructure, false, false): 3,
	NewTile(City, Road, Farm, Road, NilStructure, false, false): 3,
	NewTile(Road, Farm, Road, Farm, NilStructure, false, false): 8,
	NewTile(Farm, Farm, Road, Road, NilStructure, false, false): 9,
	NewTile(Farm, Road, Road, Road, NilStructure, false, false): 4,
	NewTile(Road, Road, Road, Road, NilStructure, false, false): 1,
}
