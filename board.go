package go_carcassonne

import (
	"fmt"
)

// connection used in graph search algorithms
type connection struct {
	tile *tile
	side string
}

// board - +X right +Y up
type board struct {
	board          []*tile // list of all tiles in order of added
	completeCities []*structure
	completeRoads  []*structure
}

func newBoard() *board {
	start := startTile.copy()
	start.X, start.Y = 0, 0
	return &board{
		board: []*tile{start},
	}
}

func (b *board) Place(t *tile, x, y int) error {
	if b.tile(x, y) != nil {
		return fmt.Errorf("tile aready at %d,%d", x, y)
	}
	// get all adj tiles and check if valid placement
	sides := make(map[string]*tile)
	for _, boardTile := range b.board {
		if boardTile.X == x+1 && boardTile.Y == y {
			if boardTile.Sides[SideLeft] != t.Sides[SideRight] {
				return fmt.Errorf("invalid tile placement")
			}
			sides[SideRight] = boardTile
		} else if boardTile.X == x && boardTile.Y == y+1 {
			if boardTile.Sides[SideBottom] != t.Sides[SideTop] {
				return fmt.Errorf("invalid tile placement")
			}
			sides[SideTop] = boardTile
		} else if boardTile.X == x-1 && boardTile.Y == y {
			if boardTile.Sides[SideRight] != t.Sides[SideLeft] {
				return fmt.Errorf("invalid tile placement")
			}
			sides[SideLeft] = boardTile
		} else if boardTile.X == x && boardTile.Y == y-1 {
			if boardTile.Sides[SideTop] != t.Sides[SideBottom] {
				return fmt.Errorf("invalid tile placement")
			}
			sides[SideBottom] = boardTile
		}
	}
	if len(sides) <= 0 {
		return fmt.Errorf("cannot add a disconnected tile to the board")
	}
	// update adjacent pointer values
	if sides[SideTop] != nil {
		t.adjacent[SideTop] = sides[SideTop]
		sides[SideTop].adjacent[SideBottom] = t
	}
	if sides[SideRight] != nil {
		t.adjacent[SideRight] = sides[SideRight]
		sides[SideRight].adjacent[SideLeft] = t
	}
	if sides[SideBottom] != nil {
		t.adjacent[SideBottom] = sides[SideBottom]
		sides[SideBottom].adjacent[SideTop] = t
	}
	if sides[SideLeft] != nil {
		t.adjacent[SideLeft] = sides[SideLeft]
		sides[SideLeft].adjacent[SideRight] = t
	}
	// update x, y
	t.X, t.Y = x, y
	// add tile to list of tiles
	b.board = append(b.board, t)
	return nil
}

// get a tile at x,y or return nil
func (b *board) tile(x, y int) *tile {
	for _, tile := range b.board {
		if tile.X == x && tile.Y == y {
			return tile
		}
	}
	return nil
}

// given a tile location and side that contains a city section, get the current city structure
func (b *board) generateCity(x, y int, side string) (*structure, error) {
	tile := b.tile(x, y)
	if tile == nil {
		return nil, fmt.Errorf("tile does not exist at %d,%d", x, y)
	}
	if tile.Sides[side] != City {
		return nil, fmt.Errorf("side %s does not contain city section at tile %d,%d", side, x, y)
	}
	// perform BFS
	complete := true
	seen := make([]*node, 0)         // keeps track of tile and sides in city
	visited := make(map[string]bool) // key is "XY" of tile
	queue := make([]*connection, 0)

	queue = append(queue, &connection{
		tile: tile,
		side: side,
	})
	visited[fmt.Sprintf("%d%d", x, y)] = true
	for len(queue) > 0 {
		front := queue[0]
		queue = queue[1:]

		sides, err := front.tile.connectedCitySides(front.side)
		if err != nil {
			return nil, err
		}
		sides = append(sides, front.side)
		seenNode := &node{
			tile:  front.tile,
			sides: sides,
		}
		seen = append(seen, seenNode)
		for _, s := range sides {
			adjacentTile := front.tile.adjacent[s]
			if adjacentTile == nil {
				complete = false // city not yet completed
			} else if visited[fmt.Sprintf("%d%d", adjacentTile.X, adjacentTile.Y)] {
				// edge case for disconnected city sides that end up being part of the same city
				for _, n := range seen {
					if adjacentTile.X == n.tile.X && adjacentTile.Y == n.tile.Y {
						if !contains(n.sides, AcrossSide[s]) {
							n.sides = append(n.sides, AcrossSide[s])
						}
					}
				}
			} else if !visited[fmt.Sprintf("%d%d", adjacentTile.X, adjacentTile.Y)] {
				queue = append(queue, &connection{
					tile: adjacentTile,
					side: AcrossSide[s],
				})
				visited[fmt.Sprintf("%d%d", adjacentTile.X, adjacentTile.Y)] = true
			}
		}
		// continued edge case for disconnected city sides that end up being part of the same city
		if len(queue) == 0 && !front.tile.ConnectedCitySides {
			for side, section := range front.tile.Sides {
				if section == City {
					x := front.tile.X
					y := front.tile.Y
					if side == SideTop {
						y++
					} else if side == SideRight {
						x++
					} else if side == SideBottom {
						y--
					} else if side == SideLeft {
						x--
					}
					for _, n := range seen {
						if n.tile.X == x && n.tile.Y == y {
							if contains(n.sides, AcrossSide[side]) && !contains(seenNode.sides, side) {
								seenNode.sides = append(seenNode.sides, side)
							}
						}
					}
				}
			}
		}
	}
	return &structure{
		typ:      City,
		complete: complete,
		nodes:    seen,
	}, nil
}

// given a tile location and side that contains a road section, get the current road structure
func (b *board) generateRoad(x, y int, side string) (*structure, error) {
	tile := b.tile(x, y)
	if tile == nil {
		return nil, fmt.Errorf("tile does not exist at %d,%d", x, y)
	}
	if tile.Sides[side] != Road {
		return nil, fmt.Errorf("side %s does not contain road section at tile %d,%d", side, x, y)
	}
	// perform BFS
	complete := true
	seen := make([]*node, 0)         // keeps track of tile and sides in road
	visited := make(map[string]bool) // key is "XY" of tile
	queue := make([]*connection, 0)

	queue = append(queue, &connection{
		tile: tile,
		side: side,
	})
	visited[fmt.Sprintf("%d%d", x, y)] = true
	for len(queue) > 0 {
		front := queue[0]
		queue = queue[1:]

		sides, err := front.tile.connectedRoadSides(front.side)
		if err != nil {
			return nil, err
		}
		sides = append(sides, front.side)
		seenNode := &node{
			tile:  front.tile,
			sides: sides,
		}
		seen = append(seen, seenNode)
		for _, s := range sides {
			adjacentTile := front.tile.adjacent[s]
			if adjacentTile == nil {
				complete = false // road not yet completed
			} else if visited[fmt.Sprintf("%d%d", adjacentTile.X, adjacentTile.Y)] {
				// edge case for disconnected road sides that end up being part of the same road
				for _, n := range seen {
					if adjacentTile.X == n.tile.X && adjacentTile.Y == n.tile.Y {
						if !contains(n.sides, AcrossSide[s]) {
							n.sides = append(n.sides, AcrossSide[s])
						}
					}
				}
			} else if !visited[fmt.Sprintf("%d%d", adjacentTile.X, adjacentTile.Y)] {
				queue = append(queue, &connection{
					tile: adjacentTile,
					side: AcrossSide[s],
				})
				visited[fmt.Sprintf("%d%d", adjacentTile.X, adjacentTile.Y)] = true
			}
		}
		// continued edge case for disconnected road sides that end up being part of the same road
		if len(queue) == 0 {
			for side, section := range front.tile.Sides {
				if section == Road {
					x := front.tile.X
					y := front.tile.Y
					if side == SideTop {
						y++
					} else if side == SideRight {
						x++
					} else if side == SideBottom {
						y--
					} else if side == SideLeft {
						x--
					}
					for _, n := range seen {
						if n.tile.X == x && n.tile.Y == y {
							if contains(n.sides, AcrossSide[side]) && !contains(seenNode.sides, side) {
								seenNode.sides = append(seenNode.sides, side)
							}
						}
					}
				}
			}
		}
	}
	return &structure{
		typ:      Road,
		complete: complete,
		nodes:    seen,
	}, nil
}

// given a tile location and farmSide that contains farmland i.e. farm or road section, get the current farm structure
func (b *board) generateFarm(x, y int, farmSide string) (*structure, error) {
	tile := b.tile(x, y)
	if tile == nil {
		return nil, fmt.Errorf("tile does not exist at %d,%d", x, y)
	}
	side := farmSideToSide(farmSide)
	if tile.Sides[side] != Farm && tile.Sides[side] != Road {
		return nil, fmt.Errorf("side %s does not contain farmland at tile %d,%d", side, x, y)
	}
	// start BFS
	seen := make([]*node, 0)         // keeps track of tile and farm sides in farm
	visited := make(map[string]bool) // key is "XY" of tile
	queue := make([]*connection, 0)

	queue = append(queue, &connection{
		tile: tile,
		side: farmSide,
	})
	visited[fmt.Sprintf("%d%d%s", x, y, farmSide)] = true

	for len(queue) > 0 {
		front := queue[0]
		queue = queue[1:]

		sides, err := front.tile.connectedFarmSides(front.side)
		if err != nil {
			return nil, err
		}
		sides = append(sides, front.side)
		// if tile not already in seen add it
		found := false
		for _, s := range seen {
			if fmt.Sprintf("%d%d", s.tile.X, s.tile.Y) == fmt.Sprintf("%d%d", front.tile.X, front.tile.Y) {
				// check if all sides have been found and if not add those sides
				for _, side := range sides {
					if !contains(s.sides, side) {
						s.sides = append(s.sides, side)
					}
				}
				found = true
				break
			}
		}
		if !found {
			seen = append(seen, &node{
				tile:  front.tile,
				sides: sides,
			})
		}
		// perform BFS
		for _, farmSide := range sides {
			side := farmSideToSide(farmSide)
			adjacentTile := front.tile.adjacent[side]
			if adjacentTile != nil && !visited[fmt.Sprintf("%d%d%s", adjacentTile.X, adjacentTile.Y, farmSide)] {
				queue = append(queue, &connection{
					tile: adjacentTile,
					side: AcrossFarmSide[farmSide],
				})
				visited[fmt.Sprintf("%d%d%s", adjacentTile.X, adjacentTile.Y, farmSide)] = true
			}
		}
	}
	return &structure{
		typ:   Farm,
		nodes: seen,
	}, nil
}

// given a tile location that contains a cloister, get the number of tiles surrounding the cloister
func (b *board) tilesSurroundingCloister(x, y int) (int, error) {
	t := b.tile(x, y)
	if t == nil {
		return 0, fmt.Errorf("tile does not exist at %d,%d", x, y)
	}
	count := 0
	locations := [][]int{{x + 1, y}, {x - 1, y}, {x, y + 1}, {x, y - 1}, {x + 1, y + 1}, {x - 1, y + 1}, {x + 1, y - 1}, {x - 1, y - 1}}
	for _, location := range locations {
		if b.tile(location[0], location[1]) != nil {
			count++
		}
	}
	return count, nil
}

func (b *board) playable(t *tile) bool {
	emptySpaces := b.getEmptySpaces()
	// go through empty spaces looking for at least one place the tile can be placed
	copied := t.copy()
	for i := 0; i < 4; i++ {
		for _, emptySpace := range emptySpaces {
			valid := true
			for _, side := range Sides {
				if emptySpace.adjacent[side] != nil && emptySpace.adjacent[side].Sides[AcrossSide[side]] != copied.Sides[side] {
					valid = false
				}
			}
			if valid {
				return true
			}
		}
		copied.RotateRight()
	}
	return false
}

// gets a list of empty spaces that are potential place locations
func (b *board) getEmptySpaces() []*tile {
	// go through board and get all empty spaces
	emptySpaces := make(map[string]*tile)
	for _, boardTile := range b.board {
		for _, side := range Sides {
			x := boardTile.X
			y := boardTile.Y
			if side == SideTop {
				y++
			} else if side == SideRight {
				x++
			} else if side == SideBottom {
				y--
			} else if side == SideLeft {
				x--
			}
			if boardTile.adjacent[side] == nil && emptySpaces[fmt.Sprintf("%d%d", x, y)] == nil {
				emptySpaces[fmt.Sprintf("%d%d", x, y)] = emptySpaceTile(x, y)
			}
		}
	}
	// go through board and get all tiles adjacent to empty spaces
	for _, boardTile := range b.board {
		for _, side := range Sides {
			x := boardTile.X
			y := boardTile.Y
			if side == SideTop {
				y++
			} else if side == SideRight {
				x++
			} else if side == SideBottom {
				y--
			} else if side == SideLeft {
				x--
			}
			if boardTile.adjacent[side] == nil {
				emptySpaces[fmt.Sprintf("%d%d", x, y)].adjacent[AcrossSide[side]] = boardTile
			}
		}
	}
	result := make([]*tile, 0)
	for _, emptySpace := range emptySpaces {
		result = append(result, emptySpace)
	}
	return result
}
