package game

import (
	"time"
)

type Position struct {
	X int
	Y int
}

type MergeOperation struct {
	Position  Position
	MergeTime time.Time
}

type Game struct {
	Grid             [][]*Block
	ActiveBlock      *Block
	FallCounter      float64
	IsFastFalling    bool
	lastLeftPressed  bool
	lastRightPressed bool
	MergeQueue       []MergeOperation
	HasWon           bool
	HasLost          bool
	Score            int
	ComboMultiplier  int
	ChainReactions   int
	CurrentFallSpeed float64
}

func NewGame() *Game {
	game := &Game{
		CurrentFallSpeed: Config.BaseFallSpeed,
		Grid:             makeGrid(Config.GridWidth, Config.GridHeight),
	}
	game.populateInitialBlocks()
	game.resetGameState()
	game.SpawnNewBlock()
	return game
}
func (game *Game) resetGameState() {
	game.Score = 0
	game.ComboMultiplier = 1
	game.ChainReactions = 0
	game.HasWon = false
	game.HasLost = false
}

func (game *Game) triggerLoseCondition() {
	game.HasLost = true
}

func (game *Game) handleControls() {
	game.MoveLeft()
	game.MoveRight()
	game.Boost()
	game.StartGame()
}

func (game *Game) Update() error {
	game.ComboMultiplier = 1
	game.ChainReactions = 0

	game.handleControls()

	game.FallCounter += game.CurrentFallSpeed

	if game.FallCounter >= 1.0 {

		if game.ActiveBlock != nil {
			nextY := game.ActiveBlock.Y + Config.BlockSize
			if nextY/Config.BlockSize < Config.GridHeight && game.Grid[nextY/Config.BlockSize][game.ActiveBlock.X/Config.BlockSize] == nil {
				game.ActiveBlock.Y = nextY
			} else if !game.isDirectionKeyPressed() {
				game.finalizeBlockPosition()
			}
			game.FallCounter = 0
		}
	}

	return nil
}

func (game *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return Config.GridWidth * Config.BlockSize, Config.GridHeight * Config.BlockSize
}
