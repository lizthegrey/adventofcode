package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"math/bits"
	"strings"
)

var inputFile = flag.String("inputFile", "inputs/day13.input", "Relative file path to use as input.")

func main() {
	flag.Parse()
	bytes, err := ioutil.ReadFile(*inputFile)
	if err != nil {
		return
	}
	contents := string(bytes)
	split := strings.Split(contents, "\n")

	var verticalSum, horizontalSum int
	var verticalNear, horizontalNear int

	var r int
	var rChecksums, cChecksums []uint64

	for _, s := range split {
		if s == "" {
			mirrorRow, nearRow := process(rChecksums)
			// 1-indexed values.
			verticalSum += mirrorRow + 1
			verticalNear += nearRow + 1

			mirrorCol, nearCol := process(cChecksums)
			// 1-indexed values.
			horizontalSum += mirrorCol + 1
			horizontalNear += nearCol + 1

			r = 0
			rChecksums = nil
			cChecksums = nil
			continue
		}
		if r == 0 {
			cChecksums = make([]uint64, len(s))
		}
		rChecksums = append(rChecksums, 0)
		for c, v := range s {
			if v == '#' {
				// Increment the horizontal and vertical checksums
				// for the row/column that we are in.
				rChecksums[r] += 1 << c
				cChecksums[c] += 1 << r
			}
		}
		r++
	}
	fmt.Println(100*verticalSum + horizontalSum)
	fmt.Println(100*verticalNear + horizontalNear)
}

func process(checksums []uint64) (int, int) {
	m := -1
	v := -1
	for mirror := 0; mirror < len(checksums); mirror++ {
		valid := true
		var nearMiss bool
		var offset int
		for ; mirror-offset >= 0 && mirror+offset+1 < len(checksums); offset++ {
			left := checksums[mirror-offset]
			right := checksums[mirror+offset+1]
			if checksums[mirror-offset] != checksums[mirror+offset+1] {
				valid = false
				if bits.OnesCount64(left^right) == 1 && !nearMiss {
					nearMiss = true
				} else {
					// Too many failures, bail.
					nearMiss = false
					break
				}
			}
		}
		if nearMiss {
			m = mirror
		}
		if valid && offset > 0 {
			v = mirror
		}
	}
	return v, m
}
