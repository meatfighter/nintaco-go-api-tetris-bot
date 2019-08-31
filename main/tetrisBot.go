package main

import (
	"os"
	"strings"

	"github.com/meatfighter/nintaco-go-api-tetris-bot/ai"
	"github.com/meatfighter/nintaco-go-api/nintaco"
)

const emptySquare = 0xEF

type tetrisBot struct {
	api            nintaco.API
	ai             *ai.AI
	playfieldUtil  *ai.PlayfieldUtil
	tetriminos     []int
	playfield      [][]int
	tetriminoTypes []int
	playFast       bool

	playingDelay     int
	targetTetriminoY int
	startCounter     int
	movesIndex       int
	moving           bool
	states           []*ai.State
}

func newTetrisBot(playFast bool) *tetrisBot {
	t := &tetrisBot{
		api:            nintaco.GetAPI(),
		ai:             ai.NewAI(),
		playfieldUtil:  ai.NewPlayfieldUtil(),
		tetriminos:     make([]int, ai.TetriminosSearched),
		tetriminoTypes: make([]int, 19),
		playFast:       playFast,
	}
	t.playfield = t.playfieldUtil.CreatePlayfield()
	return t
}

func (t *tetrisBot) launch() {
	/*api.addActivateListener(this::apiEnabled);
	  api.addAccessPointListener(this::updateScore, AccessPointType.PreExecute,
	      0x9C35);
	  api.addAccessPointListener(this::speedUpDrop, AccessPointType.PreExecute,
	      0x8977);
	  api.addAccessPointListener(this::tetriminoYUpdated,
	      AccessPointType.PreWrite, Addresses.TetriminoY1);
	  api.addAccessPointListener(this::tetriminoYUpdated,
	      AccessPointType.PreWrite, Addresses.TetriminoY2);
	  api.addFrameListener(this::renderFinished);
	  api.addStatusListener(this::statusChanged);*/
	t.api.Run()
}

func (t *tetrisBot) APIEnabled() {
	//t.readTetriminoTypes()
}

/*func (t *tetrisBot) readTetriminoTypes() {
	for i := 0; i < 19; i++ {
		TetriminosTypes[i] = api.readCPU(Addresses.TetriminoTypeTable + i)
	}
}*/

func main() {
	nintaco.InitRemoteAPI("localhost", 9999)
	newTetrisBot(len(os.Args) > 1 && strings.EqualFold("fast", os.Args[1])).launch()
}
