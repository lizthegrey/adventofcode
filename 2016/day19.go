package main

import (
	"container/ring"
	"fmt"
)

func regular(n int) int {
	active := ring.New(n)
	for i := 0; i < n; i++ {
		active.Value = i + 1
		active = active.Next()
	}
	for ; active.Next() != active; active = active.Next() {
		active.Unlink(1)
	}
	return active.Value.(int)
}

func revised(n int) int {
	active := ring.New(n)
	var opposite *ring.Ring
	for i := 0; i < n; i++ {
		active.Value = i + 1
		if i == (n-1)/2 {
			opposite = active
		}

		active = active.Next()
	}

	rCount := 0
	for ; rCount != n-1; active = active.Next() {
		opposite = opposite.Prev()
		opposite.Unlink(1)
		if rCount%2 == 0 {
			opposite = opposite.Next().Next()
		} else {
			opposite = opposite.Next()
		}
		rCount++
	}
	return active.Value.(int)
}

func main() {
	fmt.Println(regular(5))
	fmt.Println(regular(3014603))
	fmt.Println(revised(5))
	fmt.Println(revised(3014603))
}
