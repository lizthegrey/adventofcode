package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strconv"
)

var inputFile = flag.String("inputFile", "inputs/day03.input", "Relative file path to use as input.")

var reg = regexp.MustCompile("#([0-9]+) @ ([0-9]+),([0-9]+): ([0-9]+)x([0-9]+)")

type Claim struct {
	top, bottom, left, right int // inclusive
}

type Coord struct {
	X, Y int
}

func main() {
	flag.Parse()
	f, err := os.Open(*inputFile)
	if err != nil {
		return
	}
	defer f.Close()

	claims := make(map[int]Claim)

	r := bufio.NewReader(f)
	for {
		l, err := r.ReadString('\n')
		if err != nil || len(l) == 0 {
			break
		}
		l = l[:len(l)-1]
		parsed := reg.FindStringSubmatch(l)
		claim, err := strconv.Atoi(parsed[1])
		offsetX, err := strconv.Atoi(parsed[2])
		offsetY, err := strconv.Atoi(parsed[3])
		sizeX, err := strconv.Atoi(parsed[4])
		sizeY, err := strconv.Atoi(parsed[5])
		claims[claim] = Claim{offsetY, offsetY + sizeY - 1, offsetX, offsetX + sizeX - 1}
	}

	occupied := make(map[Coord][]int)
	overlaps := make(map[Coord]bool)
	conflicted := make(map[int]bool)
	for k, v := range claims {
		for x := v.left; x <= v.right; x++ {
			for y := v.top; y <= v.bottom; y++ {
				if occupied[Coord{x, y}] != nil {
					overlaps[Coord{x, y}] = true
					occupied[Coord{x, y}] = append(occupied[Coord{x, y}], k)
					for _, id := range occupied[Coord{x, y}] {
						conflicted[id] = true
					}
				} else {
					occupied[Coord{x, y}] = []int{k}
				}
			}
		}
	}

	fmt.Printf("Result is %d\n", len(overlaps))
	for k := range claims {
		if !conflicted[k] {
			fmt.Printf("Not conflicted: %d\n", k)
		}
	}
}
