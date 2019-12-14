package main

import (
	"flag"
	"fmt"
	"github.com/lizthegrey/adventofcode/2019/intcode"
	"math"
	"sync"
	"time"
)

var inputFile = flag.String("inputFile", "inputs/day13.input", "Relative file path to use as input.")
var debug = flag.Bool("debug", false, "Whether to print debug output.")
var partB = flag.Bool("partB", false, "Whether to use Part B logic.")

type Coord struct {
	X, Y int
}

type Board map[Coord]int

func (b Board) Display() (ballX, paddleX int) {
	minX := math.MaxInt32
	maxX := math.MinInt32
	minY := math.MaxInt32
	maxY := math.MinInt32

	for coord := range b {
		if coord.X > maxX {
			maxX = coord.X
		}
		if coord.X < minX {
			minX = coord.X
		}
		if coord.Y > maxY {
			maxY = coord.Y
		}
		if coord.Y < minY {
			minY = coord.Y
		}
	}
	if *debug {
		fmt.Print("\033[2J")
	}
	for y := minY; y <= maxY; y++ {
		for x := minX; x <= maxX; x++ {
			if *debug {
				switch b[Coord{x, y}] {
				case 0:
					fmt.Printf(" ")
				case 1:
					fmt.Printf("|")
				case 2:
					fmt.Printf("x")
				case 3:
					fmt.Printf("-")
				case 4:
					fmt.Printf("o")
				}
			}
			switch b[Coord{x, y}] {
			case 3:
				paddleX = x
			case 4:
				ballX = x
			}
		}
		if *debug {
			fmt.Println()
		}
	}
	return
}

func main() {
	flag.Parse()
	tape := intcode.ReadInput(*inputFile)
	if tape == nil {
		fmt.Println("Failed to parse input.")
		return
	}

	workingTape := tape.Copy()
	input := make(chan int, 1)
	output, done := workingTape.Process(input)

	tiles := make(Board)
outer:
	for {
		select {
		case <-done:
			break outer
		default:
			x := <-output
			y := <-output
			tile := <-output
			tiles[Coord{x, y}] = tile
		}
	}

	blockCount := 0
	for _, v := range tiles {
		if v == 2 {
			blockCount++
		}
	}

	fmt.Println(blockCount)

	if *partB {
		play(tape)
	}
}

func play(tape intcode.Tape) {
	input := make(chan int, 1)
	tape[0] = 2
	output, done := tape.Process(input)

	score := 0
	tiles := make(Board)

	var bMtx sync.Mutex
	go func() {
		for {
			var ballX, paddleX int
			time.Sleep(time.Millisecond)
			bMtx.Lock()
			if len(tiles) != 1080 {
				bMtx.Unlock()
				continue
			}
			ballX, paddleX = tiles.Display()
			fmt.Println(score)
			bMtx.Unlock()

			if ballX < paddleX {
				input <- -1
			} else if ballX == paddleX {
				input <- 0
			} else {
				input <- 1
			}
		}
	}()

outer:
	for {
		select {
		case <-done:
			break outer
		default:
			x := <-output
			y := <-output
			tile := <-output
			if x == -1 && y == 0 {
				score = tile
			} else {
				bMtx.Lock()
				tiles[Coord{x, y}] = tile
				bMtx.Unlock()
			}
		}
	}
	fmt.Println(score)
}
