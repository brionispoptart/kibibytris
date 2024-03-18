package game

import (
	"fmt"
	"image/color"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font/basicfont"
)

const blockSize = 25

var (
	fastFallSpeed   = 3.5
	normalFallSpeed = 7.5
	backgroundColor = color.RGBA{R: 240, G: 234, B: 214, A: 255}
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
	GridWidth        int
	GridHeight       int
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
	rand.Seed(time.Now().UnixNano())
	// Define the maximum number of blocks to generate per row
	blocksPerRow := 3 // For example, up to 3 blocks per row

	// Initialize random blocks in the last three rows
	for y := gridHeight - 6; y < gridHeight; y++ {
		for i := 0; i < blocksPerRow; i++ {
			x := rand.Intn(gridWidth)
			if grid[y][x] == nil { // Check if the position is already occupied
				value := 1 << (rand.Intn(6) + 1) // Values between 2^1 (2) and 2^6 (64)
				grid[y][x] = &Block{
					X:         x * blockSize,
					Y:         y * blockSize,
					Size:      blockSize,
					Value:     value,
					Color:     valueToPastelColor(value),
					Timestamp: time.Now(),
				}
			}
		}
	}
	newGame.MergeBlocks()
	newGame.SpawnNewBlock()
	newGame.ComboMultiplier = 1 // Start with no multiplier
	newGame.ChainReactions = 0  // No chain reactions at the start
	return newGame
}
func (game *Game) triggerLoseCondition() {
	game.HasLost = true
}

func (game *Game) Update() error {
	game.ComboMultiplier = 1
	game.ChainReactions = 0
	if (game.HasLost || game.HasWon) && ebiten.IsKeyPressed(ebiten.KeySpace) {

		*game = *NewGame()
		return nil
	}
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
				game.MergeBlocks()   // Process all merges.
				game.SpawnNewBlock() // Spawn a new block only after merges complete.
			}
			game.FallCounter = 0
		}
	}

	return nil
}
func (game *Game) Draw(screen *ebiten.Image) {

	// Set the background to white
	screen.Fill(backgroundColor)

	face := basicfont.Face7x13
	for rowIndex, row := range game.Grid {
		for colIndex, blk := range row {
			if blk != nil {
				// Calculate positions
				blockX := float64(colIndex * blockSize)
				blockY := float64(rowIndex * blockSize)

				// Draw a border (darker color for contrast)
				borderImage := ebiten.NewImage(blk.Size, blk.Size)
				borderColor := darkenColor(blk.Color) // Assumes you have a darkenColor function
				borderImage.Fill(borderColor)
				borderOptions := &ebiten.DrawImageOptions{}
				borderOptions.GeoM.Translate(blockX, blockY)
				screen.DrawImage(borderImage, borderOptions)

				// Draw the inner block slightly smaller to create the stroke effect
				innerSize := blk.Size - 4 // Adjust the stroke width
				innerImage := ebiten.NewImage(innerSize, innerSize)
				innerImage.Fill(blk.Color) // Use the block's actual color
				innerOptions := &ebiten.DrawImageOptions{}
				innerOptions.GeoM.Translate(blockX+2, blockY+2) // Adjust to center the inner block within the border
				screen.DrawImage(innerImage, innerOptions)

				// Draw the block's value centered
				valueStr := fmt.Sprintf("%d", blk.Value)
				text.Draw(screen, valueStr, face, int(blockX)+blk.Size/2-(len(valueStr)*7)/2, int(blockY)+blk.Size/2+(7/2), color.Black)
			}
		}
	}
	if game.ActiveBlock != nil {
		// Update the active block's color based on its value
		game.ActiveBlock.Color = valueToPastelColor(game.ActiveBlock.Value)

		// Draw the active block with border and inner color
		activeBlockImage := ebiten.NewImage(game.ActiveBlock.Size, game.ActiveBlock.Size)
		activeBorderColor := darkenColor(game.ActiveBlock.Color)
		activeBlockImage.Fill(activeBorderColor)
		options := &ebiten.DrawImageOptions{}
		options.GeoM.Translate(float64(game.ActiveBlock.X), float64(game.ActiveBlock.Y))
		screen.DrawImage(activeBlockImage, options)

		// Inner block for the active block
		innerActiveBlockImage := ebiten.NewImage(game.ActiveBlock.Size-4, game.ActiveBlock.Size-4)
		innerActiveBlockImage.Fill(game.ActiveBlock.Color)
		innerOptions := &ebiten.DrawImageOptions{}
		innerOptions.GeoM.Translate(float64(game.ActiveBlock.X)+2, float64(game.ActiveBlock.Y)+2)
		screen.DrawImage(innerActiveBlockImage, innerOptions)

		// Draw the falling block's value
		valueStr := fmt.Sprintf("%d", game.ActiveBlock.Value)
		x := game.ActiveBlock.X + game.ActiveBlock.Size/2 - (len(valueStr)*7)/2 // Centering the text on the block
		y := game.ActiveBlock.Y + game.ActiveBlock.Size/2 + 7/2                 // Adjust for centering text vertically
		text.Draw(screen, valueStr, face, x, y, color.Black)
	}
	scoreStr := fmt.Sprintf("Score: %d", game.Score)
	x := 20 // For example, 20 pixels from the left
	y := 40 // For example, 40 pixels from the top, adjust as needed
	text.Draw(screen, scoreStr, basicfont.Face7x13, x, y, color.Black)

	if game.HasWon {
		// Dim the background to make the win message stand out
		overlay := ebiten.NewImage(500, 1000)                   // Assuming 500x1000 is your screen size
		overlay.Fill(color.RGBA{R: 255, G: 255, B: 255, A: 50}) // Semi-transparent white overlay
		screen.DrawImage(overlay, nil)

		// Display the win message
		msg := "You Win!"
		x := (250 - len(msg)*7) / 2 // Adjusted for basicfont.Face7x13
		y := 500 / 2
		text.Draw(screen, msg, basicfont.Face7x13, x, y, color.RGBA{R: 0, G: 0, B: 0, A: 255})
	}
	if game.HasLost {
		// Dim the background to make the win message stand out
		overlay := ebiten.NewImage(500, 1000)                   // Assuming 500x1000 is your screen size
		overlay.Fill(color.RGBA{R: 255, G: 255, B: 255, A: 50}) // Semi-transparent white overlay
		screen.DrawImage(overlay, nil)

		// Display the win message
		msg := "Really? Have you never\nplayed Tetris before?"
		x := (240 - len(msg)*3) / 2 // Adjusted for basicfont.Face7x13
		y := 500 / 2
		text.Draw(screen, msg, basicfont.Face7x13, x, y, color.RGBA{R: 0, G: 0, B: 0, A: 255})
	}
}

func darkenColor(c color.Color) color.Color {
	// This is a simple approach. For more accurate color manipulation, consider using a library like go-colorful.
	r, g, b, a := c.RGBA()
	factor := 0.7 // Adjust for desired darkness; closer to 0 is darker
	return color.RGBA{
		R: uint8(float64(r) * factor),
		G: uint8(float64(g) * factor),
		B: uint8(float64(b) * factor),
		A: uint8(a),
	}
}
func (game *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return game.GridWidth * blockSize, game.GridHeight * blockSize
}
