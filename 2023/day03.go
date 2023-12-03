package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strings"
)

var inputFile = flag.String("inputFile", "inputs/day03.input", "Relative file path to use as input.")

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

	sum := 0
	grid := make(map[Coord]rune)
	gears := make(map[Coord][]int)
	var maxX, maxY int
	for y, row := range split[:len(split)-1] {
		for x, c := range row {
			grid[Coord{x, y}] = c
			if x > maxX {
				maxX = x
			}
		}
		if y > maxY {
			maxY = y
		}
	}
	for y := 0; y <= maxY; y++ {
		for x := 0; x <= maxX; {
			var length, num int
			for {
				c := grid[Coord{x + length, y}]
				if c < '0' || c > '9' {
					break
				}
				num = 10*num + int(c-'0')
				length++
			}
			if length == 0 {
				x++
				continue
			}
			var added bool
			for j := y - 1; j <= y+1; j++ {
				for i := x - 1; i <= x+length; i++ {
					nCoord := Coord{i, j}
					neigh := grid[nCoord]
					if neigh == 0 || neigh == '.' {
						continue
					}
					if neigh >= '0' && neigh <= '9' {
						continue
					}
					if !added {
						sum += num
						added = true
					}
					if neigh == '*' {
						gears[nCoord] = append(gears[nCoord], num)
					}
				}
			}
			x += length
		}
	}
	fmt.Println(sum)

	ratio := 0
	for _, v := range gears {
		if len(v) == 2 {
			ratio += v[0] * v[1]
		}
	}
	fmt.Println(ratio)
}
