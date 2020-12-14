package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

var inputFile = flag.String("inputFile", "inputs/day14.input", "Relative file path to use as input.")
var debug = flag.Bool("debug", false, "Whether to print debug output along the way.")
var partB = flag.Bool("partB", false, "Whether to use part B logic.")

func main() {
	flag.Parse()
	bytes, err := ioutil.ReadFile(*inputFile)
	if err != nil {
		return
	}
	contents := string(bytes)
	split := strings.Split(contents, "\n")
	split = split[:len(split)-1]

	memory := make(map[uint64]uint64)
	// Mask XXXXXXXXXXXXXXX01 ->
	// Off: 00000000000000010
	// On:  00000000000000001
	// In order to apply the mask, we will bitwise
	// AND NOT (&^) with off and OR (|) with `on`.
	var maskOff, maskOn uint64
	for _, s := range split {
		parts := strings.Split(s, " = ")
		switch parts[0][0:3] {
		case "mem":
			addr, err := strconv.Atoi(parts[0][4 : len(parts[0])-1])
			operand, err := strconv.Atoi(parts[1])
			if err != nil {
				fmt.Printf("Failed to parse %s\n", s)
				break
			}
			if !*partB {
				bits := (uint64(operand) | maskOn) &^ maskOff
				memory[uint64(addr)] = bits
				if *debug {
					fmt.Printf("Wrote: %d <- %36b\n", bits, addr)
				}
			} else {
				value := uint64(operand)
				dstOverwritten := uint64(addr) | maskOn
				// Floating bits: maskOff == 0 && maskOn == 0
				valueSet := maskOff | maskOn
				var dsts []uint64
				for bit := 0; bit < 36; bit++ {
					bitAtPos := (dstOverwritten >> bit) & 1
					floating := (0 == (valueSet>>bit)&1)
					// if it is floating, take every value out of previous dsts
					// and prefix 1,0 to each of them
					// otherwise, prefix bitAtPos to each value from previous dsts
					var newDsts []uint64
					if bit != 0 {
						for _, v := range dsts {
							if floating {
								newDsts = append(newDsts, v|uint64(1<<bit), v)
							} else {
								newDsts = append(newDsts, v|(uint64(bitAtPos)<<bit))
							}
						}
					} else {
						if floating {
							newDsts = append(newDsts, uint64(1), uint64(0))
						} else {
							newDsts = append(newDsts, uint64(bitAtPos))
						}
					}
					dsts = newDsts
				}
				for _, d := range dsts {
					memory[d] = value
					if *debug {
						fmt.Printf("Wrote: %d <- %36b\n", d, addr)
					}
				}
			}
		case "mas":
			for i, bit := range parts[1] {
				switch bit {
				case 'X':
					maskOn &^= 1 << (35 - i)
					maskOff &^= 1 << (35 - i)
				case '0':
					maskOn &^= 1 << (35 - i)
					maskOff |= 1 << (35 - i)
				case '1':
					maskOff &^= 1 << (35 - i)
					maskOn |= 1 << (35 - i)
				}
				if *debug {
					fmt.Printf("Off: %36b\nOn:  %36b\n", maskOff, maskOn)
				}
			}
		default:
			fmt.Printf("Saw unexpected operand %s\n", parts[0])
		}
	}
	var sum uint64
	for _, v := range memory {
		sum += v
	}
	fmt.Println(sum)
}
