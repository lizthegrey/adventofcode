package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"math"
	"strconv"
	"strings"
)

var inputFile = flag.String("inputFile", "inputs/day12.input", "Relative file path to use as input.")
var steps = flag.Int("steps", 1000, "Number of steps to simulate for.")

type Coord3 struct {
	X, Y, Z int32
}

type Moon struct {
	Pos, Vel *Coord3
}

func (m *Moon) TickVel(i int, ms Moons) int {
	// Avoid double counting or ticking against self.
	safe := int(math.MaxInt32)
	for j, o := range ms {
		if j >= i {
			break
		}
		xDelta := m.Pos.X - o.Pos.X
		yDelta := m.Pos.Y - o.Pos.Y
		zDelta := m.Pos.Z - o.Pos.Z
		xRVel := m.Vel.X - o.Vel.X
		yRVel := m.Vel.Y - o.Vel.Y
		zRVel := m.Vel.Z - o.Vel.Z
		if xDelta/xRVel < safe {
			safe = xDelta / xRVel
		}
		if yDelta/yRVel < safe {
			safe = yDelta / yRVel
		}
		if zDelta/yRVel < safe {
			safe = zDelta / yRVel
		}

		if m.Pos.X < o.Pos.X {
			m.Vel.X++
			o.Vel.X--
		} else if m.Pos.X > o.Pos.X {
			m.Vel.X--
			o.Vel.X++
		}
		if m.Pos.Y < o.Pos.Y {
			m.Vel.Y++
			o.Vel.Y--
		} else if m.Pos.Y > o.Pos.Y {
			m.Vel.Y--
			o.Vel.Y++
		}
		if m.Pos.Z < o.Pos.Z {
			m.Vel.Z++
			o.Vel.Z--
		} else if m.Pos.Z > o.Pos.Z {
			m.Vel.Z--
			o.Vel.Z++
		}
	}
	return int(safe)
}

func (m *Moon) TickPos() {
	m.Pos.X += m.Vel.X
	m.Pos.Y += m.Vel.Y
	m.Pos.Z += m.Vel.Z
}

func (m Moon) Kinetic() float64 {
	return math.Abs(float64(m.Vel.X)) + math.Abs(float64(m.Vel.Y)) + math.Abs(float64(m.Vel.Z))
}

func (m Moon) Potential() float64 {
	return math.Abs(float64(m.Pos.X)) + math.Abs(float64(m.Pos.Y)) + math.Abs(float64(m.Pos.Z))
}

type Moons []*Moon

func (ms Moons) TickOne() {
	for i, m := range ms {
		m.TickVel(i, ms)
	}
	for _, m := range ms {
		m.TickPos()
	}
}

func (ms Moons) TickMany() int {
	maxSafe := int(math.MaxInt32)
	for i, m := range ms {
		safe := m.TickVel(i, ms)
		if safe < maxSafe {
			maxSafe = safe
		}
	}
	for _, m := range ms {
		m.TickPos()
	}
	return 1
}

func main() {
	flag.Parse()
	bytes, err := ioutil.ReadFile(*inputFile)
	if err != nil {
		return
	}
	contents := string(bytes)
	split := strings.Split(contents, "\n")
	moons := make(Moons, 0)
	for _, s := range split {
		var pos, vel Coord3
		obj := Moon{&pos, &vel}
		if s == "" {
			continue
		}
		// <x=-4, y=3, z=15>
		coords := strings.Split(s[1:len(s)-1], ", ")
		if x, err := strconv.Atoi(coords[0][2:]); err != nil {
			fmt.Printf("Failed to parse line: %s\n", s)
		} else {
			obj.Pos.X = int32(x)
		}
		if y, err := strconv.Atoi(coords[1][2:]); err != nil {
			fmt.Printf("Failed to parse line: %s\n", s)
		} else {
			obj.Pos.Y = int32(y)
		}
		if z, err := strconv.Atoi(coords[2][2:]); err != nil {
			fmt.Printf("Failed to parse line: %s\n", s)
		} else {
			obj.Pos.Z = int32(z)
		}
		moons = append(moons, &obj)
	}

	moonCopy := make(Moons, len(moons))
	for i, m := range moons {
		pos := *m.Pos
		vel := *m.Vel
		moonCopy[i] = &Moon{&pos, &vel}
	}

	for i := 0; i < *steps; i++ {
		moonCopy.TickOne()
	}
	var energy float64
	for _, m := range moonCopy {
		energy += m.Kinetic() * m.Potential()
	}
	fmt.Println(int(energy))

	s := 0
	for {
		if s%10000000 == 0 {
			fmt.Printf("Iteration %d\n", s)
		}

		var zero Coord3
		zeroVel := true
		for _, m := range moons {
			if *m.Vel != zero {
				zeroVel = false
			}
		}
		if zeroVel && s != 0 {
			fmt.Printf("Back to 0 velocity and collapsing back on self at %d, solution is %d\n", s, 2*s)
			return
		}
		s += moons.TickMany()
	}
}
