package ai

type orientation struct {
	squares       []*point
	minX          int
	maxX          int
	maxY          int
	orientationID int
}

func newOrientation() *orientation {
	o := &orientation{
		squares: make([]*point, 4),
	}
	for i := 0; i < 4; i++ {
		o.squares[i] = newOriginPoint()
	}
	return o
}
