package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"math"
	"slices"
	"strings"
)

var inputFile = flag.String("inputFile", "inputs/day09.input", "Relative file path to use as input.")

const empty uint16 = math.MaxUint16

func main() {
	flag.Parse()
	bytes, err := ioutil.ReadFile(*inputFile)
	if err != nil {
		return
	}
	contents := string(bytes)
	split := strings.Split(contents, "\n")

	// File numbers range from 0 to 19,999. uint16 holds up to 65535
	var mem []uint16
	var file bool
	var maxFile uint16
	positions := make(map[uint16]int)
	lengths := make(map[uint16]int)
	for n, v := range split[0] {
		file = !file
		digit := int(v - '0')
		val := uint16(n / 2)
		if !file {
			val = empty
		} else {
			maxFile = val
			positions[val] = len(mem)
			lengths[val] = digit
		}
		mem = append(mem, slices.Repeat([]uint16{val}, digit)...)
	}

	a := make([]uint16, len(mem))
	copy(a, mem)
	right := len(a)
	for hole := 0; hole < len(a); hole++ {
		fill := a[hole]
		if fill != empty {
			continue
		}
		// Stop when the two pointers collide.
		if right <= hole {
			break
		}
		for right--; a[right] == empty; right-- {
		}
		a[hole] = a[right]
		a[right] = empty
	}
	fmt.Println(checksum(a))

	b := make([]uint16, len(mem))
	copy(b, mem)
	for file := maxFile; file != empty; file-- {
		length := lengths[file]
		pos := positions[file]
	search:
		for i := 0; i < pos; i++ {
			for j := 0; j < length; j++ {
				if b[i+j] != empty {
					continue search
				}
			}
			copy(b[i:i+length], slices.Repeat([]uint16{file}, length))
			copy(b[pos:pos+length], slices.Repeat([]uint16{empty}, length))
			break
		}
	}
	fmt.Println(checksum(b))
}

func checksum(mem []uint16) uint64 {
	var ret uint64
	for i, v := range mem {
		if v == empty {
			continue
		}
		ret += uint64(i) * uint64(v)
	}
	return ret
}
