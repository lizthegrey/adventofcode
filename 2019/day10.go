package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"math"
	"sort"
	"strings"
)

var inputFile = flag.String("inputFile", "inputs/day10.input", "Relative file path to use as input.")
var debug = flag.Bool("debug", false, "Whether to print debug output.")
var iterations = flag.Int("iterations", 200, "How many iterations to run.")

type Coord struct {
	Row, Col int
}

// -1,0 -> 0
// 0,1 -> pi/2 or 270 deg
// 1,0 -> pi
// 0,-1 -> 3pi/2
func (c Coord) Theta() float64 {
	return math.Atan2(float64(c.Col), -float64(c.Row))
}
func (c Coord) Distance(source Coord) float64 {
	return math.Abs(float64(source.Col-c.Col)) + math.Abs(float64(source.Row-c.Row))
}

type Grid map[Coord]bool

func main() {
	flag.Parse()
	bytes, err := ioutil.ReadFile(*inputFile)
	if err != nil {
		return
	}

	asteroids := make(Grid)

	contents := string(bytes)
	split := strings.Split(contents, "\n")
	for r, line := range split {
		for c, entry := range line {
			if entry == '#' {
				asteroids[Coord{r, c}] = true
			}
		}
	}

	mostVisible := 0
	var station Coord
	for coord := range asteroids {
		visible := 0
		offsets := make(Grid)
		for other := range asteroids {
			if other == coord {
				continue
			}
			rDelta := other.Row - coord.Row
			cDelta := other.Col - coord.Col
			offset := normalize(rDelta, cDelta)
			if offsets[offset] {
				continue
			}
			offsets[offset] = true
			visible++
		}
		if visible > mostVisible {
			mostVisible = visible
			station = coord
		}
	}
	fmt.Printf("%d asteroids visible from %d, %d\n", mostVisible, station.Col, station.Row)

	// Vaporize some asteroids.
	asteroidsByTheta := make(map[Coord][]Coord)
	thetas := make([]Coord, 0)

	for other := range asteroids {
		if other == station {
			continue
		}
		rDelta := other.Row - station.Row
		cDelta := other.Col - station.Col
		offset := normalize(rDelta, cDelta)
		asteroidsByTheta[offset] = append(asteroidsByTheta[offset], other)
	}
	for k := range asteroidsByTheta {
		// Also sort by distance.
		sort.Slice(asteroidsByTheta[k], func(i, j int) bool {
			return asteroidsByTheta[k][i].Distance(station) < asteroidsByTheta[k][j].Distance(station)
		})
		thetas = append(thetas, k)
	}
	sort.Slice(thetas, func(i, j int) bool {
		return thetas[i].Theta() < thetas[j].Theta()
	})
	var thetaIndex int
	for i, v := range thetas {
		up := Coord{-1, 0}
		if v == up {
			thetaIndex = i
			break
		}
	}
	for i := 1; i <= *iterations && i < len(asteroids); i++ {
		for {
			if thetaIndex >= len(thetas) {
				thetaIndex = 0
			}
			if len(asteroidsByTheta[thetas[thetaIndex]]) > 0 {
				break
			}
			thetaIndex++
		}
		destroyed := asteroidsByTheta[thetas[thetaIndex]][0]
		asteroidsByTheta[thetas[thetaIndex]] = asteroidsByTheta[thetas[thetaIndex]][1:]
		if *debug || i == *iterations || i+1 == len(asteroids) {
			fmt.Printf("#%d: firing at theta %f, destroyed %d,%d\n", i, thetas[thetaIndex].Theta(), destroyed.Col, destroyed.Row)
		}
		thetaIndex++
	}
}

// Normalizes 6,3 to 2,1; normalizes -10,-11 to -10,-11; normalizes 8,6 to 4,3
func normalize(r, c int) Coord {
	if c == 0 {
		return Coord{r / int(math.Abs(float64(r))), 0}
	} else if r == 0 {
		return Coord{0, c / int(math.Abs(float64(c)))}
	}

	var greater int
	var lesser int
	if r > c {
		greater = int(math.Abs(float64(r)))
		lesser = int(math.Abs(float64(c)))
	} else if r <= c {
		greater = int(math.Abs(float64(c)))
		lesser = int(math.Abs(float64(r)))
	}

	gcd := 1
	for {
		remainder := greater % lesser
		if remainder == 0 {
			gcd = lesser
			break
		}
		greater = lesser
		lesser = remainder
	}

	return Coord{r / gcd, c / gcd}
}
