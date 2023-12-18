package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
)

var inputFile = flag.String("inputFile", "inputs/day18.input", "Relative file path to use as input.")

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

	var loc Coord
	var maxX, maxY, minX, minY int
	edges := make(map[Coord]string)
	edges[loc] = "start"
	for _, s := range split[:len(split)-1] {
		parts := strings.Split(s, " ")
		var incr Coord
		switch parts[0] {
		case "U":
			incr.Y++
		case "R":
			incr.X++
		case "D":
			incr.Y--
		case "L":
			incr.X--
		}
		distance, err := strconv.Atoi(parts[1])
		if err != nil {
			log.Fatalf("Failed parsing %s: %v", parts[1], err)
		}
		for i := 0; i < distance; i++ {
			loc.X += incr.X
			loc.Y += incr.Y
			edges[loc] = parts[2]
		}
		if loc.X < minX {
			minX = loc.X
		}
		if loc.Y < minY {
			minY = loc.Y
		}
		if loc.X > maxX {
			maxX = loc.X
		}
		if loc.Y > maxY {
			maxY = loc.Y
		}
	}

	fill := make(map[Coord]bool)
	var q []Coord
	var seenEdge bool
	for x := minX; x <= maxX; x++ {
		loc := Coord{x, maxY - 1}
		if seenEdge && edges[loc] == "" {
			q = append(q, loc)
			break
		}
		if edges[loc] != "" {
			seenEdge = true
		}
	}
	for len(q) > 0 {
		cur := q[0]
		q = q[1:]
		if fill[cur] || edges[cur] != "" {
			continue
		}
		fill[cur] = true
		for _, n := range []Coord{{cur.X - 1, cur.Y}, {cur.X + 1, cur.Y}, {cur.X, cur.Y - 1}, {cur.X, cur.Y + 1}} {
			if edges[n] != "" {
				continue
			}
			if cur.X < minX || cur.X > maxX || cur.Y < minY || cur.Y > maxY {
				continue
			}
			q = append(q, n)
		}
	}
	for y := maxY; y >= minY; y-- {
		for x := minX; x <= maxX; x++ {
			loc := Coord{x, y}
			if fill[loc] {
				fmt.Printf("F")
			} else if edges[loc] != "" {
				fmt.Printf("#")
			} else {
				fmt.Printf(".")
			}
		}
		fmt.Println()
	}
	fmt.Println(len(fill) + len(edges))
}
