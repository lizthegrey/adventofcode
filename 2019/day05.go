package main

import (
	"flag"
	"fmt"
	"github.com/lizthegrey/adventofcode/2019/intcode"
)

var debug = flag.Bool("debug", false, "Print debug info as we go along.")
var inputValue = flag.Int("inputValue", 0, "The input to the input instruction.")

func main() {
	flag.Parse()
	tape := intcode.ReadInput()
	if tape == nil {
		fmt.Println("Failed to parse input.")
		return
	}
	if *debug {
		for _, n := range tape {
			fmt.Printf("%d,", n)
		}
		fmt.Println()
	}

	input := make(chan int)
	result, _ := tape.Process(input)
	input <- *inputValue
	fmt.Printf("Result: %d\n", <-result)
}
