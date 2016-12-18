package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"sort"
	"sync"
)

const winX int8 = 3
const winY int8 = 3

const seed string = "dmypynyp"

type Coord struct {
	X, Y    int8
	History []byte
}

func (c Coord) IsPassable() bool {
	if c.X < 0 || c.Y < 0 || c.X > winX || c.Y > winY {
		return false
	}
	past := c.History[:len(c.History)-1]

	in := fmt.Sprintf("%s%s", seed, string(past))
	h := md5.Sum([]byte(in))
	hash := hex.EncodeToString(h[:])

	open := func(b byte) bool {
		return b == 'b' || b == 'c' || b == 'd' || b == 'e' || b == 'f'
	}

	var ret bool
	switch c.History[len(c.History)-1] {
	case 'U':
		ret = open(hash[0])
	case 'D':
		ret = open(hash[1])
	case 'L':
		ret = open(hash[2])
	case 'R':
		ret = open(hash[3])
	}

	if ret {
		fmt.Printf("OK: %s\n", string(c.History))
	} else {
		fmt.Printf("XX: %s -> %s\n", string(c.History), hash[0:4])
	}

	return ret
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

	upHist := make([]byte, len(c.History)+1)
	copy(upHist, c.History)
	upHist[len(c.History)] = 'U'
	downHist := make([]byte, len(c.History)+1)
	copy(downHist, c.History)
	downHist[len(c.History)] = 'D'
	leftHist := make([]byte, len(c.History)+1)
	copy(leftHist, c.History)
	leftHist[len(c.History)] = 'L'
	rightHist := make([]byte, len(c.History)+1)
	copy(rightHist, c.History)
	rightHist[len(c.History)] = 'R'

	up := Coord{c.X, c.Y - 1, upHist}
	down := Coord{c.X, c.Y + 1, downHist}
	left := Coord{c.X - 1, c.Y, leftHist}
	right := Coord{c.X + 1, c.Y, rightHist}

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

func (c Coord) EvalLoop(winner *[]byte, wMtx *sync.Mutex, wg *sync.WaitGroup) {
	children := c.MakeMoves()
	for _, childPair := range children {
		child := childPair.Child
		wMtx.Lock()
		if len(*winner) == 0 || len(child.History) < len(*winner) {
			if child.Win() {
				*winner = child.History
				wMtx.Unlock()
				continue
			}
			wMtx.Unlock()
			wg.Add(1)
			go child.EvalLoop(winner, wMtx, wg)
		} else {
			wMtx.Unlock()
			// We've seen this before in equal or more moves OR we've won already in fewer moves.
		}
	}
	wg.Done()
}

func (c Coord) ProcessBoard() string {
	var winner []byte

	var mtx sync.Mutex
	var wg sync.WaitGroup

	wg.Add(1)
	c.EvalLoop(&winner, &mtx, &wg)
	wg.Wait()

	if winner != nil {
		return string(winner)
	}
	return "Not Found"
}

func main() {
	start := Coord{0, 0, []byte{}}
	fmt.Println(start.ProcessBoard())
}
