package main

import (
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

	zombies []axeZombie
)

type axeZombie struct{
	level 	int
	hp 			int
	x, y 		float64
	speed		float64
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
	spawnZombies()
}

type Game struct{}



func (g *Game) Update() error { //game logic

	tickCount++
	
	lightSaberX = float64 (player1InitX + 100)
	lightSaberY = float64 (0)


moveSpeed := 3.0
blockRange := 40.0

// MOVE RIGHT (D)
if ebiten.IsKeyPressed(ebiten.KeyD) &&
  !isBlocked(player1InitX, player1InitY, 1, 0, blockRange, zombies) {
  player1InitX += moveSpeed
}

// MOVE LEFT (A)
if ebiten.IsKeyPressed(ebiten.KeyA) &&
  !isBlocked(player1InitX, player1InitY, -1, 0, blockRange, zombies) {
  player1InitX -= moveSpeed
}

// DOWN (S)
if ebiten.IsKeyPressed(ebiten.KeyS) &&
  !isBlocked(player1InitX, player1InitY, 0, 1, blockRange, zombies) {
  player1InitY += moveSpeed
}

// UP (W)
if ebiten.IsKeyPressed(ebiten.KeyW) &&
  !isBlocked(player1InitX, player1InitY, 0, -1, blockRange, zombies) {
	player1InitY -= moveSpeed
}


for i := range zombies {
  
	zombies[i].x, zombies[i].y = enemyMovement(
    player1InitX,
  	player1InitY,
    zombies[i].x,
    zombies[i].y,
    zombies[i].speed,
  )

  hitRange := 80.0 // damage player if close
  if abs(zombies[i].x-player1InitX) < hitRange &&
  abs(zombies[i].y-player1InitY) < hitRange {
    if tickCount%150 == 0 {
      player1hp--
    }
  }
}

  fmt.Println("hp:", player1hp)

	
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
sprite := axeZombieSprites[frame]

	for _, z := range zombies {
		op := &ebiten.DrawImageOptions{}
 		op.GeoM.Translate(z.x, z.y)
  	screen.DrawImage(sprite, op)
	}
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
