package main

import (
	"flag"
	"fmt"
	"github.com/lizthegrey/adventofcode/2019/intcode"
)

var inputFile = flag.String("inputFile", "inputs/day19.input", "Relative file path to use as input.")

const maxSize int = 10000
const squareSize int = 100

type Coord struct {
	X, Y int
}

type Range struct {
	Start, End int
	Invalid    bool
}

func (r Range) Size() int {
	return r.End - r.Start
}
func (r Range) Contains(o Range) bool {
	return o.Start >= r.Start && o.End <= r.End
}

func main() {
	flag.Parse()
	tape := intcode.ReadInput(*inputFile)
	if tape == nil {
		fmt.Println("Failed to parse input.")
		return
	}
	// Part A
	viewA := computePoints(tape, 50)
	sum := 0
	for _, v := range viewA {
		if v.Invalid {
			continue
		}
		sum += v.Size() + 1
	}
	fmt.Println(sum)

	// Part B
	view := computePoints(tape, maxSize+squareSize)
	found := findSquares(view, squareSize)
	x := found.X
	y := found.Y
	fmt.Printf("100x100 square found at %d,%d. answer is %d.\n", x, y, x*10000+y)
}

func findSquares(view []Range, size int) Coord {
	for y, v := range view {
		if v.Invalid || view[y+size-1].Invalid {
			continue
		}
		if v.Size()+1 < size {
			continue
		}
		if !v.Contains(Range{view[y+size-1].Start, view[y+size-1].Start + size - 1, true}) {
			continue
		}
		return Coord{view[y+size-1].Start, y}
	}
	return Coord{-1, -1}
}

func computePoints(tape intcode.Tape, nDim int) []Range {
	view := make([]Range, nDim, nDim)
	for y := 0; y < nDim; y++ {
		if y%1000 == 0 {
			fmt.Printf("Computed up to row %d\n", y)
		}
		last := Range{0, 0, true}
		if y > 1 {
			last = view[y-1]
		}
		x := 0
		if !last.Invalid && last.Start > 0 {
			x = last.Start - 1
		}
		for x <= y && !computePoint(tape, Coord{x, y}) {
			x++
		}
		if x > y {
			view[y].Invalid = true
			continue
		}
		view[y].Start = x
		x = y
		if !last.Invalid {
			x = last.End + 2
		}
		for x >= 0 && !computePoint(tape, Coord{x, y}) {
			x--
		}
		view[y].End = x
	}
	return view
}

func computePoint(tape intcode.Tape, c Coord) bool {
	workingTape := tape.Copy()
	input := make(chan int, 2)
	output, _ := workingTape.Process(input)
	input <- c.X
	input <- c.Y
	return <-output == 1
}
