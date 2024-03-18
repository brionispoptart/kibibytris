package game

import (
	"sort"
	"time"
)

// MergeBlocks initiates the merge process for all blocks in the grid.
func (game *Game) MergeBlocks() {
	mergeOccurred := true
	for mergeOccurred {
		game.populateMergeQueue()
		if len(game.MergeQueue) > 0 {
			game.executeMergeQueue()
			mergeOccurred = true
		} else {
			mergeOccurred = false
		}
	}
	game.checkWinCondition()
}

// populateMergeQueue checks each grid position for potential merges and queues them.
func (game *Game) populateMergeQueue() {
	game.MergeQueue = []MergeOperation{}
	for y := 0; y < game.GridHeight; y++ {
		for x := 0; x < game.GridWidth; x++ {
			game.checkAndQueueMerge(x, y)
		}
	}
}

// executeMergeQueue processes merges in the queue based on their timestamp.
func (game *Game) executeMergeQueue() {
	sort.Slice(game.MergeQueue, func(i, j int) bool {
		return game.MergeQueue[i].MergeTime.Before(game.MergeQueue[j].MergeTime)
	})
	for len(game.MergeQueue) > 0 {
		nextMerge := game.MergeQueue[0]
		game.MergeQueue = game.MergeQueue[1:]
		game.processMerge(nextMerge.Position.X, nextMerge.Position.Y)
	}
}

// checkWinCondition checks if any block has reached the win condition value.
func (game *Game) checkWinCondition() {
	for _, row := range game.Grid {
		for _, block := range row {
			if block != nil && block.Value == 2048 {
				game.HasWon = true
				return
			}
		}
	}
}

// checkAndQueueMerge checks if a merge is possible at the given position and queues it.
func (game *Game) checkAndQueueMerge(x, y int) {
	block := game.Grid[y][x]
	if block == nil {
		return
	}
	oldestNeighbor, _ := game.OldestNeighbor(x, y)
	if oldestNeighbor != nil && !game.isAlreadyQueued(x, y) {
		game.MergeQueue = append(game.MergeQueue, MergeOperation{Position: Position{X: x, Y: y}, MergeTime: block.Timestamp})
	}
}

// isAlreadyQueued checks if a merge operation is already in the queue.
func (game *Game) isAlreadyQueued(x, y int) bool {
	for _, op := range game.MergeQueue {
		if op.Position.X == x && op.Position.Y == y {
			return true
		}
	}
	return false
}

// processMerge handles the merging of two blocks.
func (game *Game) processMerge(x, y int) {
	block := game.Grid[y][x]
	if block == nil {
		return
	}
	neighbors := game.Neighbors(x, y)
	for _, pos := range neighbors {
		neighbor := game.Grid[pos.Y][pos.X]
		if neighbor != nil && neighbor.Value == block.Value {
			game.mergeBlocks(block, neighbor, pos.Y, pos.X)
			break // Stop after finding the first merge to prevent double merging.
		}
	}
}

// mergeBlocks merges two blocks, updates the grid, and assigns a new color.
func (game *Game) mergeBlocks(block1, block2 *Block, y2, x2 int) {
	mergedValue := block1.Value * 2 // Assuming blocks have equal value
	game.addScoreForMerge(mergedValue, game.ComboMultiplier, game.ChainReactions > 0)
	block1.Value *= 2
	block1.Timestamp = time.Now()
	block1.Color = valueToPastelColor(block1.Value)
	game.Grid[y2][x2] = nil // Remove the merged block.
	game.ChainReactions++
}

func (game *Game) addScoreForMerge(mergedValue, comboMultiplier int, isChainReaction bool) {
	baseScore := mergedValue
	if isChainReaction {
		baseScore *= 2 // Double points for chain reactions
	}
	game.Score += baseScore * comboMultiplier
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
