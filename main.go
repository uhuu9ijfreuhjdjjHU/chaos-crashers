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
	sword *ebiten.Image

	axeZombieSprites []*ebiten.Image
	axeZombieHitSprites []*ebiten.Image	

	screenHeight = 1080
	screenWidth = 1920

	player1InitX = float64(560)
	player1InitY = float64(240)
	axeZombieInitXTemp = float64 (randFloat(1,100))
	axeZombieInitYTemp = float64 (randFloat(1,100))
	swordX float64
	swordY float64

	player1hp = 20

	tickCount = 0

	zombies []axeZombie

	swordLocation = rune ('s') //a = left, d = right, s = down, w = up
)

type Game struct{}

type axeZombie struct{
	level 	int
	hp 			int
	x, y 		float64
	speed		float64
	hit 		bool
	hitTimer int
	facingRight bool
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

	sword, _, err = ebitenutil.NewImageFromFile("assets/images/sword.png") //will not run if empty
	if err != nil {
		log.Fatal(err)
	}
	
	//~~> animation functions <~~\\
	loadAxeZombieSprites()
	loadAxeZombieHitSprites()


	spawnZombies(0.7) //loads zombies, condition changes zombie speed.
}

func (g *Game) Update() error { //game logic


	for i := range zombies {
    if zombies[i].hitTimer > 0 {
      zombies[i].hitTimer--
      zombies[i].hit = true
    } else {
      zombies[i].hit = false
    }
	}

	tickCount++

	//~~> sword direction logic <~~\\

	switch {

		case ebiten.IsKeyPressed(ebiten.KeyArrowRight):
			swordLocation = 'd'
		case ebiten.IsKeyPressed(ebiten.KeyArrowLeft):	
			swordLocation = 'a'
		case ebiten.IsKeyPressed(ebiten.KeyArrowDown):	
			swordLocation = 's'
		case ebiten.IsKeyPressed(ebiten.KeyArrowUp):
			swordLocation = 'w'
	}

	switch {
		case swordLocation == 'd':
			swordX = float64 (player1InitX + 100)
			swordY = float64 (player1InitY)
		case swordLocation == 'a':
			swordX = float64 (player1InitX - 100)
			swordY = float64 (player1InitY)
		case swordLocation == 's':
			swordX = float64 (player1InitX)
			swordY = float64 (player1InitY + 100)
		case swordLocation == 'w':
			swordX = float64 (player1InitX)
			swordY = float64 (player1InitY - 100)
	}
	
	moveSpeed := 3.0
	blockRange := 35.0 //collusion stat

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

	if zombies[i].hp <= 0 {
		continue
	}
	
   // movement (once per zombie)
  zombies[i].x, zombies[i].y = enemyMovement(
  	player1InitX,
    player1InitY,
    zombies[i].x,
    zombies[i].y,
    zombies[i].speed,
		3,
		swordLocation,
    zombies,
    i,
  )

  //~~> player damage check <~~\\
  hitRange := 80.0
  
	if abs(zombies[i].x-player1InitX) < hitRange && abs(zombies[i].y-player1InitY) < hitRange {

  	if tickCount%150 == 0 {
    	player1hp--
      fmt.Println("hp:", player1hp)
    }
  }

  //~~> sword hit detection <~~\\
  swordHitRange := 30.0
  if abs(zombies[i].x - swordX) < swordHitRange && abs(zombies[i].y - swordY) < swordHitRange {
  if (tickCount % 45 == 0){ //attack speed
		zombies[i].hp--
		zombies[i].hit = true
		zombies[i].hitTimer = 8
  	fmt.Println("Zombie", i, "hp:", zombies[i].hp)
		}
  }
}
	
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {  //called every frame, graphics.
	
	screen.DrawImage(background, nil)

	op := &ebiten.DrawImageOptions{}
	opAxeZombie := &ebiten.DrawImageOptions{}

	op.GeoM.Translate(player1InitX,player1InitY)
	opAxeZombie.GeoM.Translate(axeZombieInitXTemp,axeZombieInitYTemp)		

	opSword := &ebiten.DrawImageOptions{}
	opSword.GeoM.Translate(swordX, swordY)

	screen.DrawImage(player1, op)	

	screen.DrawImage(sword, opSword)

	frame := (tickCount / 8) % len(axeZombieSprites)
	axeZombieSpriteFrame := axeZombieSprites[frame]
	axeZombieHitSpriteFrame := axeZombieHitSprites[frame]

	for _, z := range zombies {
		if z.hp <= 0 {
			continue
		}
		
		if z.hit == true {	
			op := &ebiten.DrawImageOptions{}
 			op.GeoM.Translate(z.x, z.y)
  		screen.DrawImage(axeZombieHitSpriteFrame, op)
		} else {
			if !z.facingRight{
			op := &ebiten.DrawImageOptions{}
			w := float64(axeZombieSpriteFrame.Bounds().Dx())

			op.GeoM.Scale(-1,1)
 			op.GeoM.Translate(z.x + w, z.y)
  		screen.DrawImage(axeZombieSpriteFrame, op)
			} else {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Scale(1,1) 
 			op.GeoM.Translate(z.x, z.y)
  		screen.DrawImage(axeZombieSpriteFrame, op)
			}
		}
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
