package main

func initFloor(floor *[12][12]int, targetCount int) {
	if floor == nil || targetCount <= 0 {
		return
	}

	maxI := len(floor)
	maxJ := len(floor[0])
	
	if maxI == 0 || maxJ == 0 {
		return
	}

	initialized := make(map[[2]int]bool)
	currentCount := 0

	centerI, centerJ := 6, 6 // Initialize the center point to 1
	if centerI < maxI && centerJ < maxJ {
		floor[centerI][centerJ] = 1
		initialized[[2]int{centerI, centerJ}] = true
		currentCount++
	

	for i := 0; i < maxI; i++ { // Count and track any other already non-zero elements
		for j := 0; j < maxJ; j++ {
			if floor[i][j] != 0 && !initialized[[2]int{i, j}] {
				initialized[[2]int{i, j}] = true
				currentCount++
			}
		}
	}

	queue := make([][2]int, 0) // Keep a queue of cells to process for more even distribution
	queue = append(queue, [2]int{centerI, centerJ})
	processed := make(map[[2]int]bool)

	for currentCount < targetCount && len(queue) > 0 { // Process elements until we reach target count
		current := queue[0] // Take from front of queue for breadth-first expansion
		queue = queue[1:]
		
		i, j := current[0], current[1]
		
		if processed[current] { // Skip if already processed
			continue
		}
		processed[current] = true

		neighbors := getNeighborsRandomized(i, j, maxI, maxJ, initialized) // Get all valid neighbors in random order for even distribution
		
		if len(neighbors) == 0 {
			continue
		}

		firstNeighbor := neighbors[0] // Always initialize one neighbor 100% - 50% - 25% - 12.50%
		if !initialized[firstNeighbor] {
			floor[firstNeighbor[0]][firstNeighbor[1]] = randInt(1, 2)
			initialized[firstNeighbor] = true
			currentCount++
			queue = append(queue, firstNeighbor)

			if currentCount >= targetCount {
				return
			}

			neighbors = getNeighborsRandomized(i, j, maxI, maxJ, initialized) // Update neighbors list after first initialization
		}

		// Progressively fill remaining neighbors with decreasing probability
		probability := 50.0 // Start at 50% for second neighbor
		for _, neighbor := range neighbors {
			if currentCount >= targetCount {
				return
			}
			
			if !initialized[neighbor] && randInt(1, 100) <= int(probability) {
				floor[neighbor[0]][neighbor[1]] = randInt(1, 2)
				initialized[neighbor] = true
				currentCount++
				queue = append(queue, neighbor)
			}
			probability /= 2.0 // Halve the probability for each subsequent neighbor
		}
	}
}

func getNeighborsRandomized(i, j, maxI, maxJ int, initialized map[[2]int]bool) [][2]int { // getNeighborsRandomized returns all valid, uninitialized neighbors in random order
	neighbors := make([][2]int, 0, 4)
	
	// Up
	if j > 0 && !initialized[[2]int{i, j - 1}] {
		neighbors = append(neighbors, [2]int{i, j - 1})
	}
	// Down
	if j < maxJ-1 && !initialized[[2]int{i, j + 1}] {
		neighbors = append(neighbors, [2]int{i, j + 1})
	}
	// Left
	if i > 0 && !initialized[[2]int{i - 1, j}] {
		neighbors = append(neighbors, [2]int{i - 1, j})
	}
	// Right
	if i < maxI-1 && !initialized[[2]int{i + 1, j}] {
		neighbors = append(neighbors, [2]int{i + 1, j})
	}
	
	shuffleNeighbors(neighbors) // Shuffle neighbors to randomize direction preference
	
	return neighbors
}

func shuffleNeighbors(neighbors [][2]int) { // shuffleNeighbors randomizes the order of neighbors (Fisher-Yates shuffle)
	for i := len(neighbors) - 1; i > 0; i-- {
		j := randInt(0, i)
		neighbors[i], neighbors[j] = neighbors[j], neighbors[i]
	}
}
