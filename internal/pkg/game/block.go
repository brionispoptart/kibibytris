package game

import (
	"image/color"
	"time"
)

type Block struct {
	X, Y      int
	Size      int
	Value     int
	Color     color.Color
	Timestamp time.Time
}
