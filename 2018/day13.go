package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"sort"
)

var inputFile = flag.String("inputFile", "inputs/day13.input", "Relative file path to use as input.")
var partB = flag.Bool("partB", false, "Whether to use the Part B logic.")

type Coord struct {
	X, Y int
}

type Direction int // 0 - Up, 1 = Right, 2 = Down, 3 = Left
type Rail int      // 0: |, 1: -, 2: \, 3: /, 4: +

type Cart struct {
	Loc       Coord
	NextTurn  int // 0 = left, 1 = straight, 2 = right
	Traveling Direction
}

type ByCoord []*Cart

func (c ByCoord) Len() int      { return len(c) }
func (c ByCoord) Swap(i, j int) { c[i], c[j] = c[j], c[i] }
func (c ByCoord) Less(i, j int) bool {
	if c[i].Loc.Y != c[j].Loc.Y {
		return c[i].Loc.Y < c[j].Loc.Y
	}
	return c[i].Loc.X < c[j].Loc.X
}

type Board map[Coord]Rail

func main() {
	flag.Parse()
	f, err := os.Open(*inputFile)
	if err != nil {
		return
	}
	defer f.Close()

	carts := make([]*Cart, 0)
	board := make(Board)

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
			case ' ':
				// Nothing
			case '|':
				// Rail enums:
				// 0: |, 1: -, 2: \, 3: /, 4: +
				board[loc] = 0
			case '-':
				board[loc] = 1
			case '\\':
				board[loc] = 2
			case '/':
				board[loc] = 3
			case '+':
				board[loc] = 4
			case '^':
				// Cart travel direction enums:
				// 0 - Up, 1 = Right, 2 = Down, 3 = Left
				board[loc] = 0
				carts = append(carts, &Cart{loc, 0, 0})
			case '>':
				board[loc] = 1
				carts = append(carts, &Cart{loc, 0, 1})
			case 'v':
				board[loc] = 0
				carts = append(carts, &Cart{loc, 0, 2})
			case '<':
				board[loc] = 1
				carts = append(carts, &Cart{loc, 0, 3})
			default:
				fmt.Printf("Failed to parse case %v\n", c)
				return
			}
		}
	}

	for tick := 0; len(carts) != 1; tick++ {
		sort.Sort(ByCoord(carts))
		seenCoord := make(map[Coord]*Cart)
		toRemove := make(map[*Cart]bool)

		for _, c := range carts {
			seenCoord[c.Loc] = c
		}
		for _, c := range carts {
			delete(seenCoord, c.Loc)
			nextCoord := c.Loc
			switch c.Traveling {
			case 0: // up
				nextCoord.Y -= 1
			case 1: // right
				nextCoord.X += 1
			case 2: // down
				nextCoord.Y += 1
			case 3: // left
				nextCoord.X -= 1
			}
			if seenCoord[nextCoord] != nil {
				if !*partB {
					fmt.Printf("Crash at %d,%d\n", nextCoord.X, nextCoord.Y)
					return
				} else {
					toRemove[seenCoord[nextCoord]] = true
					toRemove[c] = true
				}
			}
			if _, ok := board[nextCoord]; !ok {
				fmt.Printf("Previous position %d,%d\n", c.Loc.X, c.Loc.Y)
				fmt.Printf("illegal move to %d,%d\n", nextCoord.X, nextCoord.Y)
				return
			}
			switch board[nextCoord] {
			// Board states: 2: \, 3: /, 4: +
			case 2: // \
				// up -> left    0 to 3
				// right -> down 1 to 2
				// down -> right 2 to 1
				// left -> up    3 to 0
				c.Traveling = 3 - c.Traveling
			case 3: // /
				// up -> right   0 to 1
				// right -> up   1 to 0
				// down -> left  2 to 3
				// left -> down  3 to 2
				c.Traveling = (5 - c.Traveling) % 4
			case 4: // +
				// NextTurn possibilities:  0 = left, 1 = straight, 2 = right
				switch c.NextTurn {
				case 0: // CCW
					// decrement direction and mod
					c.Traveling = (c.Traveling + 3) % 4
				case 1: // no change
					// Deliberately empty
				case 2: // CW
					// Increment and mod direction
					c.Traveling = (c.Traveling + 1) % 4
				}
				c.NextTurn = (c.NextTurn + 1) % 3
			}
			seenCoord[nextCoord] = c
			c.Loc = nextCoord
		}
		newCarts := make([]*Cart, 0)
		for _, c := range carts {
			if !toRemove[c] {
				newCarts = append(newCarts, c)
			}
		}
		carts = newCarts
	}

	fmt.Printf("Final cart alive is at %d,%d\n", carts[0].Loc.X, carts[0].Loc.Y)
}
