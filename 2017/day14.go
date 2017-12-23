package main

import (
	"flag"
	"fmt"
	"knot"
)

var input = flag.String("input", "flqrgnkx", "The input to use.")

func main() {
	flag.Parse()

	set := 0
	var bitField [128][128]bool
	for i := 0; i < 128; i++ {
		key := fmt.Sprintf("%s-%d", *input, i)
		r := knot.Densify(knot.Hash(256, 64, knot.Key(key)))
		bitField[i] = toBits(r)
		set += countBits(bitField[i])
	}

	fmt.Printf("Found %d set bits.\n", set)

	var addedToRegions [128][128]bool

	regions := 0
	for i := 0; i < 128; i++ {
		for j := 0; j < 128; j++ {
			if addedToRegions[i][j] || !bitField[i][j] {
				// Only consider as seeds set tiles that are not part of regions.
				continue
			}
			regions++
			expandRegion(i, j, &bitField, &addedToRegions)
		}
	}
	fmt.Printf("Found %d non-overlapping regions.\n", regions)
}

func expandRegion(x, y int, bits, used *[128][128]bool) {
	used[x][y] = true
	// Try to expand up, down, left, and right.
	if x > 0 && bits[x-1][y] && !used[x-1][y] {
		expandRegion(x-1, y, bits, used)
	}
	if x < 127 && bits[x+1][y] && !used[x+1][y] {
		expandRegion(x+1, y, bits, used)
	}
	if y > 0 && bits[x][y-1] && !used[x][y-1] {
		expandRegion(x, y-1, bits, used)
	}
	if y < 127 && bits[x][y+1] && !used[x][y+1] {
		expandRegion(x, y+1, bits, used)
	}
}

func toBits(d []int) [128]bool {
	var ret [128]bool
	for j, v := range d {
		lower := v & 15
		upper := v >> 4

		ret[8*j+0] = upper&8 > 0
		ret[8*j+1] = upper&4 > 0
		ret[8*j+2] = upper&2 > 0
		ret[8*j+3] = upper&1 > 0

		ret[8*j+4] = lower&8 > 0
		ret[8*j+5] = lower&4 > 0
		ret[8*j+6] = lower&2 > 0
		ret[8*j+7] = lower&1 > 0
	}
	return ret
}

func display(array [128]bool) string {
	ret := ""
	for _, v := range array {
		if v {
			ret += "#"
		} else {
			ret += "."
		}
	}
	return ret
}

func countBits(array [128]bool) int {
	ret := 0
	for _, v := range array {
		if v {
			ret++
		}
	}
	return ret
}
