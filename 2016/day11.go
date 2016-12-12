package main

import (
	"fmt"
	"sort"
	"sync"
)

const NumFloors int8 = 4

type FullItem struct {
	Isotope   int8
	Generator bool
	Floor     int8
}

type Item int8

func (i Item) Isotope() int8 {
	return int8(i) % 8
}

func (i Item) Generator() bool {
	return ((int8(i) >> 3) % 2) == 1
}

func (i Item) Floor() int8 {
	return (int8(i) >> 4) % 4
}
func (i Item) UpdateFloor(floor int8) Item {
	ret := int8(i)
	ret &= 0xf
	ret |= floor << 4

	return Item(ret)
}

func (f FullItem) Serialize() Item {
	var ret int8
	ret |= f.Floor << 4
	if f.Generator {
		ret |= 1 << 3
	}
	ret |= f.Isotope % 8
	return Item(ret)
}

type Board struct {
	Items         [14]Item
	ElevatorFloor int8
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
		if item.Isotope() == 0 {
			continue
		}
		if item.Floor() != NumFloors-1 {
			return false
		}
	}
	return true
}

func (b *Board) Valid() int {
	score := 0
	for f := int8(0); f < NumFloors; f++ {
		gens := make(map[int8]bool)
		for _, item := range b.Items {
			if item.Floor() != f || item.Isotope() == 0 {
				continue
			}
			if item.Generator() {
				gens[item.Isotope()] = true
				if f == NumFloors-1 {
					score += 20
				}
			}
		}
		for _, item := range b.Items {
			if item.Floor() != f || item.Isotope() == 0 {
				continue
			}
			if !item.Generator() {
				if gens[item.Isotope()] {
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

func (b *Board) MakeCopy(newFloor int8, move1, move2 int) Board {
	nb := Board{
		ElevatorFloor: newFloor,
		Items:         [14]Item{},
	}
	for i := range b.Items {
		if i == move1 || i == move2 {
			nb.Items[i] = b.Items[i].UpdateFloor(newFloor)
		} else {
			nb.Items[i] = b.Items[i]
		}
	}
	return nb
}

func (b *Board) MakeMoves() BoardList {
	ret := make(BoardList, 0)
	for nf := int8(0); nf < NumFloors; nf++ {
		if nf == b.ElevatorFloor || nf < b.ElevatorFloor-1 || nf > b.ElevatorFloor+1 {
			continue
		}
		for i, x := range b.Items {
			if x.Floor() != b.ElevatorFloor || x.Isotope() == 0 {
				continue
			}
			child := b.MakeCopy(nf, i, -1)
			if score := child.Valid(); score >= 0 {
				ret = append(ret, BoardScore{child, score})
			}
		}
		for i, x := range b.Items {
			if x.Floor() != b.ElevatorFloor || x.Isotope() == 0 {
				continue
			}
			for j, y := range b.Items {
				if i <= j || y.Floor() != b.ElevatorFloor || y.Isotope() == 0 {
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

func (b Board) EvalLoop(seen map[Board]uint16, winner *uint16, wMtx *sync.Mutex, wg *sync.WaitGroup) {
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

func (b Board) ProcessBoard() uint16 {
	b.Valid()

	winner := uint16(0)

	seen := make(map[Board]uint16)
	seen[b] = uint16(1)
	var mtx sync.Mutex
	var wg sync.WaitGroup

	wg.Add(1)
	b.EvalLoop(seen, &winner, &mtx, &wg)
	wg.Wait()

	fmt.Println(len(seen))
	return winner - seen[b]
}

func main() {
	hydrogen := int8(1)
	lithium := int8(2)
	demo := Board{
		ElevatorFloor: 0,
		Items: [14]Item{
			FullItem{hydrogen, false, 0}.Serialize(), FullItem{lithium, false, 0}.Serialize(),
			FullItem{hydrogen, true, 1}.Serialize(),
			FullItem{lithium, true, 2}.Serialize(),
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
			FullItem{thulium, true, 0}.Serialize(), FullItem{thulium, false, 0}.Serialize(), FullItem{plutonium, true, 0}.Serialize(), FullItem{strontium, true, 0}.Serialize(),
			FullItem{plutonium, false, 1}.Serialize(), FullItem{strontium, false, 1}.Serialize(),
			FullItem{promethium, true, 2}.Serialize(), FullItem{promethium, false, 2}.Serialize(), FullItem{ruthenium, true, 2}.Serialize(), FullItem{ruthenium, false, 2}.Serialize(),
		},
	}
	fmt.Println(state.ProcessBoard())

	elerium := int8(6)
	dilithium := int8(7)
	extra := Board{
		ElevatorFloor: 0,
		Items: [14]Item{
			FullItem{thulium, true, 0}.Serialize(), FullItem{thulium, false, 0}.Serialize(), FullItem{plutonium, true, 0}.Serialize(), FullItem{strontium, true, 0}.Serialize(),
			FullItem{plutonium, false, 1}.Serialize(), FullItem{strontium, false, 1}.Serialize(),
			FullItem{promethium, true, 2}.Serialize(), FullItem{promethium, false, 2}.Serialize(), FullItem{ruthenium, true, 2}.Serialize(), FullItem{ruthenium, false, 2}.Serialize(),
			FullItem{elerium, false, 0}.Serialize(), FullItem{elerium, true, 0}.Serialize(), FullItem{dilithium, false, 0}.Serialize(), FullItem{dilithium, true, 0}.Serialize(),
		},
	}
	fmt.Println(extra.ProcessBoard())
}
