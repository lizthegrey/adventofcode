package main

import (
	"fmt"
)

func main() {
	composites := 0

	for b := 106700; b != 123700; b += 17 {
		composite := false
		for d := 2; d < b; d++ {
			for e := 2; e < b; e++ {
				if d*e == b {
					composite = true
				}
			}
		}
		if composite {
			composites++
		}
	}

	fmt.Println(composites)
}
