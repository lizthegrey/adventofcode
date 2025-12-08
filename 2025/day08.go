package main

import (
	"cmp"
	"flag"
	"fmt"
	"io/ioutil"
	"slices"
	"strconv"
	"strings"
)

var inputFile = flag.String("inputFile", "inputs/day08.input", "Relative file path to use as input.")

type Coord3 struct {
	X, Y, Z int
}

type Pair struct {
	Lower, Upper Coord3
	Distance     uint64
}

func main() {
	flag.Parse()
	bytes, err := ioutil.ReadFile(*inputFile)
	if err != nil {
		return
	}
	contents := string(bytes)
	split := strings.Split(contents, "\n")

	var coords []Coord3
	var distances []Pair
	for _, s := range split[:len(split)-1] {
		parts := strings.Split(s, ",")
		x, _ := strconv.Atoi(parts[0])
		y, _ := strconv.Atoi(parts[1])
		z, _ := strconv.Atoi(parts[2])
		coords = append(coords, Coord3{x, y, z})
	}

	for i, a := range coords {
		for j, b := range coords {
			if i <= j {
				continue
			}
			distX := a.X - b.X
			distY := a.Y - b.Y
			distZ := a.Z - b.Z
			distances = append(distances, Pair{
				Lower:    a,
				Upper:    b,
				Distance: uint64(distX*distX) + uint64(distY*distY) + uint64(distZ*distZ),
			})
		}
	}
	slices.SortFunc(distances, func(a, b Pair) int {
		return cmp.Compare(a.Distance, b.Distance)
	})

	circuits := make(map[Coord3]int)
	var circuitId int
	stepsPartA := 10
	if len(coords) > 50 {
		stepsPartA = 1000
	}
	for i, toJoin := range distances {
		l, foundL := circuits[toJoin.Lower]
		u, foundU := circuits[toJoin.Upper]
		if foundL && foundU {
			for k, v := range circuits {
				if v == u {
					circuits[k] = l
				}
			}
		} else if foundL {
			circuits[toJoin.Upper] = l
		} else if foundU {
			circuits[toJoin.Lower] = u
		} else {
			circuits[toJoin.Lower] = circuitId
			circuits[toJoin.Upper] = circuitId
			circuitId++
		}

		sizes := make(map[int]int)
		for _, l := range circuits {
			sizes[l]++
		}
		var sizeSlice []int
		for _, v := range sizes {
			sizeSlice = append(sizeSlice, v)
		}
		if len(sizeSlice) == 1 && sizeSlice[0] == len(coords) {
			fmt.Println(toJoin.Upper.X * toJoin.Lower.X)
			break
		}
		if i+1 == stepsPartA {
			slices.Sort(sizeSlice)
			slices.Reverse(sizeSlice)
			fmt.Println(sizeSlice[0] * sizeSlice[1] * sizeSlice[2])
		}
	}
}
