package main

import (
	"log"
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"math"
)

var ( //declare variable for images, name *ebiten.Image.
	background *ebiten.Image
	player1 *ebiten.Image
	swordSprites []*ebiten.Image
	axeZombieDeathSprites []*ebiten.Image
	portalSprite *ebiten.Image

	axeZombieSprites []*ebiten.Image //an array of image files means it for a animation
	axeZombieHitSprites []*ebiten.Image	//see functions.go

	screenHeight = 540 //* 1.5 //= 810
	screenWidth = 960 //* 1.5 //= 1440

	//lower is faster
	axeZombieAnimationSpeed = float64(10)
	axeZombieHitAnimationSpeed = float64(5)
	//higher is faster
	axeZombieLiteralSpeed = float64(0.7)

	tickCount = 0 //for game time keeping
	
	sum int = 0
	floor [12][12]int
	floorInit bool = false

	zombies []axeZombie
	
	// Room system variables
	currentRoomX int = 6
	currentRoomY int = 6
	roomCleared [12][12]bool
	roomLocked bool = false
)

type Game struct{}

type Camera struct {
	x float64
	y float64
	following bool
}

var cam = Camera{
	x: 0,
	y: 0,
	following: true,
}

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
	x: 255, //init player location
	y: 132,	//player position

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

type portal struct {
	x, y float64
	targetRoomX, targetRoomY int
	direction string // "up", "down", "left", "right"
}

var portals []portal


func init() { //initialize images to variables here.
	var err error
	
 	background, _, err = ebitenutil.NewImageFromFile("assets/sprites/roomAssets/onexOneRoom.png") //name, _, etc.
	if err != nil {
		log.Fatal(err)
	}

	player1, _, err = ebitenutil.NewImageFromFile("assets/images/Sprite-0001.png") //will not run if empty
	if err != nil {
		log.Fatal(err)
	}
	
	portalSprite, _, err = ebitenutil.NewImageFromFile("assets/images/lightSaber.png")
	if err != nil {
		log.Fatal(err)
	}
	
	loadAxeZombieSprites() //call animation functions here
	loadAxeZombieHitSprites()
	loadSwordSprites()
	loadAxeZombieDeathSprites()

	// Don't spawn zombies at init - wait for room entry
}

func (g *Game) Update() error { //game logic

	if floorInit == false { //floor initialize
		
		floor[6][6] = 1

		for i := 0; i < len(floor); i++ {
	  	for j := 0; j < len(floor[i]); j++ {
	    	sum += floor[i][j]
			}
		}   
		
		if sum == 1 {
			
			initFloor(&floor, 20)

			for i := range floor {
				for j := range floor[i] {
					fmt.Printf("%2d ", floor[i][j])
				}
				fmt.Println()
			}
			
			// Initialize the starting room (no enemies in starting room)
			enterRoom(currentRoomX, currentRoomY, &floor, &roomCleared)
		}

		floorInit = true
	}

	tickCount++

	if tickCount == 2 {
		cam.following = false
	}

	// Toggle camera following with C key
	if inpututil.IsKeyJustPressed(ebiten.KeyC) {
		cam.following = !cam.following
		fmt.Println("Camera following:", cam.following)
	}

	if tickCount % 60 == 0 { //prints every 60 frames for time keeping.
		fmt.Println("frame", tickCount, ",", "RAM: ", GetSelfRAM(), "MB")
		for i := range zombies {
			fmt.Println("axe zombie" ,i ," frame: ", zombies[i].walkFrame)
		}
		fmt.Println("player x: ", p.x)
		fmt.Println("player y:  ", p.y)
	}

	if p.hitFrameDuration == 0 { // prevents player from attacking same enemy.
		for i := range zombies {
			if zombies[i].hitTimer == 0 {
				zombies[i].invulnerable = false
				p.attackActive = false
			}
		}
	}


	//sword logic

	switch { //player sword controls - arrow keys or Xbox ABXY
		case (ebiten.IsKeyPressed(ebiten.KeyArrowRight) || ebiten.IsStandardGamepadButtonPressed(0, ebiten.StandardGamepadButtonRightRight)) && p.hitFrameDuration == 0:
			p.swordLocation = 'd'
			p.hitFrameDuration = p.attackFrames
			p.attackFramesStart = true
		case (ebiten.IsKeyPressed(ebiten.KeyArrowLeft) || ebiten.IsStandardGamepadButtonPressed(0, ebiten.StandardGamepadButtonRightLeft)) && p.hitFrameDuration == 0:
			p.swordLocation = 'a'
			p.hitFrameDuration = p.attackFrames
			p.attackFramesStart = true
		case (ebiten.IsKeyPressed(ebiten.KeyArrowDown) || ebiten.IsStandardGamepadButtonPressed(0, ebiten.StandardGamepadButtonRightBottom)) && p.hitFrameDuration == 0:
			p.swordLocation = 's'
			p.hitFrameDuration = p.attackFrames
			p.attackFramesStart = true
		case (ebiten.IsKeyPressed(ebiten.KeyArrowUp) || ebiten.IsStandardGamepadButtonPressed(0, ebiten.StandardGamepadButtonRightTop)) && p.hitFrameDuration == 0:
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
	
	// Get analog stick input
	axisX := ebiten.StandardGamepadAxisValue(0, ebiten.StandardGamepadAxisLeftStickHorizontal)
	axisY := ebiten.StandardGamepadAxisValue(0, ebiten.StandardGamepadAxisLeftStickVertical)
	
	deadzone := 0.15 // Apply deadzone to prevent drift
	if math.Abs(axisX) < deadzone {
		axisX = 0
	}
	if math.Abs(axisY) < deadzone {
		axisY = 0
	}
	
	var moveX, moveY float64 // Calculate movement vector from all inputs
	
	if ebiten.IsKeyPressed(ebiten.KeyD) { // Keyboard input
		moveX += 1.0
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		moveX -= 1.0
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		moveY += 1.0
	}
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		moveY -= 1.0
	}
	
	if ebiten.IsStandardGamepadButtonPressed(0, ebiten.StandardGamepadButtonLeftRight) { // d-Pad input
		moveX += 1.0
	}
	if ebiten.IsStandardGamepadButtonPressed(0, ebiten.StandardGamepadButtonLeftLeft) {
		moveX -= 1.0
	}
	if ebiten.IsStandardGamepadButtonPressed(0, ebiten.StandardGamepadButtonLeftBottom) {
		moveY += 1.0
	}
	if ebiten.IsStandardGamepadButtonPressed(0, ebiten.StandardGamepadButtonLeftTop) {
		moveY -= 1.0
	}
	
	if axisX != 0 || axisY != 0 { //analog stick
		moveX = axisX
		moveY = axisY
	}
	
	if moveX != 0 && moveY != 0 { // Normalize diagonal movement to prevent faster diagonal speed
		magnitude := math.Sqrt(moveX*moveX + moveY*moveY)
		moveX /= magnitude
		moveY /= magnitude
	}
	
	if moveX > 0 && !isBlocked(p.x-25, p.y, 1, 0, blockRange, zombies) { // Apply movement with collision detection
		p.x += moveSpeed * moveX
	} else if moveX < 0 && !isBlocked(p.x, p.y, -1, 0, blockRange, zombies) {
		p.x += moveSpeed * moveX
	}
	
	if moveY > 0 && !isBlocked(p.x, p.y, 0, 1, blockRange, zombies) {
		p.y += moveSpeed * moveY
	} else if moveY < 0 && !isBlocked(p.x, p.y, 0, -1, blockRange, zombies) {
		p.y += moveSpeed * moveY
	}

	if cam.following { // Update camera to follow player
		playerWidth := float64(player1.Bounds().Dx()) // Get player sprite dimensions for proper centering
		playerHeight := float64(player1.Bounds().Dy())
		
		cam.x = p.x + playerWidth / 2 - float64(screenWidth) / 3
		cam.y = p.y + playerHeight / 2 - float64(screenHeight) / 2.35
	}

	// Check for portal collision and room transition
	portal := checkPortalCollision(p.x, p.y)
	if portal != nil {
		// Transition to new room
		currentRoomX = portal.targetRoomX
		currentRoomY = portal.targetRoomY
		
		// Reset player position based on which portal they entered
		// Room is 594x351 after scaling
		// Offset player away from the portal to prevent immediate re-trigger
		switch portal.direction {
		case "up":
			p.y = 280.0 // Enter from bottom, offset up from portal
			p.x = 297.0
		case "down":
			p.y = 70.0 // Enter from top, offset down from portal
			p.x = 297.0
		case "left":
			p.x = 520.0 // Enter from right, offset left from portal
			p.y = 175.5
		case "right":
			p.x = 74.0 // Enter from left, offset right from portal
			p.y = 175.5
		}
		
		enterRoom(currentRoomX, currentRoomY, &floor, &roomCleared)
	}
	
	// Check if current room is cleared
	checkRoomCleared(currentRoomX, currentRoomY, &roomCleared)

	zombieLogic()

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {  //called every frame, graphics
	
	// Draw background with camera offset and scaling
	opBg := &ebiten.DrawImageOptions{}
	
	// Calculate scale factors to fit the background to screen
	
	// Apply scaling first, then translate for camera
	opBg.GeoM.Scale(0.65, 0.65)
	opBg.GeoM.Translate(-cam.x, -cam.y)
	screen.DrawImage(background, opBg)

	// Draw player with camera offset
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(p.x - cam.x, p.y - cam.y)
	screen.DrawImage(player1, op)	

	// Draw zombies with camera offset
	for i := range zombies {
		z := &zombies[i]

		op := &ebiten.DrawImageOptions{}
		w := float64(axeZombieSprites[z.walkFrame].Bounds().Dx())

		if z.hp <= 0 && !z.deathAnimationPlayed {
			op.GeoM.Translate(z.x - cam.x, z.y - cam.y)
			screen.DrawImage(axeZombieDeathSprites[z.deathAnimationFrame], op)
		} else if z.hp <= 0 {
			continue
		} else if z.inHitAnimation {
			op.GeoM.Translate(z.x - cam.x, z.y - cam.y)
			screen.DrawImage(axeZombieHitSprites[z.hitFrame], op)
		} else {
			if !z.facingRight {
				op.GeoM.Scale(-1, 1)
				op.GeoM.Translate(z.x + w - cam.x, z.y - cam.y)
			} else {
				op.GeoM.Translate(z.x - cam.x, z.y - cam.y)
			}
			
			screen.DrawImage(axeZombieSprites[z.walkFrame], op)
		}
	}

	// Draw portals (only if room is not locked)
	if !roomLocked {
		for i := range portals {
			op := &ebiten.DrawImageOptions{}
			
			// lightSaber is 400x700, scale down to reasonable size (about 40x70 pixels)
			op.GeoM.Scale(0.1, 0.1)
			
			// Center the portal on its position
			pw := float64(portalSprite.Bounds().Dx()) * 0.1  // 400 * 0.1 = 40px
			ph := float64(portalSprite.Bounds().Dy()) * 0.1  // 700 * 0.1 = 70px
			op.GeoM.Translate(portals[i].x - pw/2 - cam.x, portals[i].y - ph/2 - cam.y)
			
			screen.DrawImage(portalSprite, op)
		}
	}

	// Draw sword with camera offset
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
			
			scaleX := 2.0 // Apply vertical flipping if attack count requires
			scaleY := 2.0
			
			if p.attackFlipped {
				scaleY = -2.0
			}
			
			op.GeoM.Translate(-cx, -cy) // pivot center
			
			op.GeoM.Scale(scaleX, scaleY) //scale, verticle
			
			op.GeoM.Rotate(angle) //Rotate

			op.GeoM.Translate(p.swordX + cx - cam.x, p.swordY + cy - cam.y) // Move final position (centered) with camera offset

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

		scaleX := 2.0
		scaleY := 2.0
		if p.attackFlipped {
			scaleY = -2.0
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
		op.GeoM.Translate(p.swordX + cx - cam.x, p.swordY + cy - cam.y) // with camera offset

		screen.DrawImage(frameImg, op)
	}

	// Draw camera status
	debugText := fmt.Sprintf("Camera Follow: %v (Toggle with C)\nRoom: [%d][%d] | Locked: %v", 
		cam.following, currentRoomX, currentRoomY, roomLocked)
	ebitenutil.DebugPrint(screen, debugText)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 640, 480
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Chaos Crashers")
	ebiten.SetFullscreen(true)
	
	if err := ebiten.RunGame(&Game{}); err != nil { 
		log.Fatal(err)
	}	
}
