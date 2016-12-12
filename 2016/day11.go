package main

import (
	"fmt"
	"sort"
	"sync"
)

const NumFloors int = 4

type Item struct {
	Isotope   int8
	Generator bool
	Floor     int
}

type Board struct {
	Items         [14]Item
	ElevatorFloor int
}

type BoardScore struct {
	Child Board
	Score int
}

type BoardList []BoardScore

func (a BoardList) Len() int      { return len(a) }
func (a BoardList) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a BoardList) Less(i, j int) bool {
	return a[i].Score > a[j].Score
}

func (b *Board) Win() bool {
	for _, item := range b.Items {
		if item.Isotope == 0 {
			continue
		}
		if item.Floor != NumFloors-1 {
			return false
		}
	}
	return true
}

func (b *Board) Valid() int {
	score := 0
	for f := 0; f < NumFloors; f++ {
		gens := make(map[int8]bool)
		for _, item := range b.Items {
			if item.Floor != f || item.Isotope == 0 {
				continue
			}
			if item.Generator {
				gens[item.Isotope] = true
				if f == NumFloors-1 {
					score += 20
				}
			}
		}
		for _, item := range b.Items {
			if item.Floor != f || item.Isotope == 0 {
				continue
			}
			if !item.Generator {
				if gens[item.Isotope] {
					// Shielded
					score += 10
					continue
				}
				if len(gens) != 0 {
					// Not shielded, and at least one other generator present
					return -1
				}
			}
		}
	}
	return score
}

func (b *Board) MakeCopy(newFloor int, move1, move2 int) Board {
	nb := Board{
		ElevatorFloor: newFloor,
		Items:         [14]Item{},
	}
	for i := range b.Items {
		if i == move1 || i == move2 {
			newItem := b.Items[i]
			newItem.Floor = newFloor
			nb.Items[i] = newItem
		} else {
			nb.Items[i] = b.Items[i]
		}
	}
	return nb
}

func (b *Board) MakeMoves() BoardList {
	ret := make(BoardList, 0)
	// fmt.Println(b)
	for nf := 0; nf < NumFloors; nf++ {
		if nf == b.ElevatorFloor || nf < b.ElevatorFloor-1 || nf > b.ElevatorFloor+1 {
			continue
		}
		for i, x := range b.Items {
			if x.Floor != b.ElevatorFloor || x.Isotope == 0 {
				continue
			}
			child := b.MakeCopy(nf, i, -1)
			if score := child.Valid(); score >= 0 {
				ret = append(ret, BoardScore{child, score})
			}
		}
		for i, x := range b.Items {
			if x.Floor != b.ElevatorFloor || x.Isotope == 0 {
				continue
			}
			for j, y := range b.Items {
				if i <= j || y.Floor != b.ElevatorFloor || y.Isotope == 0 {
					continue
				}
				child := b.MakeCopy(nf, i, j)
				if score := child.Valid(); score >= 0 {
					ret = append(ret, BoardScore{child, score})
				}
			}
		}
	}
	sort.Sort(ret)
	return ret
}

func (b Board) EvalLoop(seen map[Board]int, winner *int, wMtx *sync.Mutex, wg *sync.WaitGroup) {
	children := b.MakeMoves()
	for _, childPair := range children {
		child := childPair.Child
		wMtx.Lock()
		if (seen[child] == 0 || seen[child] > seen[b]+1) && (*winner == 0 || seen[b]+1 < *winner) {
			seen[child] = seen[b] + 1
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

func (b Board) ProcessBoard() int {
	b.Valid()

	winner := 0

	seen := make(map[Board]int)
	seen[b] = 1
	var mtx sync.Mutex
	var wg sync.WaitGroup

	wg.Add(1)
	b.EvalLoop(seen, &winner, &mtx, &wg)
	wg.Wait()
	return winner - seen[b]
}

func main() {
	hydrogen := int8(1)
	lithium := int8(2)
	demo := Board{
		ElevatorFloor: 0,
		Items: [14]Item{
			{hydrogen, false, 0}, {lithium, false, 0},
			{hydrogen, true, 1},
			{lithium, true, 2},
		},
	}
	fmt.Println(demo.ProcessBoard())

	thulium := int8(1)
	plutonium := int8(2)
	strontium := int8(3)
	promethium := int8(4)
	ruthenium := int8(5)
	state := Board{
		ElevatorFloor: 0,
		Items: [14]Item{
			{thulium, true, 0}, {thulium, false, 0}, {plutonium, true, 0}, {strontium, true, 0},
			{plutonium, false, 1}, {strontium, false, 1},
			{promethium, true, 2}, {promethium, false, 2}, {ruthenium, true, 2}, {ruthenium, false, 2},
		},
	}
	fmt.Println(state.ProcessBoard())

	elerium := int8(6)
	dilithium := int8(7)
	extra := Board{
		ElevatorFloor: 0,
		Items: [14]Item{
			{thulium, true, 0}, {thulium, false, 0}, {plutonium, true, 0}, {strontium, true, 0},
			{plutonium, false, 1}, {strontium, false, 1},
			{promethium, true, 2}, {promethium, false, 2}, {ruthenium, true, 2}, {ruthenium, false, 2},
			{elerium, false, 0}, {elerium, true, 0}, {dilithium, false, 0}, {dilithium, true, 0},
		},
	}
	fmt.Println(extra.ProcessBoard())
}
