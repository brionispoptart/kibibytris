package game

import (
	"image/color"
	"math"
	"math/rand"
	"time"

	"github.com/lucasb-eyer/go-colorful"
)

func makeGrid(width, height int) [][]*Block {
	grid := make([][]*Block, height)
	for i := range grid {
		grid[i] = make([]*Block, width)
	}
	return grid
}

func (game *Game) populateInitialBlocks() {
	rand.Seed(time.Now().UnixNano())
	blocksPerRow := 3
	for y := Config.GridHeight - 6; y < Config.GridHeight; y++ {
		for i := 0; i < blocksPerRow; i++ {
			x := rand.Intn(Config.GridWidth)
			if game.Grid[y][x] == nil {
				value := 1 << (rand.Intn(6) + 1)
				game.Grid[y][x] = &Block{
					X:         x * Config.BlockSize,
					Y:         y * Config.BlockSize,
					Size:      Config.BlockSize,
					Value:     value,
					Color:     ValueToPastelColor(value), // This function needs to be implemented or referenced correctly.
					Timestamp: time.Now(),
				}
			}
		}
	}
}

func ValueToPastelColor(value int) color.Color {
	const saturation = 0.3
	const lightness = 0.9

	hue := math.Mod(float64(value*137), 360.0) // Prime number for hue variation
	c := colorful.Hsl(hue, saturation, lightness)
	r, g, b := c.RGB255()
	return color.RGBA{R: r, G: g, B: b, A: 255}
}
