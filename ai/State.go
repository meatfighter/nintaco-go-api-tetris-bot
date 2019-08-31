package ai

// State ...
type State struct {
	x           int
	y           int
	rotation    int
	visited     int
	predecessor *State
	next        *State
}

func newState(x, y, rotation int) *State {
	return &State{
		x:        x,
		y:        y,
		rotation: rotation,
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
