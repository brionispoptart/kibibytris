package game

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

var (
	lastMoveTime = time.Now()
	moveDelay    = Config.InitialDelay
)

func (game *Game) handleMovement(key ebiten.Key, direction int) {
	now := time.Now()
	pressed := ebiten.IsKeyPressed(key)

	if pressed {
		if direction > 0 && !game.lastRightPressed || direction < 0 && !game.lastLeftPressed {
			game.tryMove(direction)
			lastMoveTime = now
			moveDelay = Config.InitialDelay
		} else if now.Sub(lastMoveTime) > moveDelay {
			game.tryMove(direction)
			moveDelay = Config.RepeatDelay
			lastMoveTime = now
		}
	}

	if direction > 0 {
		game.lastRightPressed = pressed
	} else {
		game.lastLeftPressed = pressed
	}
}

func (game *Game) tryMove(direction int) {
	newPosX := game.ActiveBlock.X + direction*Config.BlockSize
	if newPosX >= 0 && newPosX < Config.GridWidth*Config.BlockSize && game.Grid[game.ActiveBlock.Y/Config.BlockSize][newPosX/Config.BlockSize] == nil {
		game.ActiveBlock.X = newPosX
	}
}

func (game *Game) MoveRight() {
	game.handleMovement(ebiten.KeyRight, 1)
}

func (game *Game) MoveLeft() {
	game.handleMovement(ebiten.KeyLeft, -1)
}

func (game *Game) Boost() {
	game.CurrentFallSpeed = Config.BaseFallSpeed
	game.IsFastFalling = ebiten.IsKeyPressed(ebiten.KeyDown)
	if game.IsFastFalling {
		game.CurrentFallSpeed = Config.BoostSpeed
	} else {
		game.CurrentFallSpeed = Config.BaseFallSpeed
	}
}
func (game *Game) isDirectionKeyPressed() bool {
	return ebiten.IsKeyPressed(ebiten.KeyLeft) || ebiten.IsKeyPressed(ebiten.KeyRight)
}
func (game *Game) finalizeBlockPosition() {
	game.Grid[game.ActiveBlock.Y/Config.BlockSize][game.ActiveBlock.X/Config.BlockSize] = game.ActiveBlock
	game.MergeBlocks()   // Handle merging.
	game.SpawnNewBlock() // Spawn the next block.
}

func (game *Game) StartGame() {
	if (game.HasLost || game.HasWon) && ebiten.IsKeyPressed(ebiten.KeySpace) {
		*game = *NewGame()
	}
}
