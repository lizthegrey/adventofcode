package main

import (
	"fmt"
)

func main() {
test:
	for i := 0; i < 10; i++ {
		if i%1000 == 0 {
			fmt.Printf("Testing %d\n", i)
		}
		fmt.Println()
		last := 1

		b := 0
		c := 0
		d := i + 365*7
		a := d

		for {
			b = a
			a = 0

			c = b % 2
			a = b / 2

			var out int
			if c == 1 {
				out = 1
			} else if c == 0 {
				out = 0
			} else {
				fmt.Println("impossible")
			}

			if !emit(out, last) {
				continue test
			}
			last = out
			if a == 0 {
				continue
			}
			a = d
		}
	}
}

func emit(b int, last int) bool {
	fmt.Printf("%d", b)
	if b == last {
		return false
	}
	return true
}
