package main

import (
	"flag"
	"fmt"
)

var partB = flag.Bool("partB", true, "Whether to use the part B input.")

func main() {
	flag.Parse()

	// Register 0 starts with value 1 if in part b; value 0 if in part a.
	r0 := 0

	if *partB {
		r0 = 1
	}
	r1 := 0
	r4 := 0

	// skip ahead by 17 instructions on first run.
	// 0  addi 5 16 5
	// 17 addi 1 2 1
	r1 += 2

	// 18 mulr 1 1 1
	r1 *= r1

	// 19 mulr 5 1 1
	r1 *= 19

	// 20 muli 1 11 1
	r1 *= 11

	// 21 addi 4 1 4
	r4++

	// 22 mulr 4 5 4
	r4 *= 22

	// 23 addi 4 9 4
	r4 += 9

	// 24 addr 1 4 1
	r1 += r4

	// Flow control: if r0 == 0 then do the next one, otherwise skip 2 ahead.
	// 25 addr 5 0 5
	// 26 seti 0 5
	// skip down to instruction 1, don't execute this next block if r0 is 0.
	if r0 != 0 {
		// 27 setr 5 4
		r4 = 27

		// 28 mulr 4 5 4
		r4 *= 28

		// 29 addr 5 4 4
		r4 += 29

		// 30 mulr 5 4 4
		r4 *= 30

		// 31 muli 4 14 4
		r4 *= 14

		// 32 mulr 4 5 4
		r4 *= 32

		// 33 addr 1 4 1
		r1 += r4

		// 34 seti 0 0
		r0 = 0

		// 35 seti 0 5
		// Continue to instruction 1.
	}

	fmt.Printf("Sending input r1 = %d\n", r1)
	// r0 = operateSlow(r1)
	r0 = sumFactors(r1)
	fmt.Println(r0)
}

func sumFactors(v int) int {
	sum := 1 + v
	if v%2 == 0 {
		sum += 2 + (v / 2)
	}
	for x := 3; x*x < v; x += 2 {
		if v%x == 0 {
			sum += x
			sum += v / x
		}
	}
	return sum
}

func operateSlow(r1 int) int {
	result := 0

	// 1  seti 1 3
	r3 := 1

	for {
		// 2  seti 1 2
		r2 := 1

		for {
			// Multiply registers 3 and 2, store in register 4
			// 3  mulr 3 2 4
			r4 := r2 * r3

			// If registers 4 and 1 are equal, write 1 to register 4, otherwise 0
			// 4  eqrr 4 1 4
			// Add the value of register 4 (which should be 0 or 1) to register 5.
			// 5  addr 4 5 5
			// If register 4 and 1 were not equal, then this line runs (PATH A).
			// It adds 1 to register 5, causing another line skip.
			// 6  addi 5 1 5
			// We skip to here if registers 4 and 1 were equal before. (PATH B)
			// Add registers 3 and 0 together, and store in register 0.
			// 7  addr 3 0 0
			if r4 == r1 {
				result = result + r3
			}

			// PATHS A AND B CONVERGE HERE.
			// Add 1 to register 2, storing in register 2.
			// 8  addi 2 1 2
			r2++

			// if r2 > r1 then write 1 to r4, otherwise write 0 to r4.
			// 9  gtrr 2 1 4
			if r2 > r1 {
				// Break the inner for.
				break
			}

			// Add register 4's value (which should be 0 or 1)  to r5.
			// 10 addr 5 4 5
			// If r2 <= r1, then this line runs.
			// Set r5 to 2 (which then gets incremented to 3).
			// 11 seti 2 5
		} // end for.

		// Increment r3 by 1.
		// 12 addi 3 1 3
		r3++

		// If r3 > r1, then write 1 to r4, otherwise 0.
		// 13 gtrr 3 1 4
		// Add register 4 (which is 0 or 1) to register 5.
		// 14 addr 4 5 5
		// This runs only if r3 <= r1.
		// Set r5 to 1 (which then gets incremented to 2.
		// 15 seti 1 5
		// Loop back to instruction 2 (setting r2 to 1)
		if r3 > r1 {
			break
		}
	} // end for

	// Otherwise: if r3 > r1:
	// 16 mulr 5 5 5
	// r5 = r5 * r5, which will be way out of bounds. 16 * 16 = 256, + 1 = 257.
	return result
}
