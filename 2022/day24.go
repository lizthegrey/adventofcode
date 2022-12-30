package main

import (
	"container/heap"
	"flag"
	"fmt"
	"io/ioutil"
	"strings"
)

var inputFile = flag.String("inputFile", "inputs/day24.input", "Relative file path to use as input.")

type coord struct {
	r, c int
}

type coord3 struct {
	coord
	t int
}

type dir int

const (
	up = int8(1 << iota)
	right
	down
	left
)

type terrain struct {
	height, width int
	storms        map[coord]int8
}

type generator map[int]terrain

func memo(seed terrain) generator {
	return generator{
		0: seed,
	}
}

func (g generator) generate(t int) terrain {
	if cached, ok := g[t]; ok {
		return cached
	}
	ret := g.generate(t - 1).iterate()
	g[t] = ret
	return ret
}

func (g generator) neighbours(src coord3) []coord3 {
	board := g.generate(src.t + 1)
	ret := board.neighbours(src.coord)
	for i := range ret {
		ret[i].t = src.t + 1
	}
	return ret
}

func (t terrain) iterate() terrain {
	next := terrain{
		height: t.height,
		width:  t.width,
		storms: make(map[coord]int8),
	}
	for loc, mask := range t.storms {
		if mask&left == left {
			next.storms[coord{loc.r, (loc.c + t.width - 1) % t.width}] |= left
		}
		if mask&right == right {
			next.storms[coord{loc.r, (loc.c + 1) % t.width}] |= right
		}
		if mask&up == up {
			next.storms[coord{(loc.r + t.height - 1) % t.height, loc.c}] |= up
		}
		if mask&down == down {
			next.storms[coord{(loc.r + 1) % t.height, loc.c}] |= down
		}
	}
	return next
}

func (t terrain) traversable(c coord) bool {
	return t.storms[c] == 0
}

func (c coord) minDst(o coord) int {
	var sum int
	if c.r > o.r {
		sum += c.r - o.r
	} else {
		sum += o.r - c.r
	}
	if c.c > o.c {
		sum += c.c - o.c
	} else {
		sum += o.c - c.c
	}
	return sum
}

type HeapQueue struct {
	elems     *[]coord3
	score     map[coord3]int
	positions map[coord3]int
}

func (h HeapQueue) Len() int           { return len(*h.elems) }
func (h HeapQueue) Less(i, j int) bool { return h.score[(*h.elems)[i]] < h.score[(*h.elems)[j]] }
func (h HeapQueue) Swap(i, j int) {
	h.positions[(*h.elems)[i]], h.positions[(*h.elems)[j]] = h.positions[(*h.elems)[j]], h.positions[(*h.elems)[i]]
	(*h.elems)[i], (*h.elems)[j] = (*h.elems)[j], (*h.elems)[i]
}

func (h HeapQueue) Push(x interface{}) {
	h.positions[x.(coord3)] = len(*h.elems)
	*h.elems = append(*h.elems, x.(coord3))
}

func (h HeapQueue) Pop() interface{} {
	old := *h.elems
	n := len(old)
	x := old[n-1]
	*h.elems = old[0 : n-1]
	delete(h.positions, x)
	return x
}

func (h HeapQueue) Position(x coord3) int {
	if pos, ok := h.positions[x]; ok {
		return pos
	}
	return -1
}

func aStar(gen generator, src coord3, dst coord) int {
	visited := map[coord3]bool{
		src: true,
	}
	fScore := map[coord3]int{
		src: src.minDst(dst),
	}
	workList := HeapQueue{&[]coord3{src}, fScore, make(map[coord3]int)}
	heap.Init(&workList)

	for len(*workList.elems) != 0 {
		// Pop the current node off the worklist.
		current := heap.Pop(&workList).(coord3)

		if current.coord == dst {
			return current.t
		}
		for _, n := range gen.neighbours(current) {
			if !visited[n] {
				visited[n] = true
				fScore[n] = n.t + n.minDst(dst)
				if pos := workList.Position(n); pos == -1 {
					heap.Push(&workList, n)
				} else {
					heap.Fix(&workList, pos)
				}
			}
		}
	}
	return -1
}

func (t terrain) start() coord {
	return coord{-1, 0}
}

func (t terrain) end() coord {
	return coord{t.height, t.width - 1}
}

func (t terrain) canMove(dst coord) bool {
	// the start/end tiles are always passable.
	if dst == t.start() || dst == t.end() {
		return true
	}
	// if otherwise outside of bounds, return false
	if dst.r < 0 || dst.r >= t.height || dst.c < 0 || dst.c >= t.width {
		return false
	}
	// otherwise defer to the winds.
	return t.traversable(dst)
}

func (t terrain) maybeAppend(list []coord3, dst coord) []coord3 {
	if t.canMove(dst) {
		list = append(list, coord3{coord: dst})
	}
	return list
}

func (t terrain) neighbours(src coord) []coord3 {
	var ret []coord3
	// Staying still is a valid (and useful) move.
	ret = t.maybeAppend(ret, coord{src.r + 0, src.c + 0})

	ret = t.maybeAppend(ret, coord{src.r + 0, src.c + 1})
	ret = t.maybeAppend(ret, coord{src.r + 0, src.c - 1})
	ret = t.maybeAppend(ret, coord{src.r + 1, src.c + 0})
	ret = t.maybeAppend(ret, coord{src.r - 1, src.c + 0})
	return ret
}

func main() {
	flag.Parse()
	bytes, err := ioutil.ReadFile(*inputFile)
	if err != nil {
		return
	}
	contents := string(bytes)
	split := strings.Split(contents, "\n")

	storms := make(map[coord]int8)
	var height, width int
	for r, s := range split[1 : len(split)-2] {
		height += 1
		width = 0
		for c, char := range s[1 : len(s)-1] {
			width += 1
			loc := coord{r, c}
			switch char {
			case '^':
				storms[loc] |= up
			case '<':
				storms[loc] |= left
			case '>':
				storms[loc] |= right
			case 'v':
				storms[loc] |= down
			}
		}
	}

	seed := terrain{
		height: height,
		width:  width,
		storms: storms,
	}

	// Part A
	gen := memo(seed)
	steps := aStar(gen, coord3{coord: seed.start(), t: 0}, seed.end())
	fmt.Println(steps)

	// Part B
	steps = aStar(gen, coord3{coord: seed.end(), t: steps}, seed.start())
	steps = aStar(gen, coord3{coord: seed.start(), t: steps}, seed.end())
	fmt.Println(steps)
}
