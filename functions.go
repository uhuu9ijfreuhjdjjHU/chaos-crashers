package main

import (
	"math/rand"
	"log"
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

func spawnZombies() {
  count := randInt(3, 6)

  for i := 0; i < count; i++ {
    z := axeZombie{
    x:     randFloat(0, float64(screenWidth + 100)),
    y:     randFloat(0, float64(screenHeight + 100)),
    hp:    randInt(3, 10),
    level: randInt(1, 3),
    speed: randFloat(0.3, 1.0),
    }
    
		zombies = append(zombies, z)
  }
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


func isBlocked(px, py float64, dx, dy float64, blockRange float64, zombies []axeZombie) bool {
  for _, z := range zombies {
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
