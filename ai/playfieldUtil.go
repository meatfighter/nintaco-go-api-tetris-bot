package ai

func newInt2D(width, height int) [][]int {
	a := make([][]int, height)
	rows := make([]int, width*height)
	for y := height - 1; y >= 0; y-- {
		a[y] = rows[width*y : width*(y+1)]
	}
	return a
}

type playfieldUtil struct {
	spareRows    [][]int
	columnDepths []int
	spareIndex   int
}

func newPlayfieldUtil() *playfieldUtil {
	p := &playfieldUtil{
		spareRows:    newInt2D(playfieldWidth+1, 8*tetriminosSearched),
		columnDepths: make([]int, playfieldWidth),
	}
	for y := 0; y < len(p.spareRows); y++ {
		for x := 0; x < playfieldWidth; x++ {
			p.spareRows[y][x] = tetriminoNone
		}
	}
	return p
}

func (p *playfieldUtil) createPlayfield() [][]int {
	playfield := newInt2D(playfieldWidth+1, playfieldHeight)
	for y := 0; y < playfieldHeight; y++ {
		for x := 0; x < playfieldWidth; x++ {
			playfield[y][x] = tetriminoNone
		}
	}
	return playfield
}

func (p *playfieldUtil) lockTetrimino(playfield [][]int, tetriminoType int, s *State) {

	squares := orientations[tetriminoType][s.rotation].squares
	for i := 0; i < 4; i++ {
		square := squares[i]
		y := s.y + square.y
		if y >= 0 {
			playfield[y][s.x+square.x] = tetriminoType
			playfield[y][playfieldWidth]++
		}
	}

	startRow := s.y - 2
	endRow := s.y + 1

	if startRow < 1 {
		startRow = 1
	}
	if endRow >= playfieldHeight {
		endRow = playfieldHeight - 1
	}

	for y := startRow; y <= endRow; y++ {
		if playfield[y][playfieldWidth] == playfieldWidth {
			clearedRow := playfield[y]
			for i := y; i > 0; i-- {
				playfield[i] = playfield[i-1]
			}
			for x := 0; x < playfieldWidth; x++ {
				clearedRow[x] = tetriminoNone
			}
			clearedRow[playfieldWidth] = 0
			playfield[0] = clearedRow
		}
	}
}

func (p *playfieldUtil) evaluatePlayfield(playfield [][]int, e *playfieldEvaluation) {

	for x := 0; x < playfieldWidth; x++ {
		p.columnDepths[x] = playfieldHeight - 1
		for y := 0; y < playfieldHeight; y++ {
			if playfield[y][x] != tetriminoNone {
				p.columnDepths[x] = y
				break
			}
		}
	}

	e.wells = 0
	for x := 0; x < playfieldWidth; x++ {
		var minY int
		if x == 0 {
			minY = p.columnDepths[1]
		} else if x == playfieldWidth-1 {
			minY = p.columnDepths[playfieldWidth-2]
		} else {
			minY = max(p.columnDepths[x-1], p.columnDepths[x+1])
		}
		for y := p.columnDepths[x]; y >= minY; y-- {
			if (x == 0 || playfield[y][x-1] != tetriminoNone) &&
				(x == playfieldWidth-1 || playfield[y][x+1] != tetriminoNone) {
				e.wells++
			}
		}
	}

	e.holes = 0
	e.columnTransitions = 0
	for x := 0; x < playfieldWidth; x++ {
		solid := true
		for y := p.columnDepths[x] + 1; y < playfieldHeight; y++ {
			if playfield[y][x] == tetriminoNone {
				if playfield[y-1][x] != tetriminoNone {
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
	for y := 0; y < playfieldHeight; y++ {
		solidFound := false
		solid := true
		transitions := 0
		for x := 0; x <= playfieldWidth; x++ {
			if x == playfieldWidth {
				if !solid {
					transitions++
				}
			} else {
				if playfield[y][x] == tetriminoNone {
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

func (p *playfieldUtil) clearRows(playfield [][]int, tetriminoY int) int {

	rows := 0
	startRow := tetriminoY - 2
	endRow := tetriminoY + 1

	if startRow < 1 {
		startRow = 1
	}
	if endRow >= playfieldHeight {
		endRow = playfieldHeight
	}

	for y := startRow; y <= endRow; y++ {
		if playfield[y][playfieldWidth] == playfieldWidth {
			rows++
			p.clearRow(playfield, y)
		}
	}

	return rows
}

func (p *playfieldUtil) clearRow(playfield [][]int, y int) {
	clearedRow := playfield[y]
	clearedRow[playfieldWidth] = y
	for i := y; i > 0; i-- {
		playfield[i] = playfield[i-1]
	}
	playfield[0] = p.spareRows[p.spareIndex]
	playfield[0][playfieldWidth] = 0

	p.spareRows[p.spareIndex] = clearedRow
	p.spareIndex++
}

func (p *playfieldUtil) restoreRow(playfield [][]int) {
	p.spareIndex--
	restoredRow := p.spareRows[p.spareIndex]
	y := restoredRow[playfieldWidth]

	p.spareRows[p.spareIndex] = playfield[0]

	for i := 0; i < y; i++ {
		playfield[i] = playfield[i+1]
	}
	restoredRow[playfieldWidth] = playfieldWidth
	playfield[y] = restoredRow
}

func (p *playfieldUtil) restoreRows(playfield [][]int, rows int) {
	for i := 0; i < rows; i++ {
		p.restoreRow(playfield)
	}
}
