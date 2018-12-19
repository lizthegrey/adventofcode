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

func (c Coord) Adjacent() []Coord {
	r := make([]Coord, 0)
	for x := -1; x <= 1; x++ {
		for y := -1; y <= 1; y++ {
			if x == 0 && y == 0 {
				continue
			}
			r = append(r, Coord{c.X + x, c.Y + y})
		}
	}
	return r
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
	tiles := make(map[Coord]int)

	r := bufio.NewReader(f)
	for y := 0; ; y++ {
		l, err := r.ReadString('\n')
		if err != nil || len(l) == 0 {
			break
		}
		l = l[:len(l)-1]
		for x, c := range l {
			loc := Coord{x, y}
			switch c {
			case '.':
				tiles[loc] = 0
			case '#':
				tiles[loc] = 1
			case '|':
				tiles[loc] = 2
			}
		}
	}

	for r := 0; r < *rounds; r++ {
		newTiles := make(map[Coord]int)

		for c, v := range tiles {
			adj := c.Adjacent()

			trees := 0
			yards := 0
			for _, n := range adj {
				if tiles[n] == 1 {
					yards++
				} else if tiles[n] == 2 {
					trees++
				}
			}

			switch v {
			case 0:
				if trees >= 3 {
					newTiles[c] = 2
				} else {
					newTiles[c] = 0
				}
			case 1:
				if trees >= 1 && yards >= 1 {
					newTiles[c] = 1
				} else {
					newTiles[c] = 0
				}
			case 2:
				if yards >= 3 {
					newTiles[c] = 1
				} else {
					newTiles[c] = 2
				}
			}
		}
		tiles = newTiles
	}

	trees := 0
	yards := 0
	for _, v := range tiles {
		if v == 1 {
			yards++
		} else if v == 2 {
			trees++
		}
	}
	fmt.Printf("Found %d yards and %d trees for a result of %d.\n", yards, trees, yards*trees)
}
