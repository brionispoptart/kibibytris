package game

import (
	"fmt"
	"image/color"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font/basicfont"
)

const blockSize = 25

var (
	fastFallSpeed   = 5.0
	normalFallSpeed = 15.0
)

type Position struct {
	X int
	Y int
}

type Game struct {
	Grid             [][]*Block
	Score            int
	GridWidth        int
	GridHeight       int
	ActiveBlock      *Block
	FallCounter      float64
	IsFastFalling    bool
	lastLeftPressed  bool
	lastRightPressed bool
}

func NewGame() *Game {
	gridWidth, gridHeight := 10, 20
	grid := make([][]*Block, gridHeight)
	for rowIndex := range grid {
		grid[rowIndex] = make([]*Block, gridWidth)
	}
	newGame := &Game{
		Grid:             grid,
		Score:            0,
		GridWidth:        gridWidth,
		GridHeight:       gridHeight,
		lastLeftPressed:  false,
		lastRightPressed: false,
	}
	newGame.SpawnNewBlock()
	return newGame
}

func (game *Game) SpawnNewBlock() {
	game.ActiveBlock = &Block{
		X:         game.GridWidth / 2 * blockSize,
		Y:         0,
		Size:      blockSize,
		Value:     2, // Starting value for a new block, adjust as necessary.
		Color:     color.RGBA{R: 204, G: 204, B: 255, A: 255},
		Timestamp: time.Now(),
	}
}

func (game *Game) Update() error {
	leftPressed := ebiten.IsKeyPressed(ebiten.KeyLeft)
	if leftPressed && !game.lastLeftPressed {
		newPosX := game.ActiveBlock.X - blockSize
		if newPosX >= 0 && game.Grid[game.ActiveBlock.Y/blockSize][newPosX/blockSize] == nil {
			game.ActiveBlock.X = newPosX
		}
	}
	game.lastLeftPressed = leftPressed

	rightPressed := ebiten.IsKeyPressed(ebiten.KeyRight)
	if rightPressed && !game.lastRightPressed {
		newPosX := game.ActiveBlock.X + blockSize
		if newPosX < game.GridWidth*blockSize && game.Grid[game.ActiveBlock.Y/blockSize][newPosX/blockSize] == nil {
			game.ActiveBlock.X = newPosX
		}
	}
	game.lastRightPressed = rightPressed

	game.IsFastFalling = ebiten.IsKeyPressed(ebiten.KeyDown)
	fallSpeed := normalFallSpeed
	if game.IsFastFalling {
		fallSpeed = fastFallSpeed
	}

	game.FallCounter++
	if game.FallCounter >= fallSpeed {
		if game.ActiveBlock != nil {
			nextY := game.ActiveBlock.Y + blockSize
			if nextY/blockSize < game.GridHeight && game.Grid[nextY/blockSize][game.ActiveBlock.X/blockSize] == nil {
				game.ActiveBlock.Y = nextY
			} else {
				game.Grid[game.ActiveBlock.Y/blockSize][game.ActiveBlock.X/blockSize] = game.ActiveBlock
				game.MergeBlocks() // Call the merge function after placing the block
				game.SpawnNewBlock()
			}
			game.FallCounter = 0
		}
	}

	return nil
}
func (game *Game) Draw(screen *ebiten.Image) {
	for rowIndex, row := range game.Grid {
		for colIndex, blk := range row {
			if blk != nil {
				// Draw the block
				rectImage := ebiten.NewImage(blk.Size, blk.Size)
				rectImage.Fill(blk.Color)
				options := &ebiten.DrawImageOptions{}
				blockX := float64(colIndex * blockSize)
				blockY := float64(rowIndex * blockSize)
				options.GeoM.Translate(blockX, blockY)
				screen.DrawImage(rectImage, options)

				// Draw the block's value
				valueStr := fmt.Sprintf("%v", blk.Value)
				text.Draw(screen, valueStr, basicfont.Face7x13, int(blockX)+blk.Size/2, int(blockY)+blk.Size/2, color.Black)
			}
		}
	}
	if game.ActiveBlock != nil {
		activeBlockImage := ebiten.NewImage(game.ActiveBlock.Size, game.ActiveBlock.Size)
		activeBlockImage.Fill(game.ActiveBlock.Color)
		options := &ebiten.DrawImageOptions{}
		options.GeoM.Translate(float64(game.ActiveBlock.X), float64(game.ActiveBlock.Y))
		screen.DrawImage(activeBlockImage, options)
	}
}

func (game *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return game.GridWidth * blockSize, game.GridHeight * blockSize
}
