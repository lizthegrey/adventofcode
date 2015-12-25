package main

import (
	"fmt"
)

func main() {
	n := getItem(3010, 3019)
	fmt.Println(n)
	v := 20151125
	for i := 0; i < n; i++ {
		v *= 252533
		v %= 33554393
	}
	fmt.Println(v)
}

func getItem(row, col int) int {
	result := 0
	for {
		if row == 1 && col == 1 {
			return result
		}
		result++
		if col > 1 {
			col--
			row++
		} else {
			col = row - 1
			row = 1
		}
	}
}
