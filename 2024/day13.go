package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"math"
	"strconv"
	"strings"
)

var inputFile = flag.String("inputFile", "inputs/day13.input", "Relative file path to use as input.")

type scenario struct {
	ax, ay, bx, by, tx, ty int64
}

const offset int64 = 10000000000000

func main() {
	flag.Parse()
	bytes, err := ioutil.ReadFile(*inputFile)
	if err != nil {
		return
	}
	contents := string(bytes)
	split := strings.Split(contents, "\n")

	var scenarios []scenario
	for i := 0; i+4 <= len(split); i += 4 {
		var s scenario
		partsA := strings.Split(split[i][12:], ", Y+")
		s.ax, _ = strconv.ParseInt(partsA[0], 10, 64)
		s.ay, _ = strconv.ParseInt(partsA[1], 10, 64)
		partsB := strings.Split(split[i+1][12:], ", Y+")
		s.bx, _ = strconv.ParseInt(partsB[0], 10, 64)
		s.by, _ = strconv.ParseInt(partsB[1], 10, 64)
		target := strings.Split(split[i+2][9:], ", Y=")
		s.tx, _ = strconv.ParseInt(target[0], 10, 64)
		s.ty, _ = strconv.ParseInt(target[1], 10, 64)
		scenarios = append(scenarios, s)
	}

	// part A
	var sumA int64
	for _, s := range scenarios {
		var minCost int64 = math.MaxInt64
		for i := int64(1); i <= 100; i++ {
			for j := int64(1); j <= 100; j++ {
				if s.ax*i+s.bx*j == s.tx && s.ay*i+s.by*j == s.ty {
					cost := 3*i + j
					if cost < minCost {
						minCost = cost
					}
					// No point in going past winning point.
					break
				}
			}
		}
		if minCost != math.MaxInt64 {
			sumA += minCost
		}
	}
	fmt.Println(sumA)

	// part B
	var sumB int64
	for _, s := range scenarios {
		s.tx += offset
		s.ty += offset
		if s.tx%gcd(s.ax, s.bx) != 0 || s.ty%gcd(s.ay, s.by) != 0 {
			// There is no solution, because a divisibility factor is missing.
			continue
		}

		// two remaining cases: either the lines intersect at exactly one point
		// (which needs to be checked for being an integer), or
		// the lines exactly overlap, in which case we should choose to pick the smaller of 3*a or b
		// although this didn't come up in my input

		// a*ax + b*bx == tx
		// a*ay + b*by == ty
		// grab one equation, solve for a
		// a*ax = tx - b*bx
		// a = (tx - b*bx) / ax
		// substitute into the other equation etc
		// ay * (tx - b*bx) / ax + b * by = ty
		// ay * tx / ax - b * bx * ay / ax + b * by = ty
		// b * (by - ay * bx / ax) = ty - ay * tx / ax
		// b = (ty - ay * tx / ax) / (by - ay * bx / ax)
		// b = (ty * ax - ay * tx) / (by * ax - ay * bx)
		// then substitute back through again, or just know it's symmetrical
		// a = -(ty * bx - by * tx) / (by * ax - ay * bx)
		denom := s.by*s.ax - s.ay*s.bx
		if denom == 0 {
			// this is the case where the lines exactly overlap. skip for now.
			continue
		}
		a := (s.tx*s.by - s.ty*s.bx)
		if a%denom != 0 {
			// solution is not integral
			continue
		}
		a /= denom
		b := (s.ty*s.ax - s.tx*s.ay)
		if b%denom != 0 {
			// solution is not integral
			continue
		}
		b /= denom
		if a < 0 || b < 0 {
			continue
		}
		sumB += 3*a + b
	}
	fmt.Println(sumB)
}

func gcd(a, b int64) int64 {
	for a != b {
		if a > b {
			a -= b
		} else {
			b -= a
		}
	}
	return a
}
