package ai

import "math"

// TODO COMMENT
const (
	PlayfieldWidth  = 10
	PlayfieldHeight = 20

	TetriminosSearched = 2
)

var weights = []float64{
	1.0,
	12.885008263218383,
	15.842707182438396,
	26.89449650779595,
	27.616914062397015,
	30.18511071927904,
}

// AI ...
type AI struct {
	searchers        []*searcher
	tetriminoIndices []int
	playfieldU       *PlayfieldUtil
	e                *playfieldEvaluation
	totalRows        int
	totalDropHeight  int
	bestFitness      float64
	bestResult       *State
	result0          *State
	searchListener   iSearchListener
}

// NewAI ...
func NewAI() *AI {
	return NewAI2(nil)
}

// NewAI2 ...
func NewAI2(positionValidator iChildFilter) *AI {
	a := &AI{
		playfieldU: NewPlayfieldUtil(),
		e:          newPlayfieldEvaluation(),
	}
	a.searchListener = a
	a.searchers = make([]*searcher, TetriminosSearched)
	for i := 0; i < TetriminosSearched; i++ {
		a.searchers[i] = newSearcher(a.searchListener, positionValidator)
	}
	return a
}

func (a *AI) handleResult(playfield [][]int, tetriminoType, id int, s *State) {
	if id == 0 {
		a.result0 = s
	}

	orientation := orientations[tetriminoType][s.rotation]
	rows := a.playfieldU.ClearRows(playfield, s.y)
	originalTotalRows := a.totalRows
	originalTotalDropHeight := a.totalDropHeight
	a.totalRows += rows
	a.totalDropHeight += orientation.maxY - s.y

	nextID := id + 1

	if nextID == len(a.tetriminoIndices) {

		a.playfieldU.EvaluatePlayfield(playfield, a.e)

		fitness := a.computeFitness()
		if fitness < a.bestFitness {
			a.bestFitness = fitness
			a.bestResult = a.result0
		}
	} else {
		a.searchers[nextID].search(playfield, a.tetriminoIndices[nextID], nextID)
	}

	a.totalDropHeight = originalTotalDropHeight
	a.totalRows = originalTotalRows
	a.playfieldU.RestoreRows(playfield, rows)
}

func (a *AI) computeFitness() float64 {
	return weights[0]*float64(a.totalRows) +
		weights[1]*float64(a.totalDropHeight) +
		weights[2]*float64(a.e.wells) +
		weights[3]*float64(a.e.holes) +
		weights[4]*float64(a.e.columnTransitions) +
		weights[5]*float64(a.e.rowTransitions)
}

// Search ...
func (a *AI) Search(playfield [][]int, tetriminoIndices []int) *State {

	a.tetriminoIndices = tetriminoIndices
	a.bestResult = nil
	a.bestFitness = math.MaxFloat64

	a.searchers[0].search(playfield, tetriminoIndices[0], 0)

	return a.bestResult
}

// BuildStatesList ...
func (a *AI) BuildStatesList(state *State) []*State {
	s := state
	count := 0
	for s != nil {
		count++
		s = s.predecessor
	}
	states := make([]*State, count)
	for state != nil {
		count--
		states[count] = state
		state = state.predecessor
	}
	return states
}
