package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

var inputFile = flag.String("inputFile", "inputs/day22.input", "Relative file path to use as input.")
var debug = flag.Bool("debug", false, "Whether to print debug output.")

type board struct {
	tiles      map[coord]bool
	maxR, maxC int
}

type dir int

const (
	right dir = iota
	down
	left
	up
)

var dirs = map[dir]string{
	right: "right",
	down:  "down",
	left:  "left",
	up:    "up",
}

type coord struct {
	r, c int
}

func (t turtle) String(p board) string {
	width := p.maxC / 3
	height := p.maxR / 4
	relC := t.c % width
	relR := t.r % height
	faces := map[coord]string{
		coord{0 * height, 1 * width}: "A",
		coord{0 * height, 2 * width}: "B",
		coord{1 * height, 1 * width}: "C",
		coord{2 * height, 0 * width}: "D",
		coord{2 * height, 1 * width}: "E",
		coord{3 * height, 0 * width}: "F",
	}
	return fmt.Sprintf("(%d,%d) on %s, facing %s", relR, relC, faces[coord{t.r - relR, t.c - relC}], dirs[t.facing])
}

func (t *turtle) step() {
	switch t.facing {
	case right:
		t.c += 1
	case down:
		t.r += 1
	case left:
		t.c -= 1
	case up:
		t.r -= 1
	}
}

func (t *turtle) enforceBounds(p board, cube bool) bool {
	if !cube {
		if t.r <= -1 {
			t.r = p.maxR
		} else if t.r >= p.maxR {
			t.r = -1
		} else if t.c <= -1 {
			t.c = p.maxC
		} else if t.c >= p.maxC {
			t.c = -1
		} else {
			return false
		}
	} else {
		// Look up in the table of how we should map to a new coordinate/direction.
		// For now, just hardcoding the table because, yeah. elly tells me it's very difficult
		// to generically deduce so I will special case to my folding.
		width := p.maxC / 3
		height := p.maxR / 4
		relC := t.c % width
		relR := t.r % height

		// My input (will not work on small test input):
		// 012cr
		//  AB 0
		//  C  1
		// DE  2
		// F   3
		A := coord{0 * height, 1 * width}
		B := coord{0 * height, 2 * width}
		C := coord{1 * height, 1 * width}
		D := coord{2 * height, 0 * width}
		E := coord{2 * height, 1 * width}
		F := coord{3 * height, 0 * width}

		if t.r <= -1 {
			// Hit top edge, travelling Up.
			//  D: C Right, {r,c} -> {c, min}  (transposed)
			//  A: F Right, {r,c} -> {c, min}  (transposed)
			//  B: F Up,    {r,c} -> {max, c}  (no transform)
			switch t.c - relC {
			case D.c:
				t.facing = right
				t.coord = coord{C.r + relC, -1}
			case A.c:
				t.facing = right
				t.coord = coord{F.r + relC, -1}
			case B.c:
				t.facing = up
				t.coord = coord{p.maxR, F.c + relC}
			}
		} else if t.r >= p.maxR {
			// Hit bottom edge, travelling Down
			//  F: B Down,  {r,c} -> {min, c}  (no transform)
			//  E: F Left,  {r,c} -> {c, max}  (transposed)
			//  B: C Left,  {r,c} -> {c, max}  (transposed)
			switch t.c - relC {
			case F.c:
				t.facing = down
				t.coord = coord{-1, B.c + relC}
			case E.c:
				t.facing = left
				t.coord = coord{F.r + relC, p.maxC}
			case B.c:
				t.facing = left
				t.coord = coord{C.r + relC, p.maxC}
			}
		} else if t.c <= -1 {
			// Hit Left edge
			//  A: D Right, {r,c} -> {-r, min} (upside down)
			//  C: D Down,  {r,c} -> {min, r}  (transposed)
			//  D: A Right, {r,c} -> {-r, min} (upside down)
			//  F: A Down,  {r,c} -> {min, r}  (transposed)
			switch t.r - relR {
			case A.r:
				t.facing = right
				t.coord = coord{D.r + (height - 1 - relR), -1}
			case C.r:
				t.facing = down
				t.coord = coord{-1, D.c + relR}
			case D.r:
				t.facing = right
				t.coord = coord{A.r + (height - 1 - relR), -1}
			case F.r:
				t.facing = down
				t.coord = coord{-1, A.c + relR}
			}
		} else if t.c >= p.maxC {
			// Hit Right edge
			//  B: E Left,  {r,c} -> {-r, max} (upside down)
			//  C: B Up,    {r,c} -> {max, r}  (transposed)
			//  E: B Left,  {r,c} -> {-r, max} (upside down)
			//  F: E Up,    {r,c} -> {max, r}  (transposed)
			switch t.r - relR {
			case B.r:
				t.facing = left
				t.coord = coord{E.r + (height - 1 - relR), p.maxC}
			case C.r:
				t.facing = up
				t.coord = coord{p.maxR, B.c + relR}
			case E.r:
				t.facing = left
				t.coord = coord{B.r + (height - 1 - relR), p.maxC}
			case F.r:
				t.facing = up
				t.coord = coord{p.maxR, E.c + relR}
			}
		} else {
			return false
		}
	}
	return true
}

type turtle struct {
	coord
	facing dir
}

func (t *turtle) forward(passable board, n int, cube bool) {
	var i int
	rollback := *t
	var wrapped bool
	for i < n {
		// Provisionally move us one square with wraparound.
		proposed := *t
		proposed.step()
		if proposed.enforceBounds(passable, cube) {
			wrapped = true
		}

		if p, ok := passable.tiles[proposed.coord]; !ok {
			// This isn't a real tile. slide riiiight on over without incrementing i
			// or changing rollback value
			*t = proposed
		} else if !p {
			// We've hit a wall, we can't proceed regardless of having more allowed movement.
			// Roll back to the last known non-void position rather than stranding us in the void.
			*t = rollback
			break
		} else {
			i++
			// Change both pos and rollback
			if *debug && wrapped {
				fmt.Printf("Wrap from %s to %s\n", rollback.String(passable), proposed.String(passable))
			}
			wrapped = false
			rollback = proposed
			*t = proposed
		}
	}
}

func (t turtle) password() int {
	sum := 1000 * (t.r + 1)
	sum += 4 * (t.c + 1)
	sum += int(t.facing)
	return sum
}

func main() {
	flag.Parse()
	bytes, err := ioutil.ReadFile(*inputFile)
	if err != nil {
		return
	}
	contents := string(bytes)
	split := strings.Split(contents, "\n")

	var passable board
	passable.tiles = make(map[coord]bool)
	for r, line := range split[:len(split)-3] {
		if r > passable.maxR {
			passable.maxR = r + 1
		}
		for c, v := range line {
			if c > passable.maxC {
				passable.maxC = c + 1
			}
			switch v {
			case ' ':
				continue
			case '#':
				passable.tiles[coord{r, c}] = false
			case '.':
				passable.tiles[coord{r, c}] = true
			}
		}
	}

	// part A
	fmt.Println(passable.run(split[len(split)-2], false))
	// part B
	fmt.Println(passable.run(split[len(split)-2], true))
}

func (b board) run(commands string, cube bool) int {
	player := turtle{}
	player.forward(b, 1, cube)

	for i := 0; i < len(commands); i++ {
		switch commands[i] {
		case 'L':
			player.facing = (player.facing + 3) % 4
		case 'R':
			player.facing = (player.facing + 1) % 4
		default:
			// This is a numeric value that we need to fast forward until we encounter a non-numeric value.
			j := 1
			for i+j < len(commands) && commands[i+j] != 'L' && commands[i+j] != 'R' {
				j++
			}
			v, _ := strconv.Atoi(commands[i : i+j])
			player.forward(b, v, cube)
			i += j - 1
		}
	}
	return player.password()
}
