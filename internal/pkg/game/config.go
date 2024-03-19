package game

import "time"

var Config = struct {
	BlockSize     int
	GridWidth     int
	GridHeight    int
	BaseFallSpeed float64 // Represents blocks per frame at 60 fps for 1 block per second fall rate
	BoostSpeed    float64 // Boost speed, similarly calculated
	InitialDelay  time.Duration
	RepeatDelay   time.Duration
}{
	BlockSize:     25,
	GridWidth:     9,
	GridHeight:    18,
	BaseFallSpeed: 10.0 / 60.0, // Assuming 60 FPS, for a block to fall one blockSize per second
	BoostSpeed:    20.0 / 60.0, // Adjusted for faster fall during boost, ensuring it's more frames per second
	InitialDelay:  200 * time.Millisecond,
	RepeatDelay:   50 * time.Millisecond,
}
