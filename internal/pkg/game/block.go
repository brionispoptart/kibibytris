package game

import (
	"image/color"
	"math"
	"math/rand"
	"time"

	"github.com/lucasb-eyer/go-colorful"
)

type Block struct {
	X, Y      int
	Size      int
	Value     int
	Color     color.Color
	Timestamp time.Time
}

func valueToPastelColor(value int) color.Color {
	// Adjust these ranges to fine-tune the appearance of the pastel colors
	const saturationRange = 0.4 // Pastel-like saturation
	const lightness = 0.8       // High lightness for pastel effect

	// Use the block's value to calculate a hue. This is a simple modulo operation,
	// but you could use more complex functions for different mappings.
	// Here we assume value is always a power of 2.
	hue := math.Mod(float64(value*137), 360.0) // 137 is an arbitrary prime to scatter hues

	// Generate the pastel color from the HSL space
	c := colorful.Hsl(hue, saturationRange, lightness)

	// Convert to RGBA
	r, g, b := c.RGB255()
	return color.RGBA{R: r, G: g, B: b, A: 255}
}

func (game *Game) SpawnNewBlock() {
	// Attempt to merge blocks before spawning a new one.
	game.MergeBlocks()

	// Check if it's possible to spawn a new block.
	if game.canSpawnNewBlock() {
		// Calculate properties for the new block.
		_, _, exponent := game.calculateSpawnBlockProperties()

		// Spawn the new block at the grid's center top.
		game.ActiveBlock = &Block{
			X:         game.GridWidth / 2 * blockSize,
			Y:         0,
			Size:      blockSize,
			Value:     1 << exponent,
			Color:     valueToPastelColor(1 << exponent),
			Timestamp: time.Now(),
		}
	} else {
		// If a new block can't be spawned, trigger the lose condition.
		game.triggerLoseCondition()
	}
}

func (game *Game) canSpawnNewBlock() bool {
	return game.Grid[0][game.GridWidth/2] == nil
}

// calculateSpawnBlockProperties calculates the properties for the new block to spawn.
func (game *Game) calculateSpawnBlockProperties() (int, int, int) {
	// Similar to your previous implementation or use the previously provided logic.
	highestValue := 2
	for _, row := range game.Grid {
		for _, block := range row {
			if block != nil && block.Value > highestValue {
				highestValue = block.Value
			}
		}
	}
	maxExponent := int(math.Log2(float64(highestValue)))
	rand.Seed(time.Now().UnixNano())
	exponent := rand.Intn(maxExponent) + 1

	return highestValue, maxExponent, exponent
}
