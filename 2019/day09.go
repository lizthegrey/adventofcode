package main

import (
	"flag"
	"fmt"
	"github.com/lizthegrey/adventofcode/2019/intcode"
)

var inputFile = flag.String("inputFile", "inputs/day09.input", "Relative file path to use as input.")

func main() {
	flag.Parse()
	tape := intcode.ReadInput(*inputFile)
	if tape == nil {
		fmt.Println("Failed to parse input.")
		return
	}

	for i := 1; i <= 2; i++ {
		workingTape := tape.Copy()
		input := make(chan int, 1)
		input <- i
		output, _ := workingTape.Process(input)
		for out := range output {
			fmt.Printf("%d,", out)
		}
		fmt.Println()
	}
}
