package go_carcassonne

/*
   Sides of a tile

            Top
             ——
       Left |  | Right
             ——
           Bottom
*/
const (
	SideTop    = "Top"
	SideRight  = "Right"
	SideBottom = "Bottom"
	SideLeft   = "Left"
)

var (
	Sides                = []string{SideTop, SideRight, SideBottom, SideLeft}
	ClockwiseSide        = map[string]string{SideTop: SideRight, SideRight: SideBottom, SideBottom: SideLeft, SideLeft: SideTop}
	CounterClockwiseSide = map[string]string{SideTop: SideLeft, SideLeft: SideBottom, SideBottom: SideRight, SideRight: SideTop}
	AcrossSide           = map[string]string{SideTop: SideBottom, SideRight: SideLeft, SideBottom: SideTop, SideLeft: SideRight}
)
