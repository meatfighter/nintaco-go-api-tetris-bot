package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/meatfighter/nintaco-go-api-tetris-bot/ai"
	"github.com/meatfighter/nintaco-go-api/nintaco"
)

const emptySquare = 0xEF

type tetrisBot struct {
	api             nintaco.API
	ai              *ai.AI
	playfieldUtil   *ai.PlayfieldUtil
	tetriminos      []int
	playfield       [][]int
	tetriminosTypes []int
	playFast        bool

	playingDelay     int
	targetTetriminoY int
	startCounter     int
	movesIndex       int
	moving           bool
	states           []*ai.State
}

func newTetrisBot(playFast bool) *tetrisBot {
	t := &tetrisBot{
		api:             nintaco.GetAPI(),
		ai:              ai.NewAI(),
		playfieldUtil:   ai.NewPlayfieldUtil(),
		tetriminos:      make([]int, ai.TetriminosSearched),
		tetriminosTypes: make([]int, 19),
		playFast:        playFast,
	}
	t.playfield = t.playfieldUtil.CreatePlayfield()
	return t
}

func (t *tetrisBot) launch() {
	t.api.AddActivateListener(nintaco.ActivateFunc(t.apiEnabled))
	t.api.AddAccessPointListener(nintaco.AccessPointFunc(t.updateScore),
		nintaco.AccessPointTypePreExecute, 0x9C35)
	t.api.AddAccessPointListener(nintaco.AccessPointFunc(t.speedUpDrop),
		nintaco.AccessPointTypePreExecute, 0x8977)
	t.api.AddAccessPointListener(nintaco.AccessPointFunc(t.tetriminoYUpdated),
		nintaco.AccessPointTypePreWrite, addressTetriminoY1)
	t.api.AddAccessPointListener(nintaco.AccessPointFunc(t.tetriminoYUpdated),
		nintaco.AccessPointTypePreWrite, addressTetriminoY2)
	t.api.AddFrameListener(nintaco.FrameFunc(t.renderFinished))
	t.api.AddStatusListener(nintaco.StatusFunc(t.statusChanged))
	t.api.Run()
}

func (t *tetrisBot) apiEnabled() {
	t.readTetriminoTypes()
}

func (t *tetrisBot) tetriminoYUpdated(typ, address, tetriminoY int) int {
	if tetriminoY == 0 {
		t.targetTetriminoY = 0
	}
	if t.moving {
		return t.targetTetriminoY
	}
	return tetriminoY
}

func (t *tetrisBot) readTetriminoTypes() {
	for i := 0; i < 19; i++ {
		t.tetriminosTypes[i] = t.api.ReadCPU(addressTetriminoTypeTable + i)
	}
}

func (t *tetrisBot) resetPlayState(gameState int) {
	if gameState != 4 {
		t.api.WriteCPU(addressPlayState, 0)
	}
}

func (t *tetrisBot) updateScore(typ, address, value int) int {
	// cap the points multiplier at 30 to avoid the kill screen
	if t.api.ReadCPU(0x00A8) > 30 {
		t.api.WriteCPU(0x00A8, 30)
	}
	return -1
}

func (t *tetrisBot) speedUpDrop(typ, address, value int) int {
	t.api.SetX(0x1E)
	return -1
}

func (t *tetrisBot) setTetriminoYAddress(address, y int) {
	t.targetTetriminoY = y
	t.api.WriteCPU(address, y)
}

func (t *tetrisBot) setTetriminoY(y int) {
	t.setTetriminoYAddress(addressTetriminoY1, y)
	t.setTetriminoYAddress(addressTetriminoY2, y)
}

func (t *tetrisBot) makeMove(tetriminoType int, state *ai.State, finalMove bool) {
	if finalMove {
		t.api.WriteCPU(0x006E, 0x03)
	}
	t.api.WriteCPU(addressTetriminoX, state.X)
	t.setTetriminoY(state.Y)
	t.api.WriteCPU(addressTetriminoID, ai.Orientations[tetriminoType][state.Rotation].OrientationID)
}

func (t *tetrisBot) readTetrimino() int {
	return t.tetriminosTypes[t.api.ReadCPU(addressTetriminoID)]
}

func (t *tetrisBot) readNextTetrimino() int {
	return t.tetriminosTypes[t.api.ReadCPU(addressNextTetriminoID)]
}

func (t *tetrisBot) readPlayfield() {
	t.tetriminos[0] = t.readTetrimino()
	t.tetriminos[1] = t.readNextTetrimino()

	for i := 0; i < ai.PlayfieldHeight; i++ {
		t.playfield[i][10] = 0
		for j := 0; j < ai.PlayfieldWidth; j++ {
			if t.api.ReadCPU(addressPlayfield+10*i+j) == emptySquare {
				t.playfield[i][j] = ai.TetriminoNone
			} else {
				t.playfield[i][j] = ai.TetriminoI
				t.playfield[i][10]++
			}
		}
	}
}

func (t *tetrisBot) spawned() bool {
	currentTetrimino := t.api.ReadCPU(addressTetriminoID)
	playState := t.api.ReadCPU(addressPlayState)
	tetriminoX := t.api.ReadCPU(addressTetriminoX)
	tetriminoY := t.api.ReadCPU(addressTetriminoY1)

	return playState == 1 && tetriminoX == 5 && tetriminoY == 0 && currentTetrimino < len(t.tetriminosTypes)
}

func (t *tetrisBot) isPlaying(gameState int) bool {
	return gameState == 4 && t.api.ReadCPU(addressPlayState) < 9
}

func (t *tetrisBot) pressStart() {
	if t.startCounter > 0 {
		t.startCounter--
	} else {
		t.startCounter = 10
	}
	if t.startCounter >= 5 {
		t.api.WriteGamepad(0, nintaco.GamepadButtonStart, true)
	}
}

func (t *tetrisBot) skipCopyrightScreen(gameState int) {
	if gameState == 0 {
		if t.api.ReadCPU(addressCopyright1) > 1 {
			t.api.WriteCPU(addressCopyright1, 0)
		} else if t.api.ReadCPU(addressCopyright2) > 2 {
			t.api.WriteCPU(addressCopyright2, 1)
		}
	}
}

func (t *tetrisBot) skipTitleAndDemoScreens(gameState int) {
	if gameState == 1 || gameState == 5 {
		t.pressStart()
	} else {
		t.startCounter = 0
	}
}

func (t *tetrisBot) renderFinished() {
	gameState := t.api.ReadCPU(addressGameState)
	t.skipCopyrightScreen(gameState)
	t.skipTitleAndDemoScreens(gameState)
	t.resetPlayState(gameState)

	if t.isPlaying(gameState) {
		if t.playingDelay > 0 {
			t.playingDelay--
		} else if t.playFast {
			// skip line clearing animation
			if t.api.ReadCPU(addressPlayState) == 4 {
				t.api.WriteCPU(addressPlayState, 5)
			}
			if t.spawned() {
				t.readPlayfield()
				state := t.ai.Search(t.playfield, t.tetriminos)
				if state != nil {
					t.moving = true
					t.makeMove(t.tetriminos[0], state, true)
					t.moving = false
				}
			}
		} else {
			if t.moving && t.movesIndex < len(t.states) {
				t.makeMove(t.tetriminos[0], t.states[t.movesIndex], t.movesIndex == len(t.states)-1)
				t.movesIndex++
			} else {
				t.moving = false
				if t.spawned() {
					t.readPlayfield()
					state := t.ai.Search(t.playfield, t.tetriminos)
					if state != nil {
						t.states = t.ai.BuildStatesList(state)
						t.movesIndex = 0
						t.moving = true
					}
				}
			}
		}
	} else {
		t.states = nil
		t.moving = false
		t.playingDelay = 16
	}
}

func (t *tetrisBot) statusChanged(message string) {
	fmt.Println(message)
}

func main() {
	nintaco.InitRemoteAPI("localhost", 9999)
	newTetrisBot(len(os.Args) > 1 && strings.EqualFold("fast", os.Args[1])).launch()
}
