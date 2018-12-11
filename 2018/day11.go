package main

import (
	"flag"
	"fmt"
	"sync"
)

var serial = flag.Int("serial", 42, "The serial number of the power grid.")
var partB = flag.Bool("partB", false, "Whether to use the Part B logic.")

type Coord struct {
	X, Y int
}
type Square struct {
	X, Y, L int
}

func main() {
	flag.Parse()

	cells := make(map[Coord]int)
	for x := 1; x <= 300; x++ {
		for y := 1; y <= 300; y++ {
			rackID := x + 10
			level := (rackID * y) + *serial
			level *= rackID
			level %= 1000
			level /= 100
			level -= 5
			cells[Coord{x, y}] = level
		}
	}

	max := -1000
	var highest *Square
	var memo sync.Map
	if *partB {
		for squareSize := 1; squareSize <= 300; squareSize++ {
			score := evalSquare(cells, squareSize, &memo, max, &highest)
			if score > max {
				max = score
			}
		}
	} else {
		evalSquare(cells, 1, &memo, -1000, &highest)
		evalSquare(cells, 2, &memo, -1000, &highest)
		evalSquare(cells, 3, &memo, -1000, &highest)
	}

	fmt.Printf("Result is %d,%d,%d\n", highest.X, highest.Y, highest.L)
}

func evalSquare(cells map[Coord]int, squareSize int, memo *sync.Map, max int, highest **Square) int {
	ret := max
	var mtx sync.Mutex
	var wg sync.WaitGroup
	for x := 1; x <= 301-squareSize; x++ {
		wg.Add(1)
		go evalSquareInner(&mtx, x, &ret, cells, squareSize, memo, max, highest, &wg)
	}
	wg.Wait()
	return ret
}

func evalSquareInner(mtx *sync.Mutex, x int, ret *int, cells map[Coord]int, squareSize int, memo *sync.Map, max int, highest **Square, wg *sync.WaitGroup) {
	for y := 1; y <= 301-squareSize; y++ {
		cached, ok := memo.Load(Square{x, y, squareSize - 1})
		total := 0
		if ok {
			total = cached.(int)
		}

		for offsetX := 0; offsetX < squareSize; offsetX++ {
			total += cells[Coord{x + offsetX, y + squareSize - 1}]
		}
		for offsetY := 0; offsetY < squareSize; offsetY++ {
			total += cells[Coord{x + squareSize - 1, y + offsetY}]
		}
		total -= cells[Coord{x + squareSize - 1, y + squareSize - 1}]

		memo.Store(Square{x, y, squareSize}, total)

		mtx.Lock()
		if total > *ret {
			*ret = total
			*highest = &Square{x, y, squareSize}
		}
		mtx.Unlock()
	}
	wg.Done()
}
