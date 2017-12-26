package main

import (
	"fmt"
)

func main() {
test:
	for i := 0; ; i++ {
		successes := 0
		last := 1

		d := i + 365*7
		b := d

		for {
			c := b % 2
			b = b >> 1

			if !emit(c, last, &successes) {
				continue test
			}
			if successes == 100 {
				fmt.Println(i)
				return
			}
			last = c
			if b != 0 {
				continue
			} else {
				b = d
			}
		}
	}
}

func emit(b int, last int, successes *int) bool {
	if b == last {
		return false
	}
	*successes++
	return true
}
