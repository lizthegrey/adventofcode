package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var inputFile = flag.String("inputFile", "inputs/day17.input", "Relative file path to use as input.")
var partB = flag.Bool("partB", false, "Whether to use the Part B logic.")

type Coord struct {
	X, Y int
}

type Items map[Coord]bool

func main() {
	flag.Parse()
	f, err := os.Open(*inputFile)
	if err != nil {
		return
	}
	defer f.Close()

	r := bufio.NewReader(f)

	clay := make(Items)
	minX := 500
	maxX := 500
	maxY := 0

	for {
		l, err := r.ReadString('\n')
		if err != nil || len(l) == 0 {
			break
		}
		scan := strings.Split(l[:len(l)-1], ", ")

		coord, _ := strconv.Atoi(scan[0][2:])
		bounds := strings.Split(scan[1][2:], "..")
		lbound, _ := strconv.Atoi(bounds[0])
		ubound, _ := strconv.Atoi(bounds[1])

		if scan[0][0] == 'x' {
			x := coord
			if minX > x {
				minX = x
			}
			if maxX < x {
				maxX = x
			}
			if ubound > maxY {
				maxY = ubound
			}
			for y := lbound; y <= ubound; y++ {
				clay[Coord{x, y}] = true
			}
		} else {
			y := coord
			if minX > lbound {
				minX = lbound
			}
			if maxX < ubound {
				maxX = ubound
			}
			if y > maxY {
				maxY = y
			}
			for x := lbound; x <= ubound; x++ {
				clay[Coord{x, y}] = true
			}
		}
	}

	water := make(Items)
	traversed := make(Items)
	spring := Coord{500, 0}

	for prevLen := -1; len(water)+len(traversed) != prevLen; {
		prevLen = len(water) + len(traversed)
		dropFall(traversed, water, clay, spring, minX, maxX, maxY)
		fmt.Printf("Placed drops %d\n", len(water))
	}

	invalid := make([]Coord, 0)
	for k := range water {
		valid := true
		for x := k.X; x >= minX; x-- {
			left := Coord{x, k.Y}
			if water[left] {
				continue
			}
			if !clay[left] {
				valid = false
			}
			break
		}
		for x := k.X; x <= maxX; x++ {
			right := Coord{x, k.Y}
			if water[right] {
				continue
			}
			if !clay[right] {
				valid = false
			}
			break
		}
		if !valid {
			invalid = append(invalid, k)
		}
	}
	for _, v := range invalid {
		delete(water, v)
	}

	for y := 0; y <= maxY; y++ {
		for x := minX; x <= maxX; x++ {
			loc := Coord{x, y}
			if loc == spring {
				fmt.Printf("+")
			}
			if clay[loc] {
				fmt.Printf("#")
			} else if water[loc] {
				fmt.Printf("~")
			} else if traversed[loc] {
				fmt.Printf("|")
			} else {
				fmt.Printf(".")
			}
		}
		fmt.Println()
	}

	fmt.Println(len(traversed))
}

func dropFall(traversed, water, clay Items, spring Coord, minX, maxX, maxY int) {
	worklist := []Coord{spring}
	for len(worklist) != 0 {
		loc := worklist[0]
		newWork := loc.fallInner(traversed, water, clay, spring, minX, maxX, maxY)
		worklist = append(newWork, worklist[1:]...)
	}
}

func (drop Coord) fallInner(traversed, water, clay Items, spring Coord, minX, maxX, maxY int) []Coord {
	iter := make([]Coord, 0)

	if drop != spring {
		traversed[drop] = true
	}

	fall := drop
	fall.Y++
	if fall.Y > maxY {
		return iter
	}
	if clay[fall] {
		water[drop] = true
		return iter
	}
	if water[fall] {
		for x := fall.X; x >= minX; x-- {
			left := Coord{x, fall.Y}
			if water[left] {
				continue
			}
			if clay[left] {
				break
			}
			iter = append(iter, left)
			break
		}
		for x := fall.X; x <= maxX; x++ {
			right := Coord{x, fall.Y}
			if water[right] {
				continue
			}
			if clay[right] {
				break
			}
			iter = append(iter, right)
			break
		}
		if len(iter) == 0 {
			// Both sides cannot be iterated upon.
			traversed[drop] = true
			water[drop] = true
			return iter
		}
	} else {
		iter = append(iter, fall)
	}
	return iter
}
