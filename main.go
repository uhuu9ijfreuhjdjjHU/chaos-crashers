package main

import (	
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

var ( //declvare variable for images, name *ebiten.Image.
	background *ebiten.Image
	player1 * ebiten.Image
	//charecter2 * ebiten.Image
	screenWidth = 750
	screenHeight = 750
	player1InitW = (screenWidth / 2) + 50
	player1InitH = (screenHeight / 2) + 50
	//player2InitW = (screenWidth / 2) - 50
	//player2InitH = (screenHeight / 2) - 50
)

func init() { //initialize images to variables here.
	var err error
	background, _, err = ebitenutil.NewImageFromFile("assets/images/go.png") //name, _, etc.
	if err != nil {
		log.Fatal(err)
	}

	//player1, _, err = ebitenutil.NewImageFromFile("") //will not run if empty
	if err != nil {
		log.Fatal(err)
	}
}

type Game struct{}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.DrawImage(background, nil)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 640, 480
}

func main() {
	ebiten.SetWindowSize(screenHeight, screenWidth)
	ebiten.SetWindowTitle("Render an image")
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
