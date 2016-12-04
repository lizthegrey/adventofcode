package main

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

type coord struct {
	X, Y int
}

func main() {
	input := strings.Split("R5, R4, R2, L3, R1, R1, L4, L5, R3, L1, L1, R4, L2, R1, R4, R4, L2, L2, R4, L4, R1, R3, L3, L1, L2, R1, R5, L5, L1, L1, R3, R5, L1, R4, L5, R5, R1, L185, R4, L1, R51, R3, L2, R78, R1, L4, R188, R1, L5, R5, R2, R3, L5, R3, R4, L1, R2, R2, L4, L4, L5, R5, R4, L4, R2, L5, R2, L1, L4, R4, L4, R2, L3, L4, R2, L3, R3, R2, L2, L3, R4, R3, R1, L4, L2, L5, R4, R4, L1, R1, L5, L1, R3, R1, L2, R1, R1, R3, L4, L1, L3, R2, R4, R2, L2, R1, L5, R3, L3, R3, L1, R4, L3, L3, R4, L2, L1, L3, R2, R3, L2, L1, R4, L3, L5, L2, L4, R1, L4, L4, R3, R5, L4, L1, L1, R4, L2, R5, R1, R1, R2, R1, R5, L1, L3, L5, R2", ", ")
	dir := 0
	loc := coord{0, 0}
	var repeated *coord
	visited := make(map[coord]bool)
	visited[loc] = true
	for _, inst := range input {
		turn := string(inst[0])
		distance, err := strconv.Atoi(inst[1:])
		if err != nil {
			fmt.Printf("Invalid distance: %d\n", inst[1:])
			return
		}
		if turn == "L" {
			dir = (dir - 1) % 4
		} else if turn == "R" {
			dir = (dir + 1) % 4
		} else {
			fmt.Printf("Invalid turn: %s\n", turn)
			return
		}
		if dir < 0 {
			dir = dir + 4
		}
		for i := 0; i < distance; i++ {
			switch dir {
			case 0:
				loc.Y += 1
			case 1:
				loc.X += 1
			case 2:
				loc.Y -= 1
			case 3:
				loc.X -= 1
			default:
				fmt.Printf("Invalid direction: %d", dir)
				return
			}
			if visited[loc] && repeated == nil {
				saved := loc
				repeated = &saved
			}
			visited[loc] = true
		}
	}
	fmt.Println(math.Abs(float64(loc.X)) + math.Abs(float64(loc.Y)))
	fmt.Println(math.Abs(float64(repeated.X)) + math.Abs(float64(repeated.Y)))
}
