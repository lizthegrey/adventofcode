package main

import (
	"flag"
	"fmt"
	"github.com/lizthegrey/adventofcode/2019/intcode"
)

var inputFile = flag.String("inputFile", "inputs/day05.input", "Relative file path to use as input.")
var inputValue = flag.Int("inputValue", 0, "The input to the input instruction.")

func main() {
	flag.Parse()
	tape := intcode.ReadInput(*inputFile)
	if tape == nil {
		fmt.Println("Failed to parse input.")
		return
	}

	input := make(chan int)
	result, _ := tape.Process(input)
	input <- *inputValue
	fmt.Printf("Result: %d\n", <-result)
}
