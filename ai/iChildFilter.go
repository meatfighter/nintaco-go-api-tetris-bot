package ai

type iChildFilter interface {
	validate(playfield [][]int, tetriminoType, x, y, rotation int) bool
}
