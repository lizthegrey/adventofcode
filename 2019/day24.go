package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strings"
)

var inputFile = flag.String("inputFile", "inputs/day24.input", "Relative file path to use as input.")
var iterations = flag.Int("iterations", 200, "Number of iterations to run for part B.")

type Coord struct {
	Y, X int
}
type Grid [5][5]bool
type Grid3 map[int]Grid

func main() {
	flag.Parse()
	bytes, err := ioutil.ReadFile(*inputFile)
	if err != nil {
		return
	}
	contents := string(bytes)
	split := strings.Split(contents, "\n")

	var g Grid
	for y, s := range split {
		if s == "" {
			continue
		}
		for x, c := range s {
			switch c {
			case '#':
				g[y][x] = true
			case '.':
				// Do nothing, it's an empty space.
			default:
				fmt.Printf("Unable to parse line %s\n", s)
			}
		}
	}

	// Part A
	grid := g
	states := map[Grid]bool{
		grid: true,
	}
	for {
		grid = grid.Evolve()
		if states[grid] == true {
			break
		}
		states[grid] = true
	}
	fmt.Println(g.BioDiversity())

	// Part B
	g3 := Grid3{
		0: g,
	}
	for i := 1; i <= *iterations; i++ {
		g3 = g3.Evolve(-i, i)
	}
	alive := 0
	for _, v := range g3 {
		alive += v.Count()
	}
	fmt.Println(alive)
}

func (g Grid) BioDiversity() int {
	ret := 0
	power := 0
	for _, r := range g {
		for _, c := range r {
			if c {
				ret += 1 << power
			}
			power++
		}
	}
	return ret
}

func (g Grid) Count() int {
	ret := 0
	for _, r := range g {
		for _, c := range r {
			if c {
				ret++
			}
		}
	}
	return ret
}

func (g Grid) Evolve() Grid {
	var nextGen Grid
	for y, r := range g {
		for x, c := range r {
			if c {
				if g.NeighborsAlive(y, x) == 1 {
					nextGen[y][x] = true
				}
			} else {
				switch g.NeighborsAlive(y, x) {
				case 1:
					fallthrough
				case 2:
					nextGen[y][x] = true
				}
			}
		}
	}
	return nextGen
}

func (g Grid) NeighborsAlive(y, x int) int {
	sum := 0
	if g.BugAlive(y+1, x) {
		sum++
	}
	if g.BugAlive(y-1, x) {
		sum++
	}
	if g.BugAlive(y, x+1) {
		sum++
	}
	if g.BugAlive(y, x-1) {
		sum++
	}
	return sum
}

func (g Grid) BugAlive(y, x int) bool {
	if y < 0 || y >= 5 || x < 0 || x >= 5 {
		return false
	}
	return g[y][x]
}

func (g3 Grid3) Evolve(minZ, maxZ int) Grid3 {
	nextGen := make(Grid3)
	// Need to instantiate levels that don't exist yet, one above and below current Z.
	for z := minZ; z <= maxZ; z++ {
		next := Grid{}
		for y, r := range g3[z] {
			for x, c := range r {
				if x == 2 && y == 2 {
					continue
				}
				if c {
					if g3.NeighborsAlive(z, y, x) == 1 {
						next[y][x] = true
					}
				} else {
					switch g3.NeighborsAlive(z, y, x) {
					case 1:
						fallthrough
					case 2:
						next[y][x] = true
					}
				}
			}
		}
		nextGen[z] = next
	}
	return nextGen
}

func (g3 Grid3) NeighborsAlive(z, y, x int) int {
	sum := 0
	start := Coord{y, x}
	sum += g3.BugsAlive(start, z, y+1, x)
	sum += g3.BugsAlive(start, z, y-1, x)
	sum += g3.BugsAlive(start, z, y, x+1)
	sum += g3.BugsAlive(start, z, y, x-1)
	return sum
}

func (g3 Grid3) BugsAlive(start Coord, z, y, x int) int {
	sum := 0
	if y == 2 && x == 2 {
		inner := g3[z+1]
		if start.Y < 2 {
			// From the top
			for x := 0; x < 5; x++ {
				if inner[0][x] {
					sum++
				}
			}
		} else if start.Y > 2 {
			// From the bottom
			for x := 0; x < 5; x++ {
				if inner[4][x] {
					sum++
				}
			}
		} else if start.X < 2 {
			// From the left
			for y := 0; y < 5; y++ {
				if inner[y][0] {
					sum++
				}
			}
		} else if start.X > 2 {
			// From the right
			for y := 0; y < 5; y++ {
				if inner[y][4] {
					sum++
				}
			}
		} else {
			fmt.Println("Asked to compute illegal value.")
		}
		return sum
	}
	if y < 0 || y >= 5 || x < 0 || x >= 5 {
		outer := g3[z-1]
		var present bool
		if y < 0 {
			// From the top
			present = outer[1][2]
		} else if y >= 5 {
			// From the bottom
			present = outer[3][2]
		} else if x < 0 {
			// From the left
			present = outer[2][1]
		} else if x >= 5 {
			// From the right
			present = outer[2][3]
		} else {
			fmt.Println("Asked to compute illegal value.")
		}
		if present {
			return 1
		}
		return 0
	}
	if g3[z][y][x] {
		return 1
	}
	return 0
}
