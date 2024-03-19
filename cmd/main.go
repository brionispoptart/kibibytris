package main

import (
	"log"

	"github.com/brionispoptart/kibibytris/internal/pkg/game" // Adjusted import path
	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	ebiten.SetWindowSize(450, 900)
	ebiten.SetWindowTitle("kibibytris")

	gameInstance := game.NewGame() // Corrected to reference the game package

	if err := ebiten.RunGame(gameInstance); err != nil {
		log.Fatal(err)
	}
}
