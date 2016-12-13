package main

import (
	"fmt"
	"sort"
	"sync"
)

const winX int8 = 31
const winY int8 = 39
const printBoard bool = false

type Coord struct {
	X, Y int8
}

// Constants needed for Popcount64
const k1 uint64 = 6148914691236517205
const k2 uint64 = 3689348814741910323
const k4 uint64 = 1085102592571150095

func Popcount64(x uint64) uint8 {
	x = x - ((x >> 1) & k1)
	x = (x & k2) + ((x >> 2) & k2)
	x = (x + (x >> 4)) & k4
	x = (x * 72340172838076673) >> 56
	return uint8(x)
}

func (c Coord) IsPassable() bool {
	if c.X < 0 || c.Y < 0 {
		return false
	}
	x := uint64(c.X)
	y := uint64(c.Y)
	val := x*x + 3*x + 2*x*y + y + y*y + 1358
	return Popcount64(val)%2 == 0
}

type BoardScore struct {
	Child Coord
	Score int
}

type BoardList []BoardScore

func (a BoardList) Len() int      { return len(a) }
func (a BoardList) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a BoardList) Less(i, j int) bool {
	return a[i].Score < a[j].Score
}

func (c Coord) Win() bool {
	return c.X == winX && c.Y == winY
}

func (c Coord) Valid() int {
	if !c.IsPassable() {
		return -1
	}
	xDelta := int(c.X - winX)
	yDelta := int(c.Y - winY)
	return xDelta*xDelta + yDelta*yDelta
}

func (c Coord) MakeMoves() BoardList {
	ret := make(BoardList, 0)
	up := Coord{c.X, c.Y + 1}
	down := Coord{c.X, c.Y - 1}
	left := Coord{c.X - 1, c.Y}
	right := Coord{c.X + 1, c.Y}
	if score := up.Valid(); score >= 0 {
		ret = append(ret, BoardScore{up, score})
	}
	if score := down.Valid(); score >= 0 {
		ret = append(ret, BoardScore{down, score})
	}
	if score := left.Valid(); score >= 0 {
		ret = append(ret, BoardScore{left, score})
	}
	if score := right.Valid(); score >= 0 {
		ret = append(ret, BoardScore{right, score})
	}
	sort.Sort(ret)
	return ret
}

func (c Coord) EvalLoop(seen map[Coord]uint16, winner *uint16, wMtx *sync.Mutex, wg *sync.WaitGroup) {
	children := c.MakeMoves()
	for _, childPair := range children {
		child := childPair.Child
		wMtx.Lock()
		if (seen[child] == 0 || seen[child] > seen[c]+1) && (*winner == 0 || seen[c]+1 < *winner) {
			seen[child] = seen[c] + 1
			if child.Win() {
				*winner = seen[child]
				wMtx.Unlock()
				continue
			}
			wMtx.Unlock()
			wg.Add(1)
			go child.EvalLoop(seen, winner, wMtx, wg)
		} else {
			wMtx.Unlock()
			// We've seen this before in equal or more moves OR we've won already in fewer moves.
		}
	}
	wg.Done()
}

func (c Coord) ProcessBoard() uint16 {
	winner := uint16(0)

	seen := make(map[Coord]uint16)
	seen[c] = uint16(1)
	var mtx sync.Mutex
	var wg sync.WaitGroup

	wg.Add(1)
	c.EvalLoop(seen, &winner, &mtx, &wg)
	wg.Wait()

	fmt.Println(len(seen))
	in50 := 0
	for _, steps := range seen {
		if steps - seen[c] <= 50 {
			in50++
		}
	}
	fmt.Println(in50)
	return winner - seen[c]
}

func main() {
	start := Coord{1, 1}
	if printBoard {
		for y := int8(0); y < winY * 2; y++ {
			for x := int8(0); x < winX * 2; x++ {
				passable := Coord{x, y}.IsPassable()
				if x == start.X && y == start.Y {
						fmt.Printf("S")
				} else if x == winX && y == winY {
					fmt.Printf("F")
				} else if passable {
					fmt.Printf(" ")
				} else {
					fmt.Printf("#")
				}
			}
			fmt.Println()
		}
	}
	fmt.Println(start.ProcessBoard())
}
