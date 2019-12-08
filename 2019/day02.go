package main

import (
	"flag"
	"fmt"
	"github.com/lizthegrey/adventofcode/2019/intcode"
)

var partB = flag.Bool("partB", false, "Whether to use the Part B logic.")
var debug = flag.Bool("debug", false, "Print debug info as we go along.")

const expected = 19690720

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

	if !*partB {
		workingTape := make(intcode.Tape, len(tape))
		copy(workingTape, tape)
		workingTape[1] = 12
		workingTape[2] = 2
		_, done := workingTape.Process(nil)
		<-done
		fmt.Println(workingTape[0])
	} else {
		for noun := 0; noun < 100; noun++ {
			for verb := 0; verb < 100; verb++ {
				workingTape := make(intcode.Tape, len(tape))
				copy(workingTape, tape)
				workingTape[1] = noun
				workingTape[2] = verb
				_, done := workingTape.Process(nil)
				<-done
				if workingTape[0] == expected {
					fmt.Printf("Noun = %d, Verb = %d, Result = %d\n", noun, verb, 100*noun+verb)
					return
				}
			}
		}
		fmt.Println("Failed to find solution.")
	}
}
