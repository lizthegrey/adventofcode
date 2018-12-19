package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
)

var inputFile = flag.String("inputFile", "inputs/day18.input", "Relative file path to use as input.")
var rounds = flag.Int("rounds", 10, "The number of rounds to simulate.")

type Coord struct {
	X, Y int
}

func (c Coord) Adjacent(t [50][50]int) (int, int) {
	trees, yards := 0, 0

	for xoff := -1; xoff <= 1; xoff++ {
		for yoff := -1; yoff <= 1; yoff++ {
			if xoff == 0 && yoff == 0 {
				continue
			}
			y := c.Y + yoff
			x := c.X + xoff
			if y >= 50 || y < 0 || x >= 50 || x < 0 {
				continue
			}
			if t[y][x] == 1 {
				trees++
			} else if t[y][x] == 2 {
				yards++
			}
		}
	}
	return trees, yards
}

func main() {
	flag.Parse()
	f, err := os.Open(*inputFile)
	if err != nil {
		return
	}
	defer f.Close()

	// 0 == empty
	// 1 == lumberyard
	// 2 = trees
	var tiles [50][50]int

	reader := bufio.NewReader(f)
	for y := 0; ; y++ {
		l, err := reader.ReadString('\n')
		if err != nil || len(l) == 0 {
			break
		}

		l = l[:len(l)-1]
		for x, c := range l {
			switch c {
			case '.':
				tiles[y][x] = 0
			case '|':
				tiles[y][x] = 1
			case '#':
				tiles[y][x] = 2
			}
		}
	}

	// Look for cycles.
	seen := make(map[[50][50]int]int)

	loopLen := 0
	loopStart := 0

	for r := 0; r < *rounds; r++ {
		var newTiles [50][50]int

		for y := 0; y < 50; y++ {
			for x := 0; x < 50; x++ {
				trees, yards := Coord{x, y}.Adjacent(tiles)

				switch tiles[y][x] {
				case 0:
					if trees >= 3 {
						newTiles[y][x] = 1
					} else {
						newTiles[y][x] = 0
					}
				case 1:
					if yards >= 3 {
						newTiles[y][x] = 2
					} else {
						newTiles[y][x] = 1
					}
				case 2:
					if trees >= 1 && yards >= 1 {
						newTiles[y][x] = 2
					} else {
						newTiles[y][x] = 0
					}
				}
			}
		}
		if seen[newTiles] != 0 {
			loopStart = seen[newTiles]
			loopLen = r + 1 - loopStart
			fmt.Printf("Loop detected: %d to %d.\n", loopStart, r+1)
			break
		}
		seen[newTiles] = r + 1
		tiles = newTiles
	}

	if loopStart != 0 && loopLen != 0 {
		state := loopStart + ((*rounds - loopStart) % loopLen)
		fmt.Printf("Looking for state %d\n", state)
		for k, v := range seen {
			if v == state {
				tiles = k
				break
			}
		}
	}

	trees, yards := 0, 0
	for y := 0; y < 50; y++ {
		for x := 0; x < 50; x++ {
			if tiles[y][x] == 1 {
				trees++
			} else if tiles[y][x] == 2 {
				yards++
			}
		}
	}
	fmt.Printf("Found %d yards and %d trees for a result of %d.\n", yards, trees, yards*trees)
}
