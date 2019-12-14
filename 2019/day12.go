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
var steps = flag.Int("steps", 1000, "Number of steps to simulate before reporting energy.")

type Axis struct {
	Pos, Vel int32
}

type Moon [3]*Axis

func (m *Moon) TickVel(i int, ms Moons) {
	// Avoid double counting or ticking against self.
	for j, o := range ms {
		if j >= i {
			break
		}
		for k, a := range m {
			oa := o[k]
			if a.Pos < oa.Pos {
				a.Vel++
				oa.Vel--
			} else if a.Pos > oa.Pos {
				a.Vel--
				oa.Vel++
			}
		}
	}
}

func (m *Moon) TickPos() {
	for _, a := range m {
		a.Pos += a.Vel
	}
}

func (m Moon) Kinetic() float64 {
	var ret float64
	for _, a := range m {
		ret += math.Abs(float64(a.Vel))
	}
	return ret
}

func (m Moon) Potential() float64 {
	var ret float64
	for _, a := range m {
		ret += math.Abs(float64(a.Pos))
	}
	return ret
}

type Moons []Moon

func (ms Moons) TickOne() {
	for i, m := range ms {
		m.TickVel(i, ms)
	}
	for _, m := range ms {
		m.TickPos()
	}
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
		if s == "" {
			continue
		}
		// <x=-4, y=3, z=15>
		coords := strings.Split(s[1:len(s)-1], ", ")
		var obj Moon
		for i, v := range coords {
			if pos, err := strconv.Atoi(v[2:]); err != nil {
				fmt.Printf("Failed to parse line: %s\n", s)
			} else {
				axis := Axis{int32(pos), 0}
				obj[i] = &axis
			}
		}
		moons = append(moons, obj)
	}

	moonCopy := make(Moons, len(moons))
	for i, m := range moons {
		var cpy Moon
		for j, a := range m {
			axisCpy := *a
			cpy[j] = &axisCpy
		}
		moonCopy[i] = cpy
	}

	for t := 0; t < *steps; t++ {
		moonCopy.TickOne()
	}
	var energy float64
	for _, m := range moonCopy {
		energy += m.Kinetic() * m.Potential()
	}
	fmt.Println(int(energy))

	// Part B
	origMoons := make(Moons, len(moons))
	for i, m := range moons {
		var cpy Moon
		for j, a := range m {
			axisCpy := *a
			cpy[j] = &axisCpy
		}
		origMoons[i] = cpy
	}

	var repeatCycle [3]int
	for t := 1; ; t++ {
		moons.TickOne()
		var notRepeat [3]bool
		for i, m := range moons {
			for j, a := range m {
				if *a != *origMoons[i][j] {
					notRepeat[j] = true
				}
			}
		}
		for j, nRep := range notRepeat {
			if !nRep && repeatCycle[j] == 0 {
				repeatCycle[j] = t
			}
		}
		product := 1
		for _, v := range repeatCycle {
			product *= v
		}
		if product != 0 {
			// We need to compute the LCM of the three values.
			r := repeatCycle
			result := r[0] * r[1] / gcd(r[0], r[1])
			result = result * r[2] / gcd(result, r[2])
			fmt.Println(result)
			return
		}
	}
}

func gcd(r, c int) int {
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
	return gcd
}
