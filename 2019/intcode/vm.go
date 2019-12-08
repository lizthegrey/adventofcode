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

func (t Tape) Process(inputs chan int) (chan int, chan bool) {
	offset := 0
	output := make(chan int)
	done := make(chan bool, 1)
	go func() {
		for {
			if offset >= len(t) {
				fmt.Println("Ran off end of tape.")
				done <- true
				close(output)
				return
			}
			instr := t[offset] % 100
			if instr == 99 {
				done <- true
				close(output)
				return
			}
			pModes := t[offset] / 100
			instLen := map[int]int{
				1: 4,
				2: 4,
				3: 2,
				4: 2,
				5: 3,
				6: 3,
				7: 4,
				8: 4,
			}

			operands := make([]int, instLen[instr])
			for i := 1; i < len(operands); i++ {
				value := t[offset+i]
				if pModes%10 == 0 {
					// position mode
					operands[i] = t[value]
				} else {
					// literal mode
					operands[i] = value
				}
				pModes /= 10
			}
			dstOffset := t[offset+instLen[instr]-1]
			jumped := false
			switch instr {
			case 1:
				// ADD
				t[dstOffset] = operands[1] + operands[2]
			case 2:
				// MUL
				t[dstOffset] = operands[1] * operands[2]
			case 3:
				// INPUT
				t[dstOffset] = <-inputs
			case 4:
				// OUTPUT
				output <- operands[1]
			case 5:
				// JNZ
				if operands[1] != 0 {
					offset = operands[2]
					jumped = true
				}
			case 6:
				// JZ
				if operands[1] == 0 {
					offset = operands[2]
					jumped = true
				}
			case 7:
				// LT
				r := operands[1] < operands[2]
				if r {
					t[dstOffset] = 1
				} else {
					t[dstOffset] = 0
				}
			case 8:
				// EQ
				r := operands[1] == operands[2]
				if r {
					t[dstOffset] = 1
				} else {
					t[dstOffset] = 0
				}
			default:
				fmt.Printf("Failed to match opcode %d.\n", t[offset])
				done <- true
				close(output)
				return
			}
			if !jumped {
				offset += instLen[instr]
			}
		}
	}()
	return output, done
}
