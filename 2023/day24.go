package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
)

var inputFile = flag.String("inputFile", "inputs/day24.input", "Relative file path to use as input.")
var lo = flag.Float64("lo", 200000000000000, "Coordinate lower bounds to use.")
var hi = flag.Float64("hi", 400000000000000, "Coordinate upper bounds to use.")

type Coord struct {
	X, Y, Z int64
}

type Particle struct {
	Pos, Vel Coord
}

func main() {
	flag.Parse()
	bytes, err := ioutil.ReadFile(*inputFile)
	if err != nil {
		return
	}
	contents := string(bytes)
	split := strings.Split(contents, "\n")

	var particles []Particle
	for _, s := range split[:len(split)-1] {
		parts := strings.Split(s, " @ ")
		p := Particle{toCoord(parts[0]), toCoord(parts[1])}
		particles = append(particles, p)
	}

	var partA int
	for i, a := range particles {
		for j, b := range particles {
			if i >= j {
				// Only handle each pair once.
				continue
			}
			if intersect2d(a, b) {
				partA++
			}
		}
	}
	fmt.Println(partA)

	// And then for part B you just give up and use a CAS/linear equation solver.
	// I don't have it in me to write a CAS/linear equation solver.
	a := particles[0]
	b := particles[1]
	c := particles[2]
	fmt.Printf("https://www.wolframalpha.com/input?i2d=true&i=Solve%%5C%%2891%%29%%7B%%7Bx%%2Cy%%2Cz%%7D%%2B%%7BSubscript%%5Bv%%2C1%%5D%%2C+Subscript%%5Bv%%2C2%%5D%%2C+Subscript%%5Bv%%2C3%%5D%%7D*Subscript%%5Bt%%2C0%%5D%%3D%%3D%%7B%d%%2C%d%%2C%d%%7D%%2B%%7B%d%%2C%d%%2C%d%%7D*Subscript%%5Bt%%2C0%%5D%%2C+%%7Bx%%2Cy%%2Cz%%7D%%2B%%7BSubscript%%5Bv%%2C1%%5D%%2C+Subscript%%5Bv%%2C2%%5D%%2C+Subscript%%5Bv%%2C3%%5D%%7D*Subscript%%5Bt%%2C1%%5D%%3D%%3D%%7B%d%%2C%d%%2C%d%%7D%%2B%%7B%d%%2C%d%%2C%d%%7D*Subscript%%5Bt%%2C1%%5D%%2C+%%7Bx%%2Cy%%2Cz%%7D%%2B%%7BSubscript%%5Bv%%2C1%%5D%%2C+Subscript%%5Bv%%2C2%%5D%%2C+Subscript%%5Bv%%2C3%%5D%%7D*Subscript%%5Bt%%2C2%%5D%%3D%%3D%%7B%d%%2C%d%%2C%d%%7D%%2B%%7B%d%%2C%d%%2C%d%%7D*Subscript%%5Bt%%2C2%%5D%%7D%%5C%%2844%%29+x%%5C%%2844%%29+y%%5C%%2844%%29+z%%5C%%2893%%29\n",
		a.Pos.X, a.Pos.Y, a.Pos.Z, a.Vel.X, a.Vel.Y, a.Vel.Z,
		b.Pos.X, b.Pos.Y, b.Pos.Z, b.Vel.X, b.Vel.Y, b.Vel.Z,
		c.Pos.X, c.Pos.Y, c.Pos.Z, c.Vel.X, c.Vel.Y, c.Vel.Z,
	)
}

func intersect2d(a, b Particle) bool {
	// a.Pos+t_a*a.Vel == b.Pos+t_b*b.Vel for some t_a, t_b
	// Mathematica says the answer is:
	denominator := a.Vel.X*b.Vel.Y - a.Vel.Y*b.Vel.X
	if denominator == 0 {
		// Lines are parallel, avoid divide by 0.
		return false
	}
	pDiff := a.Pos.Sub(b.Pos)
	t1 := float64((b.Vel.X * pDiff.Y) - (b.Vel.Y * pDiff.X))
	t2 := float64((a.Vel.X * pDiff.Y) - (a.Vel.Y * pDiff.X))
	t1 /= float64(denominator)
	t2 /= float64(denominator)
	if t1 < 0 || t2 < 0 {
		return false
	}
	if x := float64(a.Vel.X)*t1 + float64(a.Pos.X); x < *lo || x > *hi {
		return false
	}
	if y := float64(a.Vel.Y)*t1 + float64(a.Pos.Y); y < *lo || y > *hi {
		return false
	}
	return true
}

func (c Coord) Sub(o Coord) Coord {
	return Coord{c.X - o.X, c.Y - o.Y, c.Z - o.Z}
}

func (c Coord) Add(o Coord) Coord {
	return Coord{c.X + o.X, c.Y + o.Y, c.Z + o.Z}
}

func (c Coord) Mul(o int64) Coord {
	return Coord{c.X * o, c.Y * o, c.Z * o}
}

func toCoord(in string) Coord {
	parts := strings.Split(in, ", ")
	var ret Coord
	var err error
	ret.X, err = strconv.ParseInt(parts[0], 10, 0)
	if err != nil {
		log.Fatalf("Failed to parse %%s: %%v", in, err)
	}
	ret.Y, err = strconv.ParseInt(parts[1], 10, 0)
	if err != nil {
		log.Fatalf("Failed to parse %%s: %%v", in, err)
	}
	ret.Z, err = strconv.ParseInt(parts[2], 10, 0)
	if err != nil {
		log.Fatalf("Failed to parse %%s: %%v", in, err)
	}
	return ret
}
