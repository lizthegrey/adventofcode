package main

import (
	"flag"
	"fmt"
)

func main() {
	flag.Parse()
	if !BoolCheck() {
		fmt.Println("Failed boolean check.")
	}
	results := Mystery()
	fmt.Printf("Solutions: %d, %d\n", results[0], results[len(results)-1])
}

func BoolCheck() bool {
	// 0  seti 123 2
	r2 := 123

	// 1  bani 2 456 2
	r2 &= 456

	// 2  eqri 2 72 2
	// 3  addr 2 4 4
	// 4  seti 0 0 4
	if r2 != 72 {
		// Loop forever, but for realsies return false.
		return false
	}
	return true
}

func Mystery() []int {
	seen := make(map[int]bool)
	ordered := make([]int, 0)

	r5 := 0

	// Main program starts here.
	// 5  seti 0 2
	v := 0

outer:
	for {
		// 6  bori 2 65536 5
		// Note: this constant is 0b10000000000000000
		r5 = v | 65536

		// 7  seti 16123384 2
		// Note: this constant is 0b111101100000010111111000
		v = 16123384

		for {
			// Mask off everything but the lower 8 bits of r5.
			// 8  bani 5 255 3
			// Note: this constant is 0b11111111
			masked := r5 & 255

			// Add at most 255 to v
			// 9  addr 2 3 2
			v += masked

			// Mask v to the lowest 24 bits.
			// 10 bani 2 16777215 2
			// Note: this constant is 0b111111111111111111111111
			v &= 16777215

			// 11 muli 2 65899 2
			v *= 65899

			// Mask v to the lowest 24 bits.
			// 12 bani 2 16777215 2
			// Note: this constant is 0b111111111111111111111111
			v &= 16777215

			// 13 gtir 256 5 3
			// 14 addr 3 4 4
			// 15 addi 4 1 4
			// NOTE: found a bug here in my solution from last night. >=, not >.
			if r5 >= 256 {
				// Replaced: r5 = SlowDivide256(r5) -- starting at instr 17
				// r5 = SlowDivide256(r5)

				r5 /= 256
				// Jump to line 8.
				// 27 seti 7 4
				continue
			}

			// Jump to line 28.
			// 16 seti 27 4

			// We may jump into this code.
			// 28 eqrr 2 0 3

			// if input == v {
			// We would terminate.
			// 29 addr 3 4 4
			// return
			// } else {

			// Record the value we would check against instead of terminating.
			if seen[v] {
				return ordered
			}
			seen[v] = true
			ordered = append(ordered, v)

			// Loop back to line 6.
			// 30 seti 5 4
			// We adjusted r5 and v.
			continue outer
		}
	}
}

func SlowDivide256(r5 int) int {
	// This loop finds the next multiple of 256 that's larger than r5.
	// e.g. we can replace it with:
	// 0-255: 0 returned
	// 256-511: 1 returned
	// ...

	// 17 seti 0 3
	for i := 0; ; i++ {
		// 18 addi 3 1 1
		// 19 muli 1 256 1
		r1 := (i + 1) * 256

		// 20 gtrr 1 5 1
		// 21 addr 1 4 4
		// 22 addi 4 1 4
		if r1 > r5 {
			// 23 seti 25 4
			// jump to line 26

			// 26 setr 3 5
			// We may jump into this code.
			r5 = i

			return i
		}
		// 24 addi 3 1 3
		// Loop back to line 18.
		// 25 seti 17 4
	}
}
