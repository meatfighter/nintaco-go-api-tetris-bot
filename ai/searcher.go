package ai

var globalMark = 1

type searcher struct {
	states            [][][]*State
	q                 *queue
	searchListener    iSearchListener
	positionValidator iChildFilter
}

func newSearcher(searchListener iSearchListener, positionValidator iChildFilter) *searcher {
	s := &searcher{
		q:                 newQueue(),
		searchListener:    searchListener,
		positionValidator: positionValidator,
	}
	s.createStates()
	return s
}

func (s *searcher) createStates() {
	s.states = newState3D(4, PlayfieldWidth, PlayfieldHeight)
	for y := 0; y < PlayfieldHeight; y++ {
		for x := 0; x < PlayfieldWidth; x++ {
			for rotation := 0; rotation < 4; rotation++ {
				s.states[y][x][rotation] = newState(x, y, rotation)
			}
		}
	}
}

func (s *searcher) lockTetrimino(playfield [][]int, tetriminoType, id int, stat *State) {
	squares := Orientations[tetriminoType][stat.Rotation].Squares
	for i := 0; i < 4; i++ {
		square := squares[i]
		y := stat.Y + square.y
		if y >= 0 {
			playfield[y][stat.X+square.x] = tetriminoType
			playfield[y][PlayfieldWidth]++
		}
	}
	s.searchListener.handleResult(playfield, tetriminoType, id, stat)
	for i := 0; i < 4; i++ {
		square := squares[i]
		y := stat.Y + square.y
		if y >= 0 {
			playfield[y][stat.X+square.x] = TetriminoNone
			playfield[y][PlayfieldWidth]--
		}
	}
}

// returns true if the position is valid even if the node is not enqueued
func (s *searcher) addChild(playfield [][]int, tetriminoType, mark int, stat *State, x, y, rotation int) bool {

	orientation := Orientations[tetriminoType][rotation]
	if x < orientation.MinX || x > orientation.MaxX || y > orientation.MaxY {
		return false
	}

	childNode := s.states[y][x][rotation]
	if childNode.Visited == mark {
		return true
	}

	squares := orientation.Squares
	for i := 0; i < 4; i++ {
		square := squares[i]
		playfieldY := y + square.y
		if playfieldY >= 0 && playfield[playfieldY][x+square.x] != TetriminoNone {
			return false
		}
	}

	if s.positionValidator != nil && !s.positionValidator.validate(playfield, tetriminoType, x, y, rotation) {
		return true
	}

	childNode.Visited = mark
	childNode.Predecessor = stat

	s.q.enqueue(childNode)
	return true
}

func (s *searcher) search(playfield [][]int, tetriminoType, id int) bool {

	maxRotation := len(Orientations[tetriminoType]) - 1

	globalMark++
	mark := globalMark

	if !s.addChild(playfield, tetriminoType, mark, nil, 5, 0, 0) {
		return false
	}

	for s.q.isNotEmpty() {
		stat := s.q.dequeue()

		if maxRotation != 0 {
			var r int
			if stat.Rotation == 0 {
				r = maxRotation
			} else {
				r = stat.Rotation - 1
			}
			s.addChild(playfield, tetriminoType, mark, stat, stat.X, stat.Y, r)
			if maxRotation != 1 {
				if stat.Rotation == maxRotation {
					r = 0
				} else {
					r = stat.Rotation + 1
				}
				s.addChild(playfield, tetriminoType, mark, stat, stat.X, stat.Y, r)
			}
		}

		s.addChild(playfield, tetriminoType, mark, stat, stat.X-1, stat.Y, stat.Rotation)
		s.addChild(playfield, tetriminoType, mark, stat, stat.X+1, stat.Y, stat.Rotation)

		if !s.addChild(playfield, tetriminoType, mark, stat, stat.X, stat.Y+1, stat.Rotation) {
			s.lockTetrimino(playfield, tetriminoType, id, stat)
		}
	}

	return true
}
