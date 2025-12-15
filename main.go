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

	screenHeight = 540
	screenWidth = 960

	//lower is faster
	axeZombieAnimationSpeed = float64(10)
	axeZombieHitAnimationSpeed = float64(5)
	//higher is faster
	axeZombieLiteralSpeed = float64(0.7)

	tickCount = 0 //for game time keeping

	zombies []axeZombie
)

type Game struct{}

type player struct {
	x float64
	y float64

	swordX float64
	swordY float64

	hp int

	hitFrameDuration int

	attackCount int
	attackFrames int
	attackFramesTimer int

	swordLocation rune // 'a','d','s','w'

	attackActive bool
	attackFlipped bool
	attackFramesStart bool
}


var p = player {
	x: 560,
	y: 240,

	swordX: 560,
	swordY: 240 + 100,

	hp: 20,

	hitFrameDuration: 0,

	attackCount: 0,
	attackFrames: 15,
	attackFramesTimer: 0,

	swordLocation: 's',

	attackActive: false,
	attackFlipped: false,
	attackFramesStart: false,
}

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
	knockbackSpeed	float64
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

	if tickCount % 60 == 0 { //prints every 60 frames for time keeping.
		fmt.Println("frame", tickCount, ",", "RAM: ", GetSelfRAM(), "MB")
		for i := range zombies {
			fmt.Println("axe zombie" ,i ," frame: ", zombies[i].walkFrame)
		}
	}

	if p.hitFrameDuration == 0 { // prevents player from attacking same enemy.
		for i := range zombies {
			if zombies[i].hitTimer == 0 {
				zombies[i].invulnerable = false
				p.attackActive = false
			}
		}
	}

	//~~> sword direction logic <~~\\

	switch { //player sword controls
		case ebiten.IsKeyPressed(ebiten.KeyArrowRight) && p.hitFrameDuration == 0:
			p.swordLocation = 'd'
			p.hitFrameDuration = p.attackFrames
			p.attackFramesStart = true
		case ebiten.IsKeyPressed(ebiten.KeyArrowLeft) && p.hitFrameDuration == 0:
			p.swordLocation = 'a'
			p.hitFrameDuration = p.attackFrames
			p.attackFramesStart = true
		case ebiten.IsKeyPressed(ebiten.KeyArrowDown) && p.hitFrameDuration == 0:
			p.swordLocation = 's'
			p.hitFrameDuration = p.attackFrames
			p.attackFramesStart = true
		case ebiten.IsKeyPressed(ebiten.KeyArrowUp) && p.hitFrameDuration == 0:
			p.swordLocation = 'w'
			p.hitFrameDuration = p.attackFrames
			p.attackFramesStart = true	
	}

	switch { //player sword direction logic, effected by player sword controls above
		case p.swordLocation == 'd':
			p.swordX = p.x + 100
			p.swordY = p.y
		case p.swordLocation == 'a':
			p.swordX = p.x - 100
			p.swordY = p.y
		case p.swordLocation == 's':
			p.swordX = p.x
			p.swordY = p.y + 100
		case p.swordLocation == 'w':
			p.swordX = p.x
			p.swordY = p.y - 100
	}
	
	moveSpeed := 3.0 //player move speed
	blockRange := 50.0 //player collusion stat

	//player movement

	// MOVE RIGHT (D)
	if ebiten.IsKeyPressed(ebiten.KeyD) &&
  !isBlocked(p.x - 25, p.y, 1, 0, blockRange, zombies) {
  	p.x += moveSpeed
	}

	// MOVE LEFT (A)
	if ebiten.IsKeyPressed(ebiten.KeyA) &&
  !isBlocked(p.x, p.y, -1, 0, blockRange, zombies) {
  	p.x -= moveSpeed
	}

	// DOWN (S)
	if ebiten.IsKeyPressed(ebiten.KeyS) &&
  !isBlocked(p.x, p.y, 0, 1, blockRange, zombies) { //even though going down should be -1
  	p.y += moveSpeed																//for deincrimaenting the vert position
	}																									//the code doesnt work that way.

	// UP (W)
	if ebiten.IsKeyPressed(ebiten.KeyW) &&
  !isBlocked(p.x, p.y, 0, -1, blockRange, zombies) {
		p.y -= moveSpeed
	}

	zombieLogic()

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {  //called every frame, graphics.
	
	screen.DrawImage(background, nil)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(p.x, p.y)
	opSword := &ebiten.DrawImageOptions{}
	opSword.GeoM.Translate(p.swordX, p.swordY)
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

	if p.attackFramesStart { // detects if player attack has started
  	if p.attackFramesTimer == p.attackFrames { // end of attack
    	p.attackFramesTimer = 0
    	p.attackFramesStart = false
    	p.attackFlipped = (p.attackCount % 2 == 0)
  	} else { // continue attack
			op := &ebiten.DrawImageOptions{}
  		frameImg := swordSprites[p.attackFramesTimer]

  		w := float64(frameImg.Bounds().Dx()) // Dimensions
  		h := float64(frameImg.Bounds().Dy())
  		cx := w / 2
  		cy := h / 2

			var angle float64 // Determine angle (base sprite faces RIGHT)
  		
			switch p.swordLocation {
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
    		
			if p.attackFlipped {
      	scaleY = -1.0
    	}
			
    	op.GeoM.Translate(-cx, -cy) // pivot center
			
    	op.GeoM.Scale(scaleX, scaleY) //scale, verticle
			
    	op.GeoM.Rotate(angle) //Rotate

    	op.GeoM.Translate(p.swordX + cx, p.swordY + cy) // Move final position (centered)

    	screen.DrawImage(frameImg, op)
			
    	p.attackFramesTimer++
    	p.attackCount++
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
    if p.attackFlipped {
      scaleY = -1.0
    }

    var angle float64 // Idle frame always faces whatever swordLocation was last set to

  	switch p.swordLocation {
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
    op.GeoM.Translate(p.swordX + cx, p.swordY + cy)

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
