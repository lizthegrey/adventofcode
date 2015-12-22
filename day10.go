package main

import (
	"fmt"
	"strconv"
)

func main() {
	in := make(chan byte)
	out := make(chan byte)

	input := out

	for i := 0; i < 50; i++ {
		in = out
		out = make(chan byte)
		go machine(in, out)
	}

	output := out

	str := "3113322113"
	go func() {
		for c := range str {
			input <- str[c]
		}
		close(input)
	}()
	count := 0
	for _ = range output {
		count++
	}
	fmt.Println(count)
}

func machine(in, out chan byte) {
	lastDigit := byte('x')
	count := 0
	for curDigit := range in {
		if lastDigit != curDigit && count != 0 {
			countStr := strconv.Itoa(count)
			for i := range countStr {
				out <- countStr[i]
			}
			out <- lastDigit
			count = 0
		}
		count++
		lastDigit = curDigit
	}
	countStr := strconv.Itoa(count)
	for i := range countStr {
		out <- countStr[i]
	}
	out <- lastDigit
	close(out)
}
