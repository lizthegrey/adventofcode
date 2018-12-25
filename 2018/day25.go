package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var inputFile = flag.String("inputFile", "inputs/day25.input", "Relative file path to use as input.")

type Coord4 struct {
	X, Y, Z, T int
}

func (c Coord4) Adjacent(t Coord4) bool {
	deltaX := c.X - t.X
	deltaY := c.Y - t.Y
	deltaZ := c.Z - t.Z
	deltaT := c.T - t.T

	sum := 0
	if deltaX < 0 {
		sum -= deltaX
	} else {
		sum += deltaX
	}
	if deltaY < 0 {
		sum -= deltaY
	} else {
		sum += deltaY
	}
	if deltaZ < 0 {
		sum -= deltaZ
	} else {
		sum += deltaZ
	}
	if deltaT < 0 {
		sum -= deltaT
	} else {
		sum += deltaT
	}

	return sum <= 3
}

func main() {
	flag.Parse()
	f, err := os.Open(*inputFile)
	if err != nil {
		return
	}
	defer f.Close()

	coords := make([]Coord4, 0)
	r := bufio.NewReader(f)
	for {
		l, err := r.ReadString('\n')
		if err != nil || len(l) == 0 {
			break
		}
		l = l[:len(l)-1]
		nums := strings.Split(l, ",")
		x, _ := strconv.Atoi(nums[0])
		y, _ := strconv.Atoi(nums[1])
		z, _ := strconv.Atoi(nums[2])
		t, _ := strconv.Atoi(nums[3])
		coords = append(coords, Coord4{x, y, z, t})
	}

	// Maps index of coord to clique number.
	cliquesByCoord := make(map[int]int)
	coordsByClique := make(map[int][]int)

	// Seed the first element.
	highestClique := 1

	changed := true
	for changed {
		changed = false
		for c := range coords {
			cClique := cliquesByCoord[c]
			if cClique == 0 {
				cliquesByCoord[c] = highestClique
				coordsByClique[highestClique] = []int{c}
				cClique = highestClique
				highestClique++
				changed = true
			}
			for t := range coords {
				tClique := cliquesByCoord[t]
				if tClique == cClique {
					continue
				}
				if coords[c].Adjacent(coords[t]) {
					if tClique == 0 {
						cliquesByCoord[t] = cClique
						coordsByClique[cClique] = append(coordsByClique[cClique], t)
						changed = true
						continue
					}
					// Merge the t Clique and the c Clique.
					changed = true
					for _, v := range coordsByClique[tClique] {
						cliquesByCoord[v] = cClique
					}
					coordsByClique[cClique] = append(coordsByClique[cClique], coordsByClique[tClique]...)
					delete(coordsByClique, tClique)
				}
			}
		}
	}
	fmt.Println(len(coordsByClique))
}
