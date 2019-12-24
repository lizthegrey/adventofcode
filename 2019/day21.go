package main

import (
	"flag"
	"fmt"
	"github.com/lizthegrey/adventofcode/2019/intcode"
	"time"
)

var inputFile = flag.String("inputFile", "inputs/day21.input", "Relative file path to use as input.")

func main() {
	flag.Parse()
	tape := intcode.ReadInput(*inputFile)
	if tape == nil {
		fmt.Println("Failed to parse input.")
		return
	}

	// Part A
	workingTape := tape.Copy()
	input := make(chan int, 1)
	output, done := workingTape.Process(input)

	go func() {
		for l := range output {
			if l > 255 {
				fmt.Println(l)
				return
			}
			fmt.Printf("%c", l)
		}
		fmt.Println()
	}()

	inputString := []string{
		// J = AND(OR(C, B, A), D)
		"NOT C J",
		"NOT B T",
		"OR T J",
		"NOT A T",
		"OR T J",
		"AND D J",
		"WALK",
	}
	for _, l := range inputString {
		for _, r := range l {
			input <- int(r)
		}
		input <- int('\n')
	}
	<-done
	time.Sleep(100 * time.Millisecond)

	// Part B
	input = make(chan int, 1)
	output, done = tape.Process(input)

	go func() {
		for l := range output {
			if l > 255 {
				fmt.Println(l)
				return
			}
			fmt.Printf("%c", l)
		}
		fmt.Println()
	}()

	inputString = []string{
		// J = AND(OR(C, B, A), D)
		// except we also need to make sure E or H are clear before jumping.
		// J = AND(J, OR(E,H))
		"NOT C J",
		"NOT B T",
		"OR T J",
		"NOT A T",
		"OR T J",
		"AND D J",
		"OR E T",
		"OR H T",
		"AND T J",
		"RUN",
	}
	for _, l := range inputString {
		for _, r := range l {
			input <- int(r)
		}
		input <- int('\n')
	}
	<-done
	time.Sleep(100 * time.Millisecond)
}
