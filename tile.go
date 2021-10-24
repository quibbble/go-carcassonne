package go_carcassonne

import "fmt"

const OutOfBounds = 999

type tile struct {
	// X and Y represent the location of the tile on the board, OutOfBounds means not in board
	X, Y int
	// Sides is a map from side i.e. Top, Right, Bottom, Left to structure type i.e. Farm, City, Road
	Sides map[string]string
	// Center represents the structure at the center, i.e. Cloister, NilStructure
	Center string
	// ConnectedCitySides determines whether the city sides are connected together or separate
	ConnectedCitySides bool
	// Banner determines whether this is a banner tile
	Banner bool
	// Teams is a map from side to list of teams that have won that side after completing the given structure
	Teams map[string][]string
	// FarmTeams is a map from farm side to list of teams that have won that farmland at the end of the game
	FarmTeams map[string][]string
	// CenterTeam represents the team that has won the structure in the center of the tile after completing the given structure
	CenterTeam string
	// adjacent is a map from side to adjacent tiles
	adjacent map[string]*tile
}

func newTile(topStructure, rightStructure, bottomStructure, leftStructure, centerStructure string, connectedCitySides, banner bool) *tile {
	return &tile{
		X:                  OutOfBounds,
		Y:                  OutOfBounds,
		Sides:              map[string]string{SideTop: topStructure, SideRight: rightStructure, SideBottom: bottomStructure, SideLeft: leftStructure},
		Center:             centerStructure,
		ConnectedCitySides: connectedCitySides,
		Banner:             banner,
		Teams:              make(map[string][]string),
		FarmTeams:          make(map[string][]string),
		CenterTeam:         "",
		adjacent:           make(map[string]*tile),
	}
}

func (t *tile) copy() *tile {
	return newTile(t.Sides[SideTop], t.Sides[SideRight], t.Sides[SideBottom], t.Sides[SideLeft], t.Center, t.ConnectedCitySides, t.Banner)
}

func (t *tile) RotateRight() {
	newSides := make(map[string]string)
	for _, side := range Sides {
		newSides[ClockwiseSide[side]] = t.Sides[side]
	}
	t.Sides = newSides
}

func (t *tile) RotateLeft() {
	newSides := make(map[string]string)
	for _, side := range Sides {
		newSides[CounterClockwiseSide[side]] = t.Sides[side]
	}
	t.Sides = newSides
}

// given a city side on the tile get all connected city sides
func (t *tile) connectedCitySides(side string) ([]string, error) {
	sides := make([]string, 0)
	if !contains(Sides, side) {
		return nil, fmt.Errorf("invalid side %s", side)
	}
	if t.Sides[side] != City {
		return nil, fmt.Errorf("cannot enter tile on a non-city side")
	}
	// city sides are not connected so city ends here
	if !t.ConnectedCitySides {
		return sides, nil
	}
	for _, s := range Sides {
		if s != side && t.Sides[s] == City {
			sides = append(sides, s)
		}
	}
	return sides, nil
}

// given a road side on the tile get all connected road sides
func (t *tile) connectedRoadSides(side string) ([]string, error) {
	sides := make([]string, 0)
	if !contains(Sides, side) {
		return nil, fmt.Errorf("invalid side %s", side)
	}
	if t.Sides[side] != Road {
		return nil, fmt.Errorf("cannot enter tile on a non-road side")
	}
	for _, s := range Sides {
		if s != side && t.Sides[s] == Road {
			sides = append(sides, s)
		}
	}
	// more than one so road intersection
	if len(sides) != 1 {
		return make([]string, 0), nil
	}
	return sides, nil
}

func (t *tile) connectedFarmSides(farmSide string) ([]string, error) {
	points := make([]string, 0)
	if !contains(FarmSides, farmSide) {
		return nil, fmt.Errorf("invalid farm side %s", farmSide)
	}
	side := farmSideToSide(farmSide)
	ab := farmSideToAB(farmSide)
	if t.Sides[side] == City {
		return nil, fmt.Errorf("cannot enter on a city side")
	}

	// special case in which all sides are returned
	if t.Center == Cloister {
		points = append(points, sideToFarmSide(side, inverseAB(ab)),
			sideToFarmSide(ClockwiseSide[side], FarmNotchA), sideToFarmSide(ClockwiseSide[side], FarmNotchB),
			sideToFarmSide(CounterClockwiseSide[side], FarmNotchA), sideToFarmSide(CounterClockwiseSide[side], FarmNotchB),
			sideToFarmSide(AcrossSide[side], FarmNotchA), sideToFarmSide(AcrossSide[side], FarmNotchB))
		return points, nil
	}

	switch t.Sides[side] {
	case Farm:
		points = append(points, sideToFarmSide(side, inverseAB(ab)))

		clockwiseSide := ClockwiseSide[side]
		counterClockwiseSide := CounterClockwiseSide[side]
		acrossSide := AcrossSide[side]
		// clockwise check
		if t.Sides[clockwiseSide] == Road {
			points = append(points, sideToFarmSide(clockwiseSide, FarmNotchA))
		} else if t.Sides[clockwiseSide] == Farm {
			points = append(points, sideToFarmSide(clockwiseSide, FarmNotchA), sideToFarmSide(clockwiseSide, FarmNotchB))
		}
		// counterclockwise check
		if t.Sides[counterClockwiseSide] == Road {
			points = append(points, sideToFarmSide(counterClockwiseSide, FarmNotchB))
		} else if t.Sides[counterClockwiseSide] == Farm {
			points = append(points, sideToFarmSide(counterClockwiseSide, FarmNotchA), sideToFarmSide(counterClockwiseSide, FarmNotchB))
		}
		// if access to across side blocked return
		if (t.Sides[clockwiseSide] == Road && t.Sides[counterClockwiseSide] == Road) ||
			(t.Sides[clockwiseSide] == Road && t.Sides[counterClockwiseSide] == City) ||
			(t.Sides[clockwiseSide] == City && t.Sides[counterClockwiseSide] == Road) ||
			(t.Sides[clockwiseSide] == City && t.Sides[counterClockwiseSide] == City && t.ConnectedCitySides) {
			return points, nil
		}
		// across check
		if t.Sides[clockwiseSide] == Road && t.Sides[acrossSide] == Road {
			points = append(points, sideToFarmSide(acrossSide, FarmNotchB))
		} else if t.Sides[counterClockwiseSide] == Road && t.Sides[acrossSide] == Road {
			points = append(points, sideToFarmSide(acrossSide, FarmNotchA))
		} else if t.Sides[acrossSide] == Road || t.Sides[acrossSide] == Farm {
			points = append(points, sideToFarmSide(acrossSide, FarmNotchA), sideToFarmSide(acrossSide, FarmNotchB))
		}
	case Road:
		adjacentSide := ClockwiseSide[side]
		blockedAdjacentSide := CounterClockwiseSide[side]
		acrossSide := AcrossSide[side]
		if ab == FarmNotchA {
			adjacentSide = CounterClockwiseSide[side]
			blockedAdjacentSide = ClockwiseSide[side]
		}
		// adjacent side check
		if t.Sides[adjacentSide] == Road {
			points = append(points, sideToFarmSide(adjacentSide, inverseAB(ab)))
			return points, nil
		} else if t.Sides[adjacentSide] == Farm {
			points = append(points, sideToFarmSide(adjacentSide, FarmNotchA), sideToFarmSide(adjacentSide, FarmNotchB))
		}
		// across side check - return if blocked by road
		if t.Sides[acrossSide] == Road {
			points = append(points, sideToFarmSide(acrossSide, inverseAB(ab)))
			return points, nil
		} else if t.Sides[acrossSide] == Farm {
			points = append(points, sideToFarmSide(acrossSide, FarmNotchA), sideToFarmSide(acrossSide, FarmNotchB))
		}
		// blocked side check
		if t.Sides[blockedAdjacentSide] == Road {
			points = append(points, sideToFarmSide(blockedAdjacentSide, inverseAB(ab)))
		} else if t.Sides[blockedAdjacentSide] == Farm {
			points = append(points, sideToFarmSide(blockedAdjacentSide, FarmNotchA), sideToFarmSide(blockedAdjacentSide, FarmNotchB))
		}
	}
	return points, nil
}
