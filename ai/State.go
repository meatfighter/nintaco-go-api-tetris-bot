package ai

// State describes the position and orientation of a piece as its being dropped.
// It represents a node in the flood fill search algorithm.
type State struct {
	X           int
	Y           int
	Rotation    int
	Visited     int
	Predecessor *State
	Next        *State
}

func newState(x, y, rotation int) *State {
	return &State{
		X:        x,
		Y:        y,
		Rotation: rotation,
	}
}

func newState3D(dx, dy, dz int) [][][]*State {
	dxdy := dx * dy
	rows := make([]*State, dxdy*dz)
	s := make([][][]*State, dz)
	for z := dz - 1; z >= 0; z-- {
		offset := z * dxdy
		s[z] = make([][]*State, dy)
		for y := dy - 1; y >= 0; y-- {
			s[z][y] = rows[offset+dx*y : offset+dx*y+dx]
		}
	}
	return s
}
