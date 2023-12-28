package main

import (
	"cmp"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"slices"
	"strconv"
	"strings"
)

var inputFile = flag.String("inputFile", "inputs/day22.input", "Relative file path to use as input.")

type Coord struct {
	X, Y, Z int
}

type Brick struct {
	Lo Coord
	Hi Coord
}

type Pieces []Brick

func main() {
	flag.Parse()
	bytes, err := ioutil.ReadFile(*inputFile)
	if err != nil {
		return
	}
	contents := string(bytes)
	split := strings.Split(contents, "\n")

	var pieces Pieces
	for _, s := range split[:len(split)-1] {
		parts := strings.Split(s, "~")
		lo := toCoord(parts[0])
		hi := toCoord(parts[1])
		pieces = append(pieces, Brick{lo, hi})
	}

	// Now settle the pieces, throwing away the return value.
	pieces.Settle()

	// Now try to disintegrate each brick and see what happens.
	var safe, cascades int
	for i := range pieces {
		// Copy the things because it's a destructive procedure.
		cpy := make(Pieces, 0, len(pieces)-1)
		for j, v := range pieces {
			if i == j {
				continue
			}
			cpy = append(cpy, v)
		}
		if result := cpy.Settle(); result == 0 {
			safe++
		} else {
			cascades += result
		}
	}
	fmt.Println(safe)
	fmt.Println(cascades)
}

func (pieces Pieces) Index() [10][10]map[*Brick]*Brick {
	// We make the observation that X and Y are only 10x10 (0-9 each)
	// and thus we can just make a set of vertically arranged stacks
	// for each x and y coord.
	var arrs [10][10][]*Brick
	for i := range pieces {
		piece := &pieces[i]
		for x := piece.Lo.X; x <= piece.Hi.X; x++ {
			for y := piece.Lo.Y; y <= piece.Hi.Y; y++ {
				arrs[x][y] = append(arrs[x][y], piece)
			}
		}
	}
	var indices [10][10]map[*Brick]*Brick
	for x := 0; x <= 9; x++ {
		for y := 0; y <= 9; y++ {
			slices.SortFunc(arrs[x][y], func(a, b *Brick) int {
				return cmp.Compare(a.Lo.Z, b.Lo.Z)
			})
			indices[x][y] = make(map[*Brick]*Brick)
			for i := range arrs[x][y] {
				if i == 0 {
					continue
				}
				indices[x][y][arrs[x][y][i]] = arrs[x][y][i-1]
			}
		}
	}
	return indices
}

func (pieces Pieces) Settle() int {
	indices := pieces.Index()
	changed := make(map[*Brick]bool)
	for {
		var modified bool
		for i := range pieces {
			piece := &pieces[i]
			var minZ int
			for x := piece.Lo.X; x <= piece.Hi.X; x++ {
				for y := piece.Lo.Y; y <= piece.Hi.Y; y++ {
					limit := 1
					if below := indices[x][y][piece]; below != nil {
						limit = below.Hi.Z + 1
					}
					if limit > minZ {
						minZ = limit
					}
				}
			}
			diff := piece.Lo.Z - minZ
			if diff > 0 {
				piece.Lo.Z -= diff
				piece.Hi.Z -= diff
				modified = true
				changed[piece] = true
			}
		}
		if !modified {
			break
		}
	}
	return len(changed)
}

func toCoord(in string) Coord {
	parts := strings.Split(in, ",")
	var ret Coord
	var err error
	ret.X, err = strconv.Atoi(parts[0])
	if err != nil {
		log.Fatalf("Failed to parse %s: %v", in, err)
	}
	ret.Y, err = strconv.Atoi(parts[1])
	if err != nil {
		log.Fatalf("Failed to parse %s: %v", in, err)
	}
	ret.Z, err = strconv.Atoi(parts[2])
	if err != nil {
		log.Fatalf("Failed to parse %s: %v", in, err)
	}
	return ret
}
