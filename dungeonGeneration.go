package main

func initFloor(floor *[12][12]int) {
	maxI := len(floor)
	maxJ := len(floor[0])

	for i := range floor {
		for j := range floor[i] {
			if floor[i][j] != 0 {
				bin := randInt(1, 4)

				switch bin {
				case 1: // up
					if j > 0 {
						floor[i][j-1] = randInt(1, 10)
					}
				case 2: // down
					if j < maxJ-1 {
						floor[i][j+1] = randInt(1, 10)
					}
				case 3: // right
					if i < maxI-1 {
						floor[i+1][j] = randInt(1, 10)
					}
				case 4: // left
					if i > 0 {
						floor[i-1][j] = randInt(1, 10)
					}
				}
			}
		}
	}
}
