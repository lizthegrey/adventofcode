package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

var inputFile = flag.String("inputFile", "inputs/day18.input", "Relative file path to use as input.")

type coord struct {
	x, y, z int
}

func (c coord) neighbours() []coord {
	return []coord{
		{c.x + 1, c.y + 0, c.z + 0},
		{c.x - 1, c.y + 0, c.z + 0},
		{c.x + 0, c.y + 1, c.z + 0},
		{c.x + 0, c.y - 1, c.z + 0},
		{c.x + 0, c.y + 0, c.z + 1},
		{c.x + 0, c.y + 0, c.z - 1},
	}
}

func (c coord) inBounds(min, max int) bool {
	return c.x >= min && c.x <= max && c.y >= min && c.y <= max && c.z >= min && c.z <= max
}

type model map[coord]bool

func main() {
	flag.Parse()
	bytes, err := ioutil.ReadFile(*inputFile)
	if err != nil {
		return
	}
	contents := string(bytes)
	split := strings.Split(contents, "\n")

	droplet := make(model)
	for _, s := range split[:len(split)-1] {
		parts := strings.Split(s, ",")
		x, _ := strconv.Atoi(parts[0])
		y, _ := strconv.Atoi(parts[1])
		z, _ := strconv.Atoi(parts[2])
		droplet[coord{x, y, z}] = true
	}

	// part A
	var totalArea int
	for k := range droplet {
		for _, n := range k.neighbours() {
			if !droplet[n] {
				totalArea++
			}
		}
	}
	fmt.Println(totalArea)

	// part B
	// Grow a void starting from the walls inward.
	// Looking at the input, it ranges from [0,19] but not negative numbers or numbers above 20.
	// We can start from any points outside the sculpture and grow from there.
	seed := coord{-1, -1, -1}
	void := make(model)
	q := []coord{seed}
	for len(q) != 0 {
		head := q[0]
		q = q[1:]
		if void[head] {
			continue
		}
		void[head] = true
		for _, n := range head.neighbours() {
			if !void[n] && !droplet[n] && n.inBounds(-1, 20) {
				q = append(q, n)
			}
		}
	}
	var exteriorArea int
	for k := range void {
		for _, n := range k.neighbours() {
			if droplet[n] {
				exteriorArea++
			}
		}
	}
	fmt.Println(exteriorArea)
}
