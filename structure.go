package go_carcassonne

import "fmt"

// Structure types
const (
	Farm         = "Farm"
	City         = "City"
	Road         = "Road"
	Cloister     = "Cloister"
	NilStructure = "NilStructure"
)

var StructureTypeToTokenType = map[string]string{Farm: Farmer, City: Knight, Road: Thief, Cloister: Monk}

// node that is part of a complete/incomplete City, Road, or Farm structure
type node struct {
	tile  *tile
	sides []string
}

// full set of nodes that make up a complete structure to create a complete/incomplete City, Road, or Farm
type structure struct {
	typ      string
	complete bool
	nodes    []*node
}

// get the score of a city structure by checking all city sections
func scoreCity(city *structure) (int, error) {
	if city.typ != City {
		return 0, fmt.Errorf("structure is not a city")
	}
	points := 0
	increment := 1
	if city.complete {
		increment = 2
	}
	for _, n := range city.nodes {
		points += increment
		if n.tile.Banner {
			points += increment
		}
	}
	return points, nil
}

// get the score of a road structure by checking all road sections
func scoreRoad(road *structure) (int, error) {
	if road.typ != Road {
		return 0, fmt.Errorf("structure is not a road")
	}
	return len(road.nodes), nil
}

// get the score of a farm structure by checking number of complete cities that are connected to the farm
func scoreFarm(farm *structure, completeCities []*structure) (int, error) {
	if farm.typ != Farm {
		return 0, fmt.Errorf("structure is not a farm")
	}
	points := 0
	for _, city := range completeCities {
		if !city.complete {
			continue
		}
	cityTouching:
		for _, fNode := range farm.nodes {
			for _, cNode := range city.nodes {
				if fNode.tile.X == cNode.tile.X && fNode.tile.Y == cNode.tile.Y {
					// find all farm sides touching this city on the tile
					touchingFarmSides := make([]string, 0)
					for _, cSide := range cNode.sides {
						clockwiseSection := cNode.tile.Sides[ClockwiseSide[cSide]]
						counterClockwiseSection := cNode.tile.Sides[CounterClockwiseSide[cSide]]
						acrossSection := cNode.tile.Sides[AcrossSide[cSide]]
						// clockwise side check
						if clockwiseSection == Road {
							touchingFarmSides = append(touchingFarmSides, sideToFarmSide(ClockwiseSide[cSide], FarmNotchA))
						} else if clockwiseSection == Farm {
							touchingFarmSides = append(touchingFarmSides, sideToFarmSide(ClockwiseSide[cSide], FarmNotchA), sideToFarmSide(ClockwiseSide[cSide], FarmNotchB))
						}
						// counterclockwise side check
						if counterClockwiseSection == Road {
							touchingFarmSides = append(touchingFarmSides, sideToFarmSide(CounterClockwiseSide[cSide], FarmNotchB))
						} else if counterClockwiseSection == Farm {
							touchingFarmSides = append(touchingFarmSides, sideToFarmSide(CounterClockwiseSide[cSide], FarmNotchA), sideToFarmSide(CounterClockwiseSide[cSide], FarmNotchB))
						}
						// cannot get across so continue
						if (clockwiseSection == Road && counterClockwiseSection == Road) ||
							(clockwiseSection == City && counterClockwiseSection == Road) ||
							(clockwiseSection == Road && counterClockwiseSection == City) {
							continue
						}
						// across side check
						if acrossSection == Farm {
							touchingFarmSides = append(touchingFarmSides, sideToFarmSide(AcrossSide[cSide], FarmNotchA), sideToFarmSide(AcrossSide[cSide], FarmNotchB))
						} else if acrossSection == Road && clockwiseSection == Road {
							touchingFarmSides = append(touchingFarmSides, sideToFarmSide(AcrossSide[cSide], FarmNotchB))
						} else if acrossSection == Road && counterClockwiseSection == Road {
							touchingFarmSides = append(touchingFarmSides, sideToFarmSide(AcrossSide[cSide], FarmNotchA))
						}
					}
					// check if farm side in farm node is touching city
					for _, touchingFarmSide := range touchingFarmSides {
						for _, farmSide := range fNode.sides {
							if farmSide == touchingFarmSide {
								points += 3
								break cityTouching
							}
						}
					}
				}
			}
		}
	}
	return points, nil
}
