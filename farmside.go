package go_carcassonne

// Farm notch which is a piece of a farm side
const (
	FarmNotchA = "A"
	FarmNotchB = "B"
)

/*
   Farm sides of a tile

             TopA  TopB
               ——  ——
       LeftB |        | RightA
       LeftA |        | RightB
               ——  ——
          BottomB  BottomA
*/
const (
	FarmSideTopA    = SideTop + FarmNotchA
	FarmSideTopB    = SideTop + FarmNotchB
	FarmSideRightA  = SideRight + FarmNotchA
	FarmSideRightB  = SideRight + FarmNotchB
	FarmSideBottomA = SideBottom + FarmNotchA
	FarmSideBottomB = SideBottom + FarmNotchB
	FarmSideLeftA   = SideLeft + FarmNotchA
	FarmSideLeftB   = SideLeft + FarmNotchB
)

var (
	FarmSides      = []string{FarmSideTopA, FarmSideTopB, FarmSideRightA, FarmSideRightB, FarmSideBottomA, FarmSideBottomB, FarmSideLeftA, FarmSideLeftB}
	AcrossFarmSide = map[string]string{FarmSideTopA: FarmSideBottomB, FarmSideTopB: FarmSideBottomA, FarmSideRightA: FarmSideLeftB, FarmSideRightB: FarmSideLeftA, FarmSideBottomA: FarmSideTopB, FarmSideBottomB: FarmSideTopA, FarmSideLeftA: FarmSideRightB, FarmSideLeftB: FarmSideRightA}
)

func farmSideToSide(farmSide string) string {
	return farmSide[:len(farmSide)-1]
}

func farmSideToAB(farmSide string) string {
	return string(farmSide[len(farmSide)-1])
}

func sideToFarmSide(side, ab string) string {
	return side + ab
}

func inverseAB(ab string) string {
	if ab == FarmNotchA {
		return FarmNotchB
	}
	return FarmNotchA
}
