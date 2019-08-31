package ai

type playfieldEvaluation struct {
	holes             int
	columnTransitions int
	rowTransitions    int
	wells             int
}

func newPlayfieldEvaluation() *playfieldEvaluation {
	return &playfieldEvaluation{}
}
