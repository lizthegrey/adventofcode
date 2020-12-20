package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

var inputFile = flag.String("inputFile", "inputs/day20.input", "Relative file path to use as input.")
var debug = flag.Bool("debug", false, "Whether to print debug output along the way.")

const SeaMonster string = `                  # 
#    ##    ##    ###
 #  #  #  #  #  #   `

type Tile [10][10]bool
type Edge uint16
type Cropped [8][8]bool

type RotatedPiece struct {
	TileID   int
	Flipped  bool
	Rotation uint8
}
type Mosaic [][]RotatedPiece

func (t Tile) Crop() Cropped {
	var ret Cropped
	for r := 1; r < 9; r++ {
		for c := 1; c < 9; c++ {
			ret[r-1][c-1] = t[r][c]
		}
	}
	return ret
}

func (t Tile) Edges() [4]Edge {
	// left to right
	var top, bottom Edge
	for i, v := range t[0] {
		if v {
			top |= 1 << i
		}
	}
	for i, v := range t[9] {
		if v {
			bottom |= 1 << i
		}
	}
	// top to bottom
	var left, right Edge
	for i, r := range t {
		if r[0] {
			left |= 1 << i
		}
		if r[9] {
			right |= 1 << i
		}
	}
	return [4]Edge{top, bottom, left, right}
}

func (e Edge) Flip() Edge {
	var ret Edge
	for i := 0; i < 10; i++ {
		ret |= ((e >> i) & 1) << (9 - i)
	}
	return ret
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
	tiles := make(map[int]Tile)
	for i := 0; i < len(split); i += 12 {
		// First line of metadata describes the tile id.
		// Tile 2477:
		header := split[i]
		n, err := strconv.Atoi(header[5 : len(header)-1])
		if err != nil {
			fmt.Printf("Failed to parse %s\n", header)
			break
		}
		var tile Tile
		for r := 0; r < 10; r++ {
			line := split[i+1+r]
			for c := 0; c < 10; c++ {
				if line[c] == '#' {
					tile[r][c] = true
				}
			}
		}
		tiles[n] = tile
	}
	if *debug {
		for i, t := range tiles {
			edges := t.Edges()
			var flipped [4]Edge
			for i, e := range edges {
				flipped[i] = e.Flip()
			}
			fmt.Printf("Tile %d: %v (flipped: %v)\n", i, edges, flipped)
		}
	}
	allEdges := make(map[Edge][]int)
	for i, t := range tiles {
		edges := t.Edges()
		for _, e := range edges {
			f := e.Flip()
			allEdges[e] = append(allEdges[e], i)
			allEdges[f] = append(allEdges[f], i)
		}
	}
	singletonTiles := make(map[int]int)
	for _, v := range allEdges {
		if len(v) < 2 {
			singletonTiles[v[0]]++
		}
	}

	var corners []int
	sides := make(map[int]bool)

	product := 1
	for k, v := range singletonTiles {
		if v <= 2 {
			sides[k] = true
		}
		if v > 2 && v < 5 {
			corners = append(corners, k)
			product *= k
		}
	}
	fmt.Println(product)
	if *debug {
		fmt.Printf("Number of corners found: %d; number of sides found: %d\n", len(corners), len(sides))
	}

	var image Mosaic
	for r := 0; r*r < len(tiles); r++ {
		var line []RotatedPiece
		for c := 0; c*c < len(tiles); c++ {
			line = append(line, RotatedPiece{})
		}
		image = append(image, line)
	}
	image[0][0].TileID = corners[0]

	// Fill in the top row.
	for c := 1; (c+1)*(c+1) < len(tiles); c++ {
	}
}
