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
			c := grid[Coord{x, y}]
			if c < '0' || c > '9' {
				x++
				continue
			}
			// We've found a number. Keep reading to the right.
			num := int(c - '0')
			length := 1
			for i := 1; ; i++ {
				n := grid[Coord{x + i, y}]
				if n < '0' || n > '9' {
					break
				}
				num = 10*num + int(n-'0')
				length++
			}
		outer:
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
					sum += num
					if neigh == '*' {
						gears[nCoord] = append(gears[nCoord], num)
					}
					break outer
				}
			}
			x = x + length
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
