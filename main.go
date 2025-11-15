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
	player1InitX = float64((screenWidth / 2) + 50)
	player1InitY = float64((screenHeight / 2) + 50)
	//player2InitX = (screenWidth / 2) - 50
	//player2InitY = (screenHeight / 2) - 50
)

func init() { //initialize images to variables here.
	var err error
	background, _, err = ebitenutil.NewImageFromFile("assets/images/go.png") //name, _, etc.
	if err != nil {
		log.Fatal(err)
	}

	player1, _, err = ebitenutil.NewImageFromFile("assets/images/Sprite-0001.png") //will not run if empty
	if err != nil {
		log.Fatal(err)
	}
}

type Game struct{}

func (g *Game) Update() error {
	
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) == true {
		player1InitX = player1InitX + 2
	}

	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) == true {
		player1InitX = player1InitX - 2
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {  //called every frame, game tick rate
	screen.DrawImage(background, nil)

	x := player1InitX
	y := player1InitY

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(x,y)

	screen.DrawImage(player1, op)
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
