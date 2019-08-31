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
	s.states = newState3D(4, playfieldWidth, playfieldHeight)
	for y := 0; y < playfieldHeight; y++ {
		for x := 0; x < playfieldWidth; x++ {
			for rotation := 0; rotation < 4; rotation++ {
				s.states[y][x][rotation] = newState(x, y, rotation)
			}
		}
	}
}

func (s *searcher) lockTetrimino(playfield [][]int, tetriminoType, id int, stat *State) {
	squares := orientations[tetriminoType][stat.rotation].squares
	for i := 0; i < 4; i++ {
		square := squares[i]
		y := stat.y + square.y
		if y >= 0 {
			playfield[y][stat.x+square.x] = tetriminoType
			playfield[y][playfieldWidth]++
		}
	}
	s.searchListener.handleResult(playfield, tetriminoType, id, stat)
	for i := 0; i < 4; i++ {
		square := squares[i]
		y := stat.y + square.y
		if y >= 0 {
			playfield[y][stat.x+square.x] = tetriminoNone
			playfield[y][playfieldWidth]--
		}
	}
}

// returns true if the position is valid even if the node is not enqueued
func (s *searcher) addChild(playfield [][]int, tetriminoType, mark int, stat *State, x, y, rotation int) bool {

	orientation := orientations[tetriminoType][rotation]
	if x < orientation.minX || x > orientation.maxX || y > orientation.maxY {
		return false
	}

	childNode := s.states[y][x][rotation]
	if childNode.visited == mark {
		return true
	}

	squares := orientation.squares
	for i := 0; i < 4; i++ {
		square := squares[i]
		playfieldY := y + square.y
		if playfieldY >= 0 && playfield[playfieldY][x+square.x] != tetriminoNone {
			return false
		}
	}

	if s.positionValidator != nil && !s.positionValidator.validate(playfield, tetriminoType, x, y, rotation) {
		return true
	}

	childNode.visited = mark
	childNode.predecessor = stat

	s.q.enqueue(childNode)
	return true
}

func (s *searcher) search(playfield [][]int, tetriminoType, id int) bool {

	maxRotation := len(orientations[tetriminoType]) - 1

	globalMark++
	mark := globalMark

	if !s.addChild(playfield, tetriminoType, mark, nil, 5, 0, 0) {
		return false
	}

	for s.q.isNotEmpty() {
		stat := s.q.dequeue()

		if maxRotation != 0 {
			var r int
			if stat.rotation == 0 {
				r = maxRotation
			} else {
				r = stat.rotation - 1
			}
			s.addChild(playfield, tetriminoType, mark, stat, stat.x, stat.y, r)
			if maxRotation != 1 {
				if stat.rotation == maxRotation {
					r = 0
				} else {
					r = stat.rotation + 1
				}
				s.addChild(playfield, tetriminoType, mark, stat, stat.x, stat.y, r)
			}
		}

		s.addChild(playfield, tetriminoType, mark, stat, stat.x-1, stat.y, stat.rotation)
		s.addChild(playfield, tetriminoType, mark, stat, stat.x+1, stat.y, stat.rotation)

		if !s.addChild(playfield, tetriminoType, mark, stat, stat.x, stat.y+1, stat.rotation) {
			s.lockTetrimino(playfield, tetriminoType, id, stat)
		}
	}

	return true
}
