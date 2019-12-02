package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

var inputFile = flag.String("inputFile", "inputs/day02.input", "Relative file path to use as input.")
var partB = flag.Bool("partB", false, "Whether to use the Part B logic.")
var debug = flag.Bool("debug", false, "Print debug info as we go along.")

const expected = 19690720

func main() {
	flag.Parse()
	bytes, err := ioutil.ReadFile(*inputFile)
	if err != nil {
		return
	}
	contents := string(bytes[:len(bytes)-1])
	split := strings.Split(contents, ",")
	tape := make([]int, len(split))
	for i, s := range split {
		n, err := strconv.Atoi(s)
		if err != nil {
			fmt.Printf("Failed to parse %s\n", s)
			break
		}
		tape[i] = n
	}
	if *debug {
		for _, n := range tape {
			fmt.Printf("%d,", n)
		}
		fmt.Println()
	}

	if !*partB {
		workingTape := make([]int, len(tape))
		copy(workingTape, tape)
		workingTape[1] = 12
		workingTape[2] = 2
		fmt.Println(process(workingTape))
	} else {
		for noun := 0; noun < 100; noun++ {
			for verb := 0; verb < 100; verb++ {
				workingTape := make([]int, len(tape))
				copy(workingTape, tape)
				workingTape[1] = noun
				workingTape[2] = verb
				if process(workingTape) == expected {
					fmt.Printf("Noun = %d, Verb = %d, Result = %d\n", noun, verb, 100*noun+verb)
					return
				}
			}
		}
		fmt.Println("Failed to find solution.")
	}
}

func process(tape []int) int {
	offset := 0
	for {
		if offset >= len(tape) {
			fmt.Println("Ran off end of tape.")
			return -1
		}
		if tape[offset] == 99 {
			if *debug {
				fmt.Printf("Exited normally at offset %d.\n", offset)
			}
			return tape[0]
		}
		a := tape[offset+1]
		b := tape[offset+2]
		dstOffset := tape[offset+3]
		opLength := 0
		switch tape[offset] {
		case 1:
			tape[dstOffset] = tape[a] + tape[b]
			opLength = 4
		case 2:
			tape[dstOffset] = tape[a] * tape[b]
			opLength = 4
		default:
			fmt.Printf("Failed to match opcode %d.\n", tape[offset])
			return -1
		}
		offset += opLength
	}
}
