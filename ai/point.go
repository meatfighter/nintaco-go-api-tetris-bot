package ai

type point struct {
	x int
	y int
}

func newOriginPoint() *point {
	return &point{}
}

func newPoint(x, y int) *point {
	return &point{
		x: x,
		y: y,
	}
}
