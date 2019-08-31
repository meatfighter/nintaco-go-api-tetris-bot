package ai

import "math"

func max(x, y int) int {
	if x < y {
		return y
	}
	return x
}

func min(x, y int) int {
	if x > y {
		return y
	}
	return x
}

// The pieces.
const (
	TetriminoNone = -1
	TetriminoT    = 0
	TetriminoJ    = 1
	TetriminoZ    = 2
	TetriminoO    = 3
	TetriminoS    = 4
	TetriminoL    = 5
	TetriminoI    = 6
)

var patterns = [][][][]int{
	{{{-1, 0}, {0, 0}, {1, 0}, {0, 1}}, // Td (spawn)
		{{0, -1}, {-1, 0}, {0, 0}, {0, 1}}, // Tl
		{{-1, 0}, {0, 0}, {1, 0}, {0, -1}}, // Tu
		{{0, -1}, {0, 0}, {1, 0}, {0, 1}}}, // Tr

	{{{-1, 0}, {0, 0}, {1, 0}, {1, 1}}, // Jd (spawn)
		{{0, -1}, {0, 0}, {-1, 1}, {0, 1}},  // Jl
		{{-1, -1}, {-1, 0}, {0, 0}, {1, 0}}, // Ju
		{{0, -1}, {1, -1}, {0, 0}, {0, 1}}}, // Jr

	{{{-1, 0}, {0, 0}, {0, 1}, {1, 1}}, // Zh (spawn)
		{{1, -1}, {0, 0}, {1, 0}, {0, 1}}}, // Zv

	{{{-1, 0}, {0, 0}, {-1, 1}, {0, 1}}}, // O  (spawn)

	{{{0, 0}, {1, 0}, {-1, 1}, {0, 1}}, // Sh (spawn)
		{{0, -1}, {0, 0}, {1, 0}, {1, 1}}}, // Sv

	{{{-1, 0}, {0, 0}, {1, 0}, {-1, 1}}, // Ld (spawn)
		{{-1, -1}, {0, -1}, {0, 0}, {0, 1}}, // Ll
		{{1, -1}, {-1, 0}, {0, 0}, {1, 0}},  // Lu
		{{0, -1}, {0, 0}, {0, 1}, {1, 1}}},  // Lr

	{{{-2, 0}, {-1, 0}, {0, 0}, {1, 0}}, // Ih (spawn)
		{{0, -2}, {0, -1}, {0, 0}, {0, 1}}}, // Iv
}

var orientationIDs = []int{
	0x02, 0x03, 0x00, 0x01, 0x07, 0x04, 0x05, 0x06, 0x08, 0x09,
	0x0A, 0x0B, 0x0C, 0x0E, 0x0F, 0x10, 0x0D, 0x12, 0x11}

// Orientations is a table of all the possible orientations of the pieces.
var Orientations = func() [][]*Orientation {
	o := make([][]*Orientation, len(patterns))
	for i, idIndex := 0, 0; i < len(patterns); i++ {
		tetriminos := []*Orientation{}
		for j := 0; j < len(patterns[i]); j++ {
			tetrimino := newOrientation()
			tetriminos = append(tetriminos, tetrimino)
			minX := math.MaxInt32
			maxX := math.MinInt32
			maxY := math.MinInt32
			for k := 0; k < 4; k++ {
				p := patterns[i][j][k]
				tetrimino.Squares[k].x = p[0]
				tetrimino.Squares[k].y = p[1]
				minX = min(minX, p[0])
				maxX = max(maxX, p[0])
				maxY = max(maxY, p[1])
			}
			tetrimino.MinX = -minX
			tetrimino.MaxX = PlayfieldWidth - maxX - 1
			tetrimino.MaxY = PlayfieldHeight - maxY - 1
			tetrimino.OrientationID = orientationIDs[idIndex]
			idIndex++
		}
		o[i] = tetriminos
	}
	return o
}()
