package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strings"
)

var inputFile = flag.String("inputFile", "inputs/day21.input", "Relative file path to use as input.")

type Coord struct {
	X, Y int
}

type Grid map[Coord]bool
type Distances map[Coord]int

const MaxSteps = 26501365
const TileSide = 131
const TileRadius = (MaxSteps - (TileSide / 2)) / TileSide

func main() {
	flag.Parse()
	bytes, err := ioutil.ReadFile(*inputFile)
	if err != nil {
		return
	}
	contents := string(bytes)
	split := strings.Split(contents, "\n")

	raw := make(Grid)
	var start, bounds Coord
	for y, line := range split[:len(split)-1] {
		for x, v := range line {
			loc := Coord{x, y}
			switch v {
			case '.':
				raw[loc] = true
			case '#':
				// do nothing. default map value is false.
				// raw[loc] = false
			case 'S':
				raw[loc] = true
				start = loc
			}
		}
		bounds.Y = y + 1
		bounds.X = len(line)
	}

	passable := make(Grid)
	for k, v := range raw {
		transformed := Coord{k.X - start.X, k.Y - start.Y}
		passable[transformed] = v
	}

	// Basic algorithm: we are supposed to go up to 26501365 steps.
	// Everything we've seen on a previous step can be ignored (modulo odd/even).
	// It takes us 65 steps to get to the side of our current square.
	// (and yes, there are no obstructions directly to sides or along edges)
	// Each square we cross in entirety then consumes 131 steps.
	// we can go straight up/down/left/right 202,300 tiles.
	// However, that last tile is not valid to assume we can reach all of; that
	// needs manual evaluation.
	// So, roughly, we will have a tile grid that ranges from -202300 to +202300
	// in each axis. If abs(x)+abs(y) > 202300, then discard it from
	// consideration. if abs(x)+abs(y) <= 202299, then we know every square within
	// can be reached and treat it as an even or odd tile.
	// Otherwise, if exactly 202300, we need to do the lookup to see how many
	// points within can be reached from its nearest corner to (0,0).
	// So we should start by memoising the number of points reachable with each
	// number of remaining moves from each corner and edge center.
	// There are no weird cases that force us to treat the set just inside the
	// edges as also suspect.
	var entrypoints [3][3]Distances
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			origin := Coord{(i - 1) * start.X, (j - 1) * start.Y}
			entrypoints[i][j] = passable.Distances(origin)
		}
	}

	// Part A
	fmt.Println(entrypoints[1][1].Count(64))

	if bounds.X != TileSide {
		// The example given in Part A does not work for Part B optimisation.
		return
	}

	// Part B
	scalingFactor := (TileRadius - 1) / 2

	var sum uint64
	// The interior region is made up of:
	// * One starting tile ("odd").
	// * N even and M odd whole tiles
	// The odds are the same as the original.
	odd := entrypoints[1][1].Count(MaxSteps)
	sum += odd * uint64(1+4*scalingFactor*(scalingFactor+1))
	even := entrypoints[1][1].Count(MaxSteps - 1)
	sum += even * uint64(4*(scalingFactor+1)*(scalingFactor+1))

	// * And all the diagonal tiles, of two kinds (small & large)
	for i := 0; i < 4; i++ {
		remaining := MaxSteps - 1 - TileSide*(TileRadius-1)
		big := entrypoints[2*(i%2)][2*(i/2)].Count(remaining)
		small := entrypoints[2*(i%2)][2*(i/2)].Count(remaining - TileSide)
		sum += TileRadius*small + (TileRadius-1)*big
	}

	// * And all of the axis tiles.
	sum += entrypoints[0][1].Count(TileSide - 1)
	sum += entrypoints[2][1].Count(TileSide - 1)
	sum += entrypoints[1][0].Count(TileSide - 1)
	sum += entrypoints[1][2].Count(TileSide - 1)

	fmt.Println(sum)
}

func (d Distances) Count(steps int) uint64 {
	var count uint64
	for _, v := range d {
		if v%2 == steps%2 && v <= steps {
			count++
		}
	}
	return count
}

func (g Grid) Distances(origin Coord) Distances {
	distances := Distances{
		origin: 0,
	}
	for {
		modified := false
		add := make(Distances)
		for loc, steps := range distances {
			for _, neigh := range []Coord{{loc.X + 1, loc.Y}, {loc.X - 1, loc.Y}, {loc.X, loc.Y + 1}, {loc.X, loc.Y - 1}} {
				if !g[neigh] {
					continue
				}
				if _, ok := distances[neigh]; !ok {
					add[neigh] = steps + 1
					modified = true
				}
			}
		}
		if !modified {
			break
		}
		for k, v := range add {
			distances[k] = v
		}
	}
	return distances
}
