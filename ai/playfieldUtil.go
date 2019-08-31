package ai

func newInt2D(width, height int) [][]int {
	a := make([][]int, height)
	rows := make([]int, width*height)
	for y := height - 1; y >= 0; y-- {
		a[y] = rows[width*y : width*(y+1)]
	}
	return a
}

// PlayfieldUtil ...
type PlayfieldUtil struct {
	spareRows    [][]int
	columnDepths []int
	spareIndex   int
}

// NewPlayfieldUtil ...
func NewPlayfieldUtil() *PlayfieldUtil {
	p := &PlayfieldUtil{
		spareRows:    newInt2D(PlayfieldWidth+1, 8*TetriminosSearched),
		columnDepths: make([]int, PlayfieldWidth),
	}
	for y := 0; y < len(p.spareRows); y++ {
		for x := 0; x < PlayfieldWidth; x++ {
			p.spareRows[y][x] = TetriminoNone
		}
	}
	return p
}

// CreatePlayfield ...
func (p *PlayfieldUtil) CreatePlayfield() [][]int {
	playfield := newInt2D(PlayfieldWidth+1, PlayfieldHeight)
	for y := 0; y < PlayfieldHeight; y++ {
		for x := 0; x < PlayfieldWidth; x++ {
			playfield[y][x] = TetriminoNone
		}
	}
	return playfield
}

// LockTetrimino ...
func (p *PlayfieldUtil) LockTetrimino(playfield [][]int, tetriminoType int, s *State) {

	squares := Orientations[tetriminoType][s.Rotation].Squares
	for i := 0; i < 4; i++ {
		square := squares[i]
		y := s.Y + square.y
		if y >= 0 {
			playfield[y][s.X+square.x] = tetriminoType
			playfield[y][PlayfieldWidth]++
		}
	}

	startRow := s.Y - 2
	endRow := s.Y + 1

	if startRow < 1 {
		startRow = 1
	}
	if endRow >= PlayfieldHeight {
		endRow = PlayfieldHeight - 1
	}

	for y := startRow; y <= endRow; y++ {
		if playfield[y][PlayfieldWidth] == PlayfieldWidth {
			clearedRow := playfield[y]
			for i := y; i > 0; i-- {
				playfield[i] = playfield[i-1]
			}
			for x := 0; x < PlayfieldWidth; x++ {
				clearedRow[x] = TetriminoNone
			}
			clearedRow[PlayfieldWidth] = 0
			playfield[0] = clearedRow
		}
	}
}

// EvaluatePlayfield ...
func (p *PlayfieldUtil) EvaluatePlayfield(playfield [][]int, e *playfieldEvaluation) {

	for x := 0; x < PlayfieldWidth; x++ {
		p.columnDepths[x] = PlayfieldHeight - 1
		for y := 0; y < PlayfieldHeight; y++ {
			if playfield[y][x] != TetriminoNone {
				p.columnDepths[x] = y
				break
			}
		}
	}

	e.wells = 0
	for x := 0; x < PlayfieldWidth; x++ {
		var minY int
		if x == 0 {
			minY = p.columnDepths[1]
		} else if x == PlayfieldWidth-1 {
			minY = p.columnDepths[PlayfieldWidth-2]
		} else {
			minY = max(p.columnDepths[x-1], p.columnDepths[x+1])
		}
		for y := p.columnDepths[x]; y >= minY; y-- {
			if (x == 0 || playfield[y][x-1] != TetriminoNone) &&
				(x == PlayfieldWidth-1 || playfield[y][x+1] != TetriminoNone) {
				e.wells++
			}
		}
	}

	e.holes = 0
	e.columnTransitions = 0
	for x := 0; x < PlayfieldWidth; x++ {
		solid := true
		for y := p.columnDepths[x] + 1; y < PlayfieldHeight; y++ {
			if playfield[y][x] == TetriminoNone {
				if playfield[y-1][x] != TetriminoNone {
					e.holes++
				}
				if solid {
					solid = false
					e.columnTransitions++
				}
			} else if !solid {
				solid = true
				e.columnTransitions++
			}
		}
	}

	e.rowTransitions = 0
	for y := 0; y < PlayfieldHeight; y++ {
		solidFound := false
		solid := true
		transitions := 0
		for x := 0; x <= PlayfieldWidth; x++ {
			if x == PlayfieldWidth {
				if !solid {
					transitions++
				}
			} else {
				if playfield[y][x] == TetriminoNone {
					if solid {
						solid = false
						transitions++
					}
				} else {
					solidFound = true
					if !solid {
						solid = true
						transitions++
					}
				}
			}
		}
		if solidFound {
			e.rowTransitions += transitions
		}
	}
}

// ClearRows ...
func (p *PlayfieldUtil) ClearRows(playfield [][]int, tetriminoY int) int {

	rows := 0
	startRow := tetriminoY - 2
	endRow := tetriminoY + 1

	if startRow < 1 {
		startRow = 1
	}
	if endRow >= PlayfieldHeight {
		endRow = PlayfieldHeight - 1
	}

	for y := startRow; y <= endRow; y++ {
		if playfield[y][PlayfieldWidth] == PlayfieldWidth {
			rows++
			p.clearRow(playfield, y)
		}
	}

	return rows
}

func (p *PlayfieldUtil) clearRow(playfield [][]int, y int) {
	clearedRow := playfield[y]
	clearedRow[PlayfieldWidth] = y
	for i := y; i > 0; i-- {
		playfield[i] = playfield[i-1]
	}
	playfield[0] = p.spareRows[p.spareIndex]
	playfield[0][PlayfieldWidth] = 0

	p.spareRows[p.spareIndex] = clearedRow
	p.spareIndex++
}

func (p *PlayfieldUtil) restoreRow(playfield [][]int) {
	p.spareIndex--
	restoredRow := p.spareRows[p.spareIndex]
	y := restoredRow[PlayfieldWidth]

	p.spareRows[p.spareIndex] = playfield[0]

	for i := 0; i < y; i++ {
		playfield[i] = playfield[i+1]
	}
	restoredRow[PlayfieldWidth] = PlayfieldWidth
	playfield[y] = restoredRow
}

// RestoreRows ...
func (p *PlayfieldUtil) RestoreRows(playfield [][]int, rows int) {
	for i := 0; i < rows; i++ {
		p.restoreRow(playfield)
	}
}
