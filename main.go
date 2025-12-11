package main

import (
	"log"
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"math"
)

var ( //declvare variable for images, name *ebiten.Image.
	background *ebiten.Image
	player1 *ebiten.Image
	swordSprites []*ebiten.Image
	axeZombieDeathSprites []*ebiten.Image

	axeZombieSprites []*ebiten.Image //an array of image files means it for a animation
	axeZombieHitSprites []*ebiten.Image	//see functions.go

	screenHeight = 1080
	screenWidth = 1920

	player1InitX = float64(560)
	player1InitY = float64(240)

	axeZombieInitXTemp = float64 (randFloat(1,100))
	axeZombieInitYTemp = float64 (randFloat(1,100))
	swordX float64
	swordY float64
	
	//lower is faster
	axeZombieAnimationSpeed = float64(10)
	axeZombieHitAnimationSpeed = float64(5)
	//higher is faster
	axeZombieLiteralSpeed = float64(0.7)

	tickCount = 0 //for game time keeping

	zombies []axeZombie

	player1hp = 20
	hitFrameDuration = int(0)
	playerAttackCount = int (0)
	playerAttackFrames = int(15) //frame length of player attack. hit frame duration will call to this at runtime, do not use magic numbers.
	playerAttackFramesTimer = int(0)		
	swordLocation = rune ('s') //a = left, d = right, s = down, w = up
	playerAttackActive = bool(false)
	playerAttackFlipped = bool(false)
	playerAttackFramesStart = bool(false)
)

type Game struct{}

type axeZombie struct {
	level         int
	hp            int
	x, y          float64
	speed         float64
	hit           bool
	hitTimer      int
	hitFrame      int
	walkFrame     int
	facingRight   bool
	invulnerable  bool
	walkTimer     float64 
	hitAnimTimer  float64  
	inHitAnimation bool
	deathAnimationPlayed bool
	deathAnimationTimer float64
	deathAnimationFrame int
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
	
	loadAxeZombieSprites() //call animation functions here
	loadAxeZombieHitSprites()
	loadSwordSprites()
	loadAxeZombieDeathSprites()

	spawnAxeZombies(axeZombieLiteralSpeed) //loads zombies, condition changes zombie speed.
}

func (g *Game) Update() error { //game logic

	tickCount++
	zombieWalkCycleUpdate(axeZombieAnimationSpeed)
	zombieHitAnimationUpdate(axeZombieHitAnimationSpeed)	
	zombieDeathAnimationUpdate(3)


	if tickCount % 60 == 0 { //prints every 60 frames for time keeping.
		fmt.Println("frame", tickCount, ",", "RAM: ", GetSelfRAM(), "MB")
		for i := range zombies {
			fmt.Println("axe zombie" ,i ," frame: ", zombies[i].walkFrame)
		}
	}

	for i := range zombies { //keeps track of how long zombies should be "hit" for
    if zombies[i].hitTimer > 0 {
      zombies[i].hitTimer--
      zombies[i].hit = true
    } else if zombies[i].hitTimer == 0 {
      zombies[i].hit = false
    }
	}

	if hitFrameDuration == 0 { // prevents player from attacking same enemy.
		for i := range zombies {
			if zombies[i].hitTimer == 0 {
				zombies[i].invulnerable = false
				playerAttackActive = false
			}
		}
	}

	//~~> sword direction logic <~~\\

	switch { //player sword controls
		case ebiten.IsKeyPressed(ebiten.KeyArrowRight) && hitFrameDuration == 0:
			swordLocation = 'd'
			hitFrameDuration = playerAttackFrames
			playerAttackFramesStart = true
		case ebiten.IsKeyPressed(ebiten.KeyArrowLeft) && hitFrameDuration == 0:
			swordLocation = 'a'
			hitFrameDuration = playerAttackFrames
			playerAttackFramesStart = true
		case ebiten.IsKeyPressed(ebiten.KeyArrowDown) && hitFrameDuration == 0:
			swordLocation = 's'
			hitFrameDuration = playerAttackFrames
			playerAttackFramesStart = true
		case ebiten.IsKeyPressed(ebiten.KeyArrowUp) && hitFrameDuration == 0:
			swordLocation = 'w'
			hitFrameDuration = playerAttackFrames
			playerAttackFramesStart = true	
	}

	switch { //player sword direction logic, effected by player sword controls above
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
	
	moveSpeed := 3.0 //player move speed
	blockRange := 35.0 //player collusion stat

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

	for i := range zombies { //zombie ai / logic

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
  	
		if abs(zombies[i].x-player1InitX) < hitRange && 
		abs(zombies[i].y-player1InitY) < hitRange && tickCount % 150 == 0 {	
    	player1hp--
    	fmt.Println("hp:", player1hp) 
  	}

		swordHitRange := 30.0

  	if abs(zombies[i].x - swordX) < swordHitRange && !zombies[i].invulnerable &&
		abs(zombies[i].y - swordY) < swordHitRange && playerAttackActive && 
		zombies[i].hitTimer <= 0 {
			zombies[i].hp--
			zombies[i].hit = true
			zombies[i].inHitAnimation = true
			zombies[i].hitTimer = hitFrameDuration
			zombies[i].hitFrame = 0
			zombies[i].hitAnimTimer = 0
			zombies[i].invulnerable = true
  		fmt.Println("Zombie", i, "hp:", zombies[i].hp)
		}

		if zombies[i].hit {
			zombies[i].invulnerable = true
		}
	}

    if hitFrameDuration > 0 {
		playerAttackActive = true
		hitFrameDuration--
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {  //called every frame, graphics.
	
	screen.DrawImage(background, nil)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(player1InitX,player1InitY)
	opSword := &ebiten.DrawImageOptions{}
	opSword.GeoM.Translate(swordX, swordY)
	screen.DrawImage(player1, op)	

	for i := range zombies {
		z := &zombies[i]

		op := &ebiten.DrawImageOptions{}
		w := float64(axeZombieSprites[z.walkFrame].Bounds().Dx())



		if z.hp <= 0 && !z.deathAnimationPlayed {
			op.GeoM.Translate(z.x, z.y)
			screen.DrawImage(axeZombieDeathSprites[z.deathAnimationFrame], op)
		} else if z.hp <= 0 {
			continue
		} else if z.inHitAnimation {
			op.GeoM.Translate(z.x, z.y)
			screen.DrawImage(axeZombieHitSprites[z.hitFrame], op)
		} else {
			if !z.facingRight {
				op.GeoM.Scale(-1, 1)
				op.GeoM.Translate(z.x + w, z.y)
			} else {
				op.GeoM.Translate(z.x, z.y)
			}
			
			screen.DrawImage(axeZombieSprites[z.walkFrame], op)
		}
	}

	if playerAttackFramesStart { // detects if player attack has started
  	if playerAttackFramesTimer == playerAttackFrames { // end of attack
    	playerAttackFramesTimer = 0
    	playerAttackFramesStart = false
    	playerAttackFlipped = (playerAttackCount % 2 == 0)
  	} else { // continue attack
			op := &ebiten.DrawImageOptions{}
  		frameImg := swordSprites[playerAttackFramesTimer]

  		w := float64(frameImg.Bounds().Dx()) // Dimensions
  		h := float64(frameImg.Bounds().Dy())
  		cx := w / 2
  		cy := h / 2

			var angle float64 // Determine angle (base sprite faces RIGHT)
  		
			switch swordLocation {
    		case 'd': // right
      		angle = 0
    		case 'a': // left
      		angle = math.Pi
    		case 's': // down
      		angle = math.Pi / 2
    		case 'w': // up
      		angle = -math.Pi / 2
    	}
    	
			scaleX := 1.0 // Apply vertical flipping if attack count requires
    	scaleY := 1.0
    		
			if playerAttackFlipped {
      	scaleY = -1.0
    	}
			
    	op.GeoM.Translate(-cx, -cy) // pivot center
			
    	op.GeoM.Scale(scaleX, scaleY) //scale, verticle
			
    	op.GeoM.Rotate(angle) //Rotate

    	op.GeoM.Translate(swordX+cx, swordY+cy) // Move final position (centered)

    	screen.DrawImage(frameImg, op)
			
    	playerAttackFramesTimer++
    	playerAttackCount++
  	}
	} else { // idle sword frame
    op := &ebiten.DrawImageOptions{}
    frameImg := swordSprites[0]
 		
    w := float64(frameImg.Bounds().Dx()) // Dimensions for idle as well
    h := float64(frameImg.Bounds().Dy())
    cx := w / 2
    cy := h / 2

    scaleX := 1.0
    scaleY := 1.0
    if playerAttackFlipped {
      scaleY = -1.0
    }

    var angle float64 // Idle frame always faces whatever swordLocation was last set to

  	switch swordLocation {
    	case 'd':
      	angle = 0
    	case 'a':
      	angle = math.Pi
    	case 's':
      	angle = math.Pi / 2
    	case 'w':
      	angle = -math.Pi / 2
    }

    op.GeoM.Translate(-cx, -cy)
    op.GeoM.Scale(scaleX, scaleY)
    op.GeoM.Rotate(angle)
    op.GeoM.Translate(swordX+cx, swordY+cy)

    screen.DrawImage(frameImg, op)
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
