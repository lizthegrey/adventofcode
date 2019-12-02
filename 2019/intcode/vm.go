package intcode

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

var inputFile = flag.String("inputFile", "inputs/day02.input", "Relative file path to use as input.")

type Tape []int

func ReadInput() Tape {
	bytes, err := ioutil.ReadFile(*inputFile)
	if err != nil {
		return nil
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
	return tape
}

func (t Tape) Process() int {
	offset := 0
	for {
		if offset >= len(t) {
			fmt.Println("Ran off end of tape.")
			return -1
		}
		if t[offset] == 99 {
			return t[0]
		}
		a := t[offset+1]
		b := t[offset+2]
		dstOffset := t[offset+3]
		opLength := 0
		switch t[offset] {
		case 1:
			t[dstOffset] = t[a] + t[b]
			opLength = 4
		case 2:
			t[dstOffset] = t[a] * t[b]
			opLength = 4
		default:
			fmt.Printf("Failed to match opcode %d.\n", t[offset])
			return -1
		}
		offset += opLength
	}
}
