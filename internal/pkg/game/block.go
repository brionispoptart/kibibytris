package game

import (
	"image/color"
	"math"
	"math/rand"
	"time"
)

type Block struct {
	X, Y      int
	Size      int
	Value     int
	Color     color.Color
	Timestamp time.Time
}

func NewBlock(x, y, value int) *Block {
	return &Block{
		X:         x,
		Y:         y,
		Size:      Config.BlockSize,
		Value:     value,
		Color:     ValueToPastelColor(value), // Assuming ValueToPastelColor is now in utils.go
		Timestamp: time.Now(),
	}
}

func (game *Game) SpawnNewBlock() {
	game.MergeBlocks()
	if game.canSpawnNewBlock() {
		value := 1 << (rand.Intn(game.calculateMaxExponent()) + 1)
		xPos := Config.GridWidth / 2 * Config.BlockSize
		game.ActiveBlock = NewBlock(xPos, 0, value)
	} else {
		game.triggerLoseCondition()
	}
}

func (game *Game) canSpawnNewBlock() bool {
	midPoint := Config.GridWidth / 2
	return game.Grid[0][midPoint] == nil
}

func (game *Game) calculateMaxExponent() int {
	highestValue := game.highestBlockValue()
	return int(math.Log2(float64(highestValue))) - 1
}

func (game *Game) highestBlockValue() int {
	highestValue := 2
	for _, row := range game.Grid {
		for _, block := range row {
			if block != nil && block.Value > highestValue {
				highestValue = block.Value
			}
		}
	}
	return highestValue
}
