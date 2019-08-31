package ai

// Orientation ...
type Orientation struct {
	Squares       []*point
	MinX          int
	MaxX          int
	MaxY          int
	OrientationID int
}

func newOrientation() *Orientation {
	o := &Orientation{
		Squares: make([]*point, 4),
	}
	for i := 0; i < 4; i++ {
		o.Squares[i] = newOriginPoint()
	}
	return o
}
