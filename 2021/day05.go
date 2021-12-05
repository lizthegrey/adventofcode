package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

var inputFile = flag.String("inputFile", "inputs/day05.input", "Relative file path to use as input.")

type Coord struct {
	X, Y int
}

func main() {
	flag.Parse()
	bytes, err := ioutil.ReadFile(*inputFile)
	if err != nil {
		return
	}
	contents := string(bytes)
	split := strings.Split(contents, "\n")
	split = split[:len(split)-1]

	cloudsA := make(map[Coord]int)
	cloudsB := make(map[Coord]int)
	for _, s := range split {
		parts := strings.Split(s, " -> ")
		startRaw := strings.Split(parts[0], ",")
		endRaw := strings.Split(parts[1], ",")
		startX, _ := strconv.Atoi(startRaw[0])
		startY, _ := strconv.Atoi(startRaw[1])
		endX, _ := strconv.Atoi(endRaw[0])
		endY, _ := strconv.Atoi(endRaw[1])

		if startX == endX {
			if startY < endY {
				for y := startY; y <= endY; y++ {
					cloudsA[Coord{startX, y}]++
					cloudsB[Coord{startX, y}]++
				}
			} else {
				for y := startY; y >= endY; y-- {
					cloudsA[Coord{startX, y}]++
					cloudsB[Coord{startX, y}]++
				}
			}
		} else if startY == endY {
			if startX < endX {
				for x := startX; x <= endX; x++ {
					cloudsA[Coord{x, startY}]++
					cloudsB[Coord{x, startY}]++
				}
			} else {
				for x := startX; x >= endX; x-- {
					cloudsA[Coord{x, startY}]++
					cloudsB[Coord{x, startY}]++
				}
			}
		} else {
			if startX < endX {
				yIncr := 1
				if startY > endY {
					yIncr = -1
				}
				y := startY
				for x := startX; x <= endX; x++ {
					cloudsB[Coord{x, y}]++
					y += yIncr
				}
			} else {
				yIncr := 1
				if startY > endY {
					yIncr = -1
				}
				y := startY
				for x := startX; x >= endX; x-- {
					cloudsB[Coord{x, y}]++
					y += yIncr
				}
			}
		}
	}
	fmt.Println(countMultiples(cloudsA))
	fmt.Println(countMultiples(cloudsB))
}

func countMultiples(clouds map[Coord]int) int {
	var count int
	for _, v := range clouds {
		if v >= 2 {
			count++
		}
	}
	return count
}
