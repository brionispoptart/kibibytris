package game

import (
	"fmt"
	"image/color"
	"log"
	"os"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

var face font.Face
var (
	backgroundColor = color.RGBA{R: 240, G: 234, B: 214, A: 255}
)

func init() {
	fontBytes, err := os.ReadFile("internal/assets/font/early_gameboy.ttf")
	if err != nil {
		log.Fatalf("error reading font file: %v", err)
	}

	tt, err := opentype.Parse(fontBytes)
	if err != nil {
		log.Fatalf("error parsing font file: %v", err)
	}

	const dpi = 72
	face, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    8, // This could potentially be made configurable as well
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatalf("error creating font face: %v", err)
	}
}

func (game *Game) Draw(screen *ebiten.Image) {
	drawBackground(screen, backgroundColor)
	game.drawAllBlocks(screen)
	drawScore(screen, game.Score)

	if game.HasWon {
		drawOverlayMessage(screen, "You did what you were supposed to do,\nand your score is "+fmt.Sprint(game.Score)+"\n\nPress 'Space' to start again")
	}

	if game.HasLost {
		drawOverlayMessage(screen, "Really? Have you never\nplayed Tetris before?\n\nPress 'Space'\nto start again")
	}
}

func drawBackground(screen *ebiten.Image, bgColor color.Color) {
	screen.Fill(bgColor)
}

func (game *Game) drawAllBlocks(screen *ebiten.Image) {
	for rowIndex, row := range game.Grid {
		for colIndex, blk := range row {
			if blk != nil {
				game.drawBlock(screen, blk, colIndex, rowIndex)
			}
		}
	}
	if game.ActiveBlock != nil {
		game.drawBlock(screen, game.ActiveBlock, game.ActiveBlock.X/Config.BlockSize, game.ActiveBlock.Y/Config.BlockSize)
	}
}

func (game *Game) drawBlock(screen *ebiten.Image, blk *Block, colIndex, rowIndex int) {
	blockX, blockY := float64(colIndex*Config.BlockSize), float64(rowIndex*Config.BlockSize)
	borderColor := game.darkenColor(blk.Color) // Assuming darkenColor is moved to utils.go and exported as DarkenColor
	borderImage, innerImage := createBlockImages(blk.Size, borderColor, blk.Color)
	drawImageWithBorder(screen, borderImage, innerImage, blockX, blockY)
	drawBlockValue(screen, blk.Value, blockX, blockY, borderColor)
}

func createBlockImages(size int, borderColor, innerColor color.Color) (borderImage, innerImage *ebiten.Image) {
	borderImage = ebiten.NewImage(size, size)
	borderImage.Fill(borderColor)

	innerSize := size - 4 // Adjust the stroke width
	innerImage = ebiten.NewImage(innerSize, innerSize)
	innerImage.Fill(innerColor)

	return
}

func drawImageWithBorder(screen, borderImage, innerImage *ebiten.Image, x, y float64) {
	borderOptions := &ebiten.DrawImageOptions{}
	borderOptions.GeoM.Translate(x, y)
	screen.DrawImage(borderImage, borderOptions)

	innerOptions := &ebiten.DrawImageOptions{}
	innerOptions.GeoM.Translate(x+2, y+2) // Adjust to center the inner block within the border
	screen.DrawImage(innerImage, innerOptions)
}

func drawBlockValue(screen *ebiten.Image, value int, x, y float64, borderColor color.Color) {
	valueStr := fmt.Sprintf("%d", value)
	text.Draw(screen, valueStr, face, int(x)+Config.BlockSize/2-(len(valueStr)*7)/2, int(y)+Config.BlockSize/2+(7/2), borderColor)
}

func drawScore(screen *ebiten.Image, score int) {
	scoreStr := fmt.Sprintf("Score: %d", score)
	text.Draw(screen, scoreStr, face, 20, 40, color.Black)
}

func drawOverlayMessage(screen *ebiten.Image, msg string) {
	overlay := ebiten.NewImage(450, 900)                     // Update dimensions if needed
	overlay.Fill(color.RGBA{R: 240, G: 234, B: 214, A: 255}) // Semi-transparent white overlay

	screen.DrawImage(overlay, nil)

	lines := strings.Split(msg, "\n")

	textHeight := len(lines) * face.Metrics().Height.Ceil()

	y := (screen.Bounds().Dy() - textHeight) / 2

	for _, line := range lines {
		bound, _ := font.BoundString(face, line)
		width := (bound.Max.X - bound.Min.X).Ceil()

		x := (screen.Bounds().Dx() - width) / 2

		text.Draw(screen, line, face, x, y, color.Black)

		y += face.Metrics().Height.Ceil()
	}
}

func (game *Game) darkenColor(c color.Color) color.Color {
	r, g, b, a := c.RGBA()
	factor := .7 // Adjust for desired darkness; closer to 0 is darker
	return color.RGBA{
		R: uint8(float64(r) * factor),
		G: uint8(float64(g) * factor),
		B: uint8(float64(b) * factor),
		A: uint8(a),
	}
}
