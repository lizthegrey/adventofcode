package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

var inputFile = flag.String("inputFile", "inputs/day15.input", "Relative file path to use as input.")

type coord struct {
	x, y int
}

func (c coord) radius(o coord) int {
	var ret int
	if c.x < o.x {
		ret += o.x - c.x
	} else {
		ret += c.x - o.x
	}
	if c.y < o.y {
		ret += o.y - c.y
	} else {
		ret += c.y - o.y
	}
	return ret
}

type readings map[coord]coord

type interval struct {
	lo, hi int
}

func (i interval) width() int {
	return i.hi - i.lo + 1
}

func (i interval) contains(o interval) bool {
	return i.lo <= o.lo && o.hi <= i.hi
}

// There's an overlap if other's bound falls within my bounds.
// Doesn't cover full contains case but that's what contains is for.
func (i interval) overlaps(o interval) bool {
	return (i.lo <= o.lo && o.lo <= i.hi) || (i.lo <= o.hi && o.hi <= i.hi)
}

func (i interval) union(o interval) *interval {
	if i.contains(o) {
		return &i
	}
	if o.contains(i) {
		return &o
	}
	if !i.overlaps(o) {
		return nil
	}
	// We have a partial overlap.
	var combined interval
	if i.lo < o.lo {
		combined.lo = i.lo
	} else {
		combined.lo = o.lo
	}
	if i.hi > o.hi {
		combined.hi = i.hi
	} else {
		combined.hi = o.hi
	}
	return &combined
}

// intervals _must_ be already non-overlapping and sorted by
// lo value, ascending. Use union to insert into a blank intervals
// in order to preserve this invariant.
type intervals []interval

func (is intervals) intersect(bounds interval) intervals {
	ret := make(intervals, 0, len(is))
	for _, v := range is {
		if v.hi < bounds.lo {
			// Out of range, skip
			continue
		}
		if v.lo > bounds.hi {
			// Out of range, we're done since list is sorted.
			break
		}
		// Definitely starts or ends in range. Clamp any weird bounds.
		lo := v.lo
		if lo < bounds.lo {
			lo = bounds.lo
		}
		hi := v.hi
		if hi > bounds.hi {
			hi = bounds.hi
		}
		ret = append(ret, interval{lo, hi})
	}
	return ret
}

func (is intervals) union(os ...interval) intervals {
	// special case: adding to an empty intervals.
	if is == nil || len(is) == 0 {
		return os
	}

	// First, sort the intervals into one list, then step through and cull overlapping ones.
	tmp := make(intervals, 0, len(is)+len(os))
	for len(is)+len(os) > 0 {
		if len(is) == 0 || (len(os) > 0 && os[0].lo < is[0].lo) {
			tmp = append(tmp, os[0])
			os = os[1:]
		} else {
			tmp = append(tmp, is[0])
			is = is[1:]
		}
	}

	ret := make(intervals, 0, len(tmp))
	cur := tmp[0]
	for i := 1; i < len(tmp); i++ {
		if overlap := cur.union(tmp[i]); overlap != nil {
			cur = *overlap
		} else {
			ret = append(ret, cur)
			cur = tmp[i]
		}
	}
	ret = append(ret, cur)
	return ret
}

func (is intervals) totalWidth() int {
	var sum int
	for _, i := range is {
		sum += i.width()
	}
	return sum
}

func main() {
	flag.Parse()
	bytes, err := ioutil.ReadFile(*inputFile)
	if err != nil {
		return
	}
	contents := string(bytes)
	split := strings.Split(contents, "\n")

	nearest := make(readings)
	for _, s := range split[:len(split)-1] {
		// Sensor at x=391282, y=2038170: closest beacon is at x=-532461, y=2166525
		parts := strings.Split(s, " ")
		locX, _ := strconv.Atoi(parts[2][2 : len(parts[2])-1])
		locY, _ := strconv.Atoi(parts[3][2 : len(parts[3])-1])

		becX, _ := strconv.Atoi(parts[8][2 : len(parts[8])-1])
		becY, _ := strconv.Atoi(parts[9][2:])

		loc := coord{locX, locY}
		bec := coord{becX, becY}
		nearest[loc] = bec
	}

	// part A
	rowY := 2000000
	rowExclusions := nearest.computeExclusions(rowY)

	// If the sensor is in that position, it should get excluded.
	// Subtract back out the places where there actually is a beacon.
	seen := make(map[coord]bool)
	for _, v := range nearest {
		if v.y == rowY {
			seen[v] = true
		}
	}
	fmt.Println(rowExclusions.totalWidth() - len(seen))

	// part B
	maxCoord := 4000000
	bounds := interval{0, maxCoord}
	for testY := 0; testY < maxCoord; testY++ {
		intersection := nearest.computeExclusions(testY).intersect(bounds)
		if intersection.totalWidth() == maxCoord-1+1 {
			// If we've found exactly 4000000 excluded positions, then we want to know what
			// the non-excluded position is.
			x := intersection[0].hi + 1
			if x != intersection[1].lo-1 {
				fmt.Printf("Something is wrong: %v\n", intersection)
			}
			tuning := x*maxCoord + testY
			fmt.Println(tuning)
			break
		}
	}
}

// Basic algorithm: the radius of an exclusion zone is equal to
// (radius - abs(y_test - y_sensor)) centered around x_sensor
func (r readings) computeExclusions(rowY int) intervals {
	var exclusions intervals
	for k, v := range r {
		midpoint := coord{k.x, rowY}
		deltaX := k.radius(v) - k.radius(midpoint)
		exclusions = exclusions.union(interval{k.x - deltaX, k.x + deltaX})
	}
	return exclusions
}
