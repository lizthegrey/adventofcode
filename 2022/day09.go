package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

var inputFile = flag.String("inputFile", "inputs/day09.input", "Relative file path to use as input.")

type coord struct {
	x, y int
}

func (c coord) dragged(o coord) coord {
	diffX := o.x - c.x
	diffY := o.y - c.y
	if diffX < 2 && diffX > -2 && diffY < 2 && diffY > -2 {
		// Within +/- 1 in each dimension, stay in place.
		return c
	}
	// We will always move 1 in the direction that we were off by 2.
	// and we will also move 1 in the other direction if there is a diff.
	if diffX > 0 {
		c.x += (diffX + 1) / 2
	} else {
		c.x += (diffX - 1) / 2
	}
	if diffY > 0 {
		c.y += (diffY + 1) / 2
	} else {
		c.y += (diffY - 1) / 2
	}
	return c
}

func main() {
	flag.Parse()
	bytes, err := ioutil.ReadFile(*inputFile)
	if err != nil {
		return
	}
	contents := string(bytes)
	split := strings.Split(contents, "\n")

	// part A
	visited := make(map[coord]bool)
	var head, tail coord
	visited[tail] = true
	for _, s := range split[:len(split)-1] {
		c, _ := strconv.Atoi(s[2:])
		for i := 0; i < c; i++ {
			switch s[0] {
			case 'U':
				head.y++
			case 'D':
				head.y--
			case 'L':
				head.x--
			case 'R':
				head.x++
			}
			tail = tail.dragged(head)
			visited[tail] = true
		}
	}
	fmt.Println(len(visited))

	// part B
	visited = make(map[coord]bool)
	var positions [10]coord
	visited[positions[len(positions)-1]] = true
	for _, s := range split[:len(split)-1] {
		c, _ := strconv.Atoi(s[2:])
		for i := 0; i < c; i++ {
			switch s[0] {
			case 'U':
				positions[0].y++
			case 'D':
				positions[0].y--
			case 'L':
				positions[0].x--
			case 'R':
				positions[0].x++
			}
			for j := 1; j < len(positions); j++ {
				positions[j] = positions[j].dragged(positions[j-1])
			}
			visited[positions[len(positions)-1]] = true
		}
	}
	fmt.Println(len(visited))
}
