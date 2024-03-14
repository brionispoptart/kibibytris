package game

import (
	"time"
)

func (game *Game) MergeBlocks() {
	for y := 0; y < game.GridHeight; y++ {
		for x := 0; x < game.GridWidth; x++ {
			currentBlock := game.Grid[y][x]
			if currentBlock == nil {
				continue
			}

			// Check for neighboring blocks and merge into the oldest
			oldestNeighbor, oldestNeighborPos := game.OldestNeighbor(x, y)
			if oldestNeighbor != nil {
				// Merge current block into the oldest neighbor
				oldestNeighbor.Value *= 2
				game.Grid[y][x] = nil // Remove the current block
				game.Grid[oldestNeighborPos.Y][oldestNeighborPos.X] = oldestNeighbor
				game.MergeInto(oldestNeighborPos.X, oldestNeighborPos.Y, oldestNeighbor.Timestamp) // Check for further merges
			}
		}
	}
}

func (game *Game) OldestNeighbor(x, y int) (*Block, *Position) {
	directions := []Position{{X: 0, Y: -1}, {X: -1, Y: 0}, {X: 1, Y: 0}, {X: 0, Y: 1}}
	var oldestBlock *Block
	var oldestBlockPos *Position

	for _, dir := range directions {
		neighborX, neighborY := x+dir.X, y+dir.Y
		if neighborX >= 0 && neighborX < game.GridWidth && neighborY >= 0 && neighborY < game.GridHeight {
			neighbor := game.Grid[neighborY][neighborX]
			if neighbor != nil && neighbor.Value == game.Grid[y][x].Value {
				if oldestBlock == nil || neighbor.Timestamp.Before(oldestBlock.Timestamp) {
					oldestBlock = neighbor
					oldestBlockPos = &Position{X: neighborX, Y: neighborY}
				}
			}
		}
	}

	return oldestBlock, oldestBlockPos
}

func (game *Game) MergeInto(x, y int, mergeTime time.Time) {
	currentBlock := game.Grid[y][x]
	if currentBlock == nil || currentBlock.Timestamp.After(mergeTime) {
		return
	}

	neighborsToMerge := game.Neighbors(x, y)
	for _, pos := range neighborsToMerge {
		neighbor := game.Grid[pos.Y][pos.X]
		if neighbor.Timestamp.Before(currentBlock.Timestamp) {
			currentBlock.Value += neighbor.Value
			game.Grid[pos.Y][pos.X] = nil
			currentBlock.Timestamp = time.Now()
		} else {
			neighbor.Value += currentBlock.Value
			game.Grid[y][x] = nil
			game.MergeInto(pos.X, pos.Y, neighbor.Timestamp)
			return
		}
	}
}

func (game *Game) Neighbors(x, y int) []Position {
	var neighbors []Position
	directions := []Position{{X: 0, Y: -1}, {X: -1, Y: 0}, {X: 1, Y: 0}, {X: 0, Y: 1}}

	for _, dir := range directions {
		neighborX, neighborY := x+dir.X, y+dir.Y
		if neighborX >= 0 && neighborX < game.GridWidth && neighborY >= 0 && neighborY < game.GridHeight {
			neighbor := game.Grid[neighborY][neighborX]
			if neighbor != nil && neighbor.Value == game.Grid[y][x].Value {
				neighbors = append(neighbors, Position{X: neighborX, Y: neighborY})
			}
		}
	}

	return neighbors
}
