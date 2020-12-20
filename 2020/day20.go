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

var SeaMonster = strings.Split(`                  # 
#    ##    ##    ###
 #  #  #  #  #  #   `, " ")

const (
	TOP = iota
	RIGHT
	BOTTOM
	LEFT
)

type Tile [10][10]bool
type Edge uint16
type Cropped [8][8]bool

type Coord struct {
	Row, Col int
}

type RotatedPiece struct {
	TileID   int
	Flipped  bool
	Rotation int
}
type Mosaic [][]RotatedPiece

func (m Mosaic) PixelAtCoord(tiles map[int]Tile, r, c int) bool {
	mosaicRow := r / 8
	mosaicCol := c / 8
	t := tiles[m[mosaicRow][mosaicCol].TileID]
	return t.Crop()[r%8][c%8]
}

func (m Mosaic) MonsterTopLeftCoord(tiles map[int]Tile, r, c int) bool {
	for rOffset, row := range SeaMonster {
		for cOffset, char := range row {
			if char != '#' {
				continue
			}
			if !m.PixelAtCoord(tiles, r+rOffset, c+cOffset) {
				return false
			}
		}
	}
	return true
}

func (m Mosaic) FindMonsters(tiles map[int]Tile) []Coord {
	// Find the sea monster.
	var monsters []Coord
	for r := 0; (r+len(SeaMonster))*(r+len(SeaMonster)) < 100*len(tiles); r++ {
		for c := 0; (c+len(SeaMonster[0]))*(c+len(SeaMonster[0])) < 100*len(tiles); c++ {
			if m.MonsterTopLeftCoord(tiles, r, c) {
				monsters = append(monsters, Coord{r, c})
			}
		}
	}
	return monsters
}

func (r RotatedPiece) GetEdge(tiles map[int]Tile, side int) Edge {
	t := tiles[r.TileID]
	edges := t.Edges()
	idx := (side + 4 - r.Rotation) % 4
	if r.Flipped && (idx == RIGHT || idx == LEFT) {
		idx = (idx + 2) % 4
	}
	edge := edges[idx]
	if r.Flipped {
		edge = edge.Flip()
	}
	return edge
}

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
	// Return in clockwise orientation.
	return [4]Edge{top, right, bottom.Flip(), left.Flip()}
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

	// Fill in the top row.
	var topLeft int
	for i := 0; i < 4; i++ {
		matchingRight := allEdges[tiles[corners[i]].Edges()[RIGHT]]
		matchingDown := allEdges[tiles[corners[i]].Edges()[BOTTOM]]
		if len(matchingRight) < 2 || len(matchingDown) < 2 {
			// This is not the top left corner.
			continue
		}
		topLeft = corners[i]
	}

	image[0][0].TileID = topLeft
	used := map[int]bool{topLeft: true}
	count := image.Traverse(tiles, used, 0, 0, RIGHT)
	if count == 0 {
		fmt.Println("Failed to traverse right from what should be top left.")
		return
	}

	for c := 0; c*c < len(tiles); c++ {
		count = image.Traverse(tiles, used, 0, c, BOTTOM)
		if count == 0 {
			fmt.Printf("Failed to traverse down column %d\n", c)
			return
		}
	}
	fmt.Println(image)

	monsters := image.FindMonsters(tiles)
	fmt.Println(len(monsters))
}

func (m Mosaic) Traverse(tiles map[int]Tile, used map[int]bool, r, c int, dir int) int {
	loose := m[r][c].GetEdge(tiles, dir)
	var rIncr, cIncr int
	switch dir {
	case RIGHT:
		cIncr = 1
	case BOTTOM:
		rIncr = 1
	case LEFT:
		cIncr = -1
	case TOP:
		rIncr = -1
	}
	i := 1
outer:
	for ; ; i++ {
		row := r + rIncr*i
		col := c + cIncr*i
		if row < 0 || row*row >= len(tiles) {
			break
		}
		if row < 0 || col*col >= len(tiles) {
			break
		}
		for s, t := range tiles {
			if used[s] {
				// Don't re-use the same piece twice.
				continue
			}
			edges := t.Edges()
			for d, e := range edges {
				if loose == e {
					m[row][col].TileID = s
					m[row][col].Flipped = !m[row-rIncr][col-cIncr].Flipped
					m[row][col].Rotation = (4 + dir - d) % 4
					loose = edges[(d+2)%4]
					used[s] = true
					continue outer
				}
				if loose.Flip() == e {
					m[row][col].TileID = s
					m[row][col].Flipped = m[row-rIncr][col-cIncr].Flipped
					m[row][col].Rotation = (4 + dir - d) % 4
					loose = edges[(d+2)%4].Flip()
					used[s] = true
					continue outer
				}
				// Otherwise, this edge didn't match. Continue on to other edges.
			}
			// This piece hasn't matched, continue on to other pieces.
		}
		// We didn't match any pieces, abort.
		return 0
	}
	return i - 1
}
