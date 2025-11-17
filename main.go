package main

import (
	"math/rand"
	"log"
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

var ( //declvare variable for images, name *ebiten.Image.
	background *ebiten.Image
	player1 *ebiten.Image

	axeZombieSprites []*ebiten.Image

	lightSaber *ebiten.Image

	screenHeight = 1080
	screenWidth = 1920

	player1InitX = float64(560)
	player1InitY = float64(240)
	axeZombieInitXTemp = float64 (randFloat(1,100))
	axeZombieInitYTemp = float64 (randFloat(1,100))
	lightSaberX float64
	lightSaberY float64

	player1hp = 20
	tickCount = 0
)

type axeZombie struct{

}

func enemyMovement(targetX, targetY, enemyX, enemyY, speed float64) (float64, float64) {
	if enemyX < (targetX - 80){ //enemie movement
		enemyX += speed
	}
	if enemyX > (targetX + 80){
		enemyX -= speed
	}
	if enemyY < (targetY - 80){
		enemyY += speed
	}
	if enemyY > (targetY + 80){
		enemyY -= speed
	}

	return enemyX, enemyY
}

func loadAxeZombieSprites() {
  axeZombieSprites = make([]*ebiten.Image, 8)

  for i := 1; i <= 8; i++ {
    filename := fmt.Sprintf("assets/sprites/enemies/axeZombie/axeZombieSprite%02d.png", i)

    img, _, err := ebitenutil.NewImageFromFile(filename)
    if err != nil {
    log.Fatal(err)
    }
    axeZombieSprites[i-1] = img
	}
}

func abs(f float64) float64 {
  if f < 0 {
    return -f
  }
  return f
}

func randFloat(min, max float64) float64 {
  return min + rand.Float64()*(max-min)
}

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

	lightSaber, _, err = ebitenutil.NewImageFromFile("assets/images/lightSaber.png") //will not run if empty
	if err != nil {
		log.Fatal(err)
	}

	loadAxeZombieSprites()
}

type Game struct{}

func (g *Game) Update() error { //game logic

	tickCount++
	
	lightSaberX = float64 (player1InitX + 100)
	lightSaberY = float64 (0)

	//any movement code cannot be a switch because it will prevent diagnol
	if ebiten.IsKeyPressed(ebiten.KeyD) == true { //player movement
		player1InitX = player1InitX + 3
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) == true {
		player1InitX = player1InitX - 3
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) == true {
		player1InitY = player1InitY + 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyW) == true {
		player1InitY = player1InitY - 3
	} 

	axeZombieInitXTemp, axeZombieInitYTemp = enemyMovement(
    player1InitX,
    player1InitY,
    axeZombieInitXTemp,
    axeZombieInitYTemp,
		0.5,
	)

	// enemy damage when close enough
	hitRange := 100.0 // adjust to taste

	if abs(axeZombieInitXTemp - player1InitX) < hitRange && abs(axeZombieInitYTemp - player1InitY) < hitRange {
	if tickCount % 150 == 0 {
		player1hp--
	}
  fmt.Println("hp:", player1hp)
}
	
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {  //called every frame, graphics.
	
	screen.DrawImage(background, nil)

	op := &ebiten.DrawImageOptions{}	
	opAxeZombie := &ebiten.DrawImageOptions{}

	op.GeoM.Translate(player1InitX,player1InitY)
	opAxeZombie.GeoM.Translate(axeZombieInitXTemp,axeZombieInitYTemp)		

	screen.DrawImage(player1, op)	

	opLightSaber := &ebiten.DrawImageOptions{} //todo: fix
	opLightSaber.GeoM.Translate(lightSaberX, lightSaberY)

	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) == true {
		screen.DrawImage(lightSaber, opLightSaber)
	}

	frame := (tickCount / 8) % len(axeZombieSprites)
	currentSprite := axeZombieSprites[frame]
	screen.DrawImage(currentSprite, opAxeZombie)

}


func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 640, 480
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Render an image")
	
	if err := ebiten.RunGame(&Game{}); err != nil { 
		log.Fatal(err)
	}	
}
