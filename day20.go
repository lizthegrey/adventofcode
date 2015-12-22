package main

import (
	"fmt"
)

func main() {
	i := 1
	for divisorSum(i) < 36000000 {
		i++
	}
	fmt.Println(i)
}

func divisorSum(k int) int {
	sum := 0
	for i := 1; i * i < k; i++ {
		if k % i == 0 {
			if i <= 50 {
				sum += k / i
			}
			if k / i <= 50 {
				sum += i
			}
		}
	}
	return sum * 11
}
