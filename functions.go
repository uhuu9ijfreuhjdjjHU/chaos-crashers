package main

import (
	"math/rand"
	"log"
	"fmt"
  "os"
  "bufio"
  "strings"
  "strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

func spawnAxeZombies(speedSelect float64) {
  count := randInt(7, 7)

  for i := 0; i < count; i++ {
    z := axeZombie{
    x:     randFloat(0, float64(screenWidth + 100)),
    y:     randFloat(0, float64(screenHeight + 100)),
    hp:    randInt(3, 10),
    level: randInt(1, 3),
    speed: speedSelect,
		facingRight: true,
		invulnerable: false,
		walkFrame: randInt(0, (len(axeZombieSprites) - 1)),
		hitFrame: 1,
		inHitAnimation: false,
		deathAnimationPlayed: false,
		deathAnimationTimer: 0,
		deathAnimationFrame: 0,
		knockbackSpeed: 5,
    }
    
		zombies = append(zombies, z)
  }
}

func zombieHitAnimationUpdate(animationSpeed float64) {
	for i := range zombies {
		z := &zombies[i]

		// Only update if the zombie is in hit animation
		if !z.inHitAnimation {
			continue
		}

		z.invulnerable = true

		z.hitAnimTimer++

		if z.hitAnimTimer >= animationSpeed {
			z.hitAnimTimer = 0
			z.hitFrame++

			if z.hitFrame >= len(axeZombieHitSprites) {
				// Hit animation finished
				z.hitFrame = 0
				z.inHitAnimation = false
				z.hit = false
				z.invulnerable = false
			}
		}
	}
}


func zombieDeathAnimationUpdate(animationSpeed float64) {
	for i := range zombies {
		z := &zombies[i]

		// Only animate dead zombies that haven't finished their animation
		if z.hp > 0 || z.deathAnimationPlayed {
			continue
		}

		z.deathAnimationTimer++

		if z.deathAnimationTimer >= animationSpeed {
			z.deathAnimationTimer = 0
			z.deathAnimationFrame++

			// When last frame reached, mark animation done
			if z.deathAnimationFrame >= len(axeZombieDeathSprites) {
				z.deathAnimationPlayed = true
				z.deathAnimationFrame = len(axeZombieDeathSprites) - 1 // freeze on last frame
			}
		}
	}
}


func zombieWalkCycleUpdate(animationSpeed float64) {
	for i := range zombies {
		z := &zombies[i]

		// increment timer
		z.walkTimer++
		
		if z.walkTimer >= animationSpeed {
			z.walkTimer = 0
			z.walkFrame++
			if z.walkFrame >= len(axeZombieSprites) {
				z.walkFrame = 0
			}
		}
	}
}


func GetSelfRAM() float64 {
  
	file, err := os.Open("/proc/self/status")

  if err != nil {
  	return -1
  }
  
	defer file.Close()

  scanner := bufio.NewScanner(file)
    for scanner.Scan() {
      line := scanner.Text()
      if strings.HasPrefix(line, "VmRSS:") {
      fields := strings.Fields(line)
      kb, _ := strconv.Atoi(fields[1]) // value in KB
      return float64(kb) / 1024        // return MB
    }
  }
  return -1
}

func enemyMovement(targetX, targetY, x, y, speed float64, knockBackSpeed float64, knockbackDirection rune, zombies []axeZombie, self int) (float64, float64) {
  // --- Chase player ---
  dx := 0.0
  dy := 0.0
	
	switch {
		case zombies[self].hit == true && knockbackDirection == 'a':
			dx -= knockBackSpeed
		case zombies[self].hit == true && knockbackDirection == 'd':
			dx += knockBackSpeed
		case zombies[self].hit == true && knockbackDirection == 's':
			dy += knockBackSpeed
		case zombies[self].hit == true && knockbackDirection == 'w':
			dy -= knockBackSpeed
	}

	if !zombies[self].hit {
  	if x < targetX - 100 {
    	dx += speed
			zombies[self].facingRight = true
  	}
  	if x > targetX + 100 {
    	dx -= speed
			zombies[self].facingRight = false
  	}
 		if y + 20 < targetY - 100 {
  		dy += speed
  	}
  	if y - 20 > targetY + 100 {
   		dy -= speed
  	}
	}

    // --- Avoid other zombies ---\\
  avoidDist := 60.0

  for i, z := range zombies {
    if i == self {
      continue
    }

    diffX := x - z.x
    diffY := y - z.y

    if abs(diffX) < avoidDist && abs(diffY) < avoidDist { // push away from nearby zombie
			if diffX > 0 {
        dx += speed * 0.5
      } else {
        dx -= speed * 0.5
      }
      if diffY > 0 {
        dy += speed * 0.5
      } else {
        dy -= speed * 0.5
      }
    }
  }
 return x + dx, y + dy
}

func loadAxeZombieDeathSprites() {
  axeZombieDeathSprites = make([]*ebiten.Image, 11)

  for i := 1; i <= 11; i++ {
    filename := fmt.Sprintf("assets/sprites/enemies/axeZombie/axeZombieDeath/zombieDeath%d.png", i)

    img, _, err := ebitenutil.NewImageFromFile(filename)
    if err != nil {
    log.Fatal(err)
    }
    axeZombieDeathSprites[i-1] = img
	}
}

func loadSwordSprites() {
  swordSprites = make([]*ebiten.Image, 15)

  for i := 1; i <= 15; i++ {
    filename := fmt.Sprintf("assets/sprites/swordSwing/coin%d.png", i)

    img, _, err := ebitenutil.NewImageFromFile(filename)
    if err != nil {
    log.Fatal(err)
    }
    swordSprites[i-1] = img
	}
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

func loadAxeZombieHitSprites() {
  axeZombieHitSprites = make([]*ebiten.Image, 8)

  for i := 1; i <= 8; i++ {
    filename := fmt.Sprintf("assets/sprites/enemies/axeZombie/axeZombieHit/axeZombieHit%02d.png", i)

    img, _, err := ebitenutil.NewImageFromFile(filename)
    if err != nil {
    log.Fatal(err)
    }
    axeZombieHitSprites[i-1] = img
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


func isBlocked(px, py float64, dx, dy float64, blockRange float64, zombies []axeZombie) bool {
  for _, z := range zombies {
  	if z.hp == 0 {
			continue
		}

		// Project the check range in the direction the player wants to move
    checkX := px + dx*blockRange
    checkY := py + dy*blockRange

    // If an enemy is near that projected point → blocked
    if abs(z.x-checkX) < 50 && abs(z.y-checkY) < 50 {
      return true
    }
  }
 
	return false
}


func randInt(min, max int) int {
	return rand.Intn(max-min+1) + min
}

func zombieLogic() {

	zombieWalkCycleUpdate(axeZombieAnimationSpeed)
	zombieHitAnimationUpdate(axeZombieHitAnimationSpeed)	
	zombieDeathAnimationUpdate(3)

	for i := range zombies { //keeps track of how long zombies should be "hit" for
    if zombies[i].hitTimer > 0 {
      zombies[i].hitTimer--
      zombies[i].hit = true
    } else if zombies[i].hitTimer == 0 {
      zombies[i].hit = false
    }
	}

	for i := range zombies { //zombie ai / logic

		if zombies[i].hp <= 0 {
			continue
		}
		
  	// movement (once per zombie)
  	zombies[i].x, zombies[i].y = enemyMovement(
  		p.x,
    	p.y,
    	zombies[i].x,
    	zombies[i].y,
    	zombies[i].speed,
			zombies[i].knockbackSpeed,
			p.swordLocation,
    	zombies,
    	i,
  	)
		
  	//zombie attack
  	hitRange := 80.0
  	
		if abs(zombies[i].x - p.x) < hitRange && 
		abs(zombies[i].y - p.y) < hitRange && tickCount % 150 == 0 {	
    	p.hp--
    	fmt.Println("hp:", p.hp) 
  	}

		//player attack
		swordHitRange := 55.0

  	if abs(zombies[i].x - p.swordX) < swordHitRange && !zombies[i].invulnerable &&
		abs(zombies[i].y - p.swordY) < swordHitRange && p.attackActive && 
		zombies[i].hitTimer <= 0 {
			zombies[i].hp--
			zombies[i].hit = true
			zombies[i].inHitAnimation = true
			zombies[i].hitTimer = p.hitFrameDuration
			zombies[i].hitFrame = 0
			zombies[i].hitAnimTimer = 0
			zombies[i].invulnerable = true
  		fmt.Println("Zombie", i, "hp:", zombies[i].hp)
		}

		if zombies[i].hit {
			zombies[i].invulnerable = true
		}
	}

    if p.hitFrameDuration > 0 {
		p.attackActive = true
		p.hitFrameDuration--
	}
}
