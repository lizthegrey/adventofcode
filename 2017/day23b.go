package main

import (
	"fmt"
)

func main() {
	composites := 0

	for b := 106700; b <= 123700; b += 17 {
		for d := 2; d*d < b; d++ {
			if b%d == 0 {
				composites++
				break
			}
		}
	}

	fmt.Println(composites)
}
