package main

import (
	"math/rand"
	"log"
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

func spawnAxeZombies(speedSelect float64) {
  count := randInt(3, 6)

  for i := 0; i < count; i++ {
    z := axeZombie{
    x:     randFloat(0, float64(screenWidth + 100)),
    y:     randFloat(0, float64(screenHeight + 100)),
    hp:    randInt(3, 10),
    level: randInt(1, 3),
    speed: speedSelect,
		facingRight: true,
		invulnerable: false,
    }
    
		zombies = append(zombies, z)
  }
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
  	if x < targetX - 80 {
    	dx += speed
			zombies[self].facingRight = true
  	}
  	if x > targetX + 80 {
    	dx -= speed
			zombies[self].facingRight = false
  	}
 		if y + 20 < targetY - 80 {
  		dy += speed
  	}
  	if y - 20 > targetY + 80 {
   		dy -= speed
  	}
	}

    // --- Avoid other zombies ---\\
  avoidDist := 40.0

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
