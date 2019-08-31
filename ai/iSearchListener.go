package ai

type iSearchListener interface {
	handleResult(playfield [][]int, tetriminoType, id int, s *State)
}
