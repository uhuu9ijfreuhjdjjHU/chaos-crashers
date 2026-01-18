package main

import "fmt"

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
		swordHitRange := 80.0

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
