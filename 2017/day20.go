package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"math"
	"regexp"
	"strconv"
	"strings"
)

var inputFile = flag.String("inputFile", "inputs/day20.input", "Relative file path to use as input.")

type Coord struct {
	X, Y, Z int
}

type Particle struct {
	P, V, A Coord
}

type Board []*Particle

func main() {
	flag.Parse()

	bytes, err := ioutil.ReadFile(*inputFile)
	if err != nil {
		fmt.Printf("Could not open file %s because %v.\n", *inputFile, err)
	}
	contents := string(bytes[:len(bytes)-1])
	lines := strings.Split(contents, "\n")

	particles := make(Board, len(lines))
	re := regexp.MustCompile("p=<([0-9-]+),([0-9-]+),([0-9-]+)>, v=<([0-9-]+),([0-9-]+),([0-9-]+)>, a=<([0-9-]+),([0-9-]+),([0-9-]+)>")
	for i, l := range lines {
		s := re.FindStringSubmatch(l)
		if len(s) != 10 {
			fmt.Printf("Failed to parse line '%s'\n", l)
			return
		}
		r := make([]int, len(s)-1)
		for i := 1; i < len(s); i++ {
			v, _ := strconv.Atoi(s[i])
			r[i-1] = v
		}
		p := Particle{Coord{r[0], r[1], r[2]}, Coord{r[3], r[4], r[5]}, Coord{r[6], r[7], r[8]}}
		particles[i] = &p
	}

	smallestAccel := math.MaxInt32
	var smallestAccelId int
	for i, p := range particles {
		a := p.A.Distance()
		if a < smallestAccel {
			smallestAccel = a
			smallestAccelId = i
		}
	}
	fmt.Printf("Particle with smallest acceleration has index %d\n", smallestAccelId)
}

func (b Board) Tick() {
	for _, p := range b {
		p.V = p.V.Add(p.A)
		p.P = p.P.Add(p.V)
	}
}

func (c Coord) Distance() int {
	r := 0
	r += int(math.Abs(float64(c.X)))
	r += int(math.Abs(float64(c.Y)))
	r += int(math.Abs(float64(c.Z)))
	return r
}

func (a Coord) Add(b Coord) Coord {
	r := a
	r.X += b.X
	r.Y += b.Y
	r.Z += b.Z
	return r
}
