package ai

// Orientation describes a rotated piece. It includes the coordinates of the squares
// relative to a central pivot and the range that the piece can move within the playfield.
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
