package intcode

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"time"
)

type Tape map[int]int

func ReadInput(file string) Tape {
	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		return nil
	}
	contents := string(bytes[:len(bytes)-1])
	split := strings.Split(contents, ",")
	tape := make(Tape)
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

func (t Tape) Copy() Tape {
	c := make(Tape)
	for k, v := range t {
		c[k] = v
	}
	return c
}

func (t Tape) Process(inputs chan int) (chan int, chan bool) {
	output, done, _ := t.ProcessInternal(inputs, true, -1)
	return output, done
}

func (t Tape) ProcessNonBlocking(inputs chan int, tid int) (chan int, chan bool, *bool) {
	return t.ProcessInternal(inputs, false, tid)
}

func (t Tape) ProcessInternal(inputs chan int, blocking bool, tid int) (chan int, chan bool, *bool) {
	offset := 0
	var output chan int
	if blocking {
		output = make(chan int)
	} else {
		output = make(chan int, 500)
	}
	done := make(chan bool, 1)
	base := 0
	awaitingInput := false
	go func() {
		for {
			instr := t[offset] % 100
			if instr == 99 {
				done <- true
				close(output)
				return
			}
			if tid >= 0 {
				fmt.Printf("Trace %d: executing %d at ip %d\n", tid, instr, offset)
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
				9: 2,
			}

			operands := make([]int, instLen[instr])
			dstOffset := t[offset+instLen[instr]-1]
			for i := 1; i < len(operands); i++ {
				value := t[offset+i]
				switch pModes % 10 {
				case 0:
					// position mode
					operands[i] = t[value]
				case 1:
					// literal mode
					operands[i] = value
				case 2:
					operands[i] = t[base+value]
					if i == len(operands)-1 {
						dstOffset += base
					}
				}
				pModes /= 10
			}
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
				ret := -1
				if blocking {
					ret = <-inputs
				} else {
					select {
					case ret = <-inputs:
						// We've successfully read our input channel.
						awaitingInput = false
						if tid >= 0 {
							fmt.Printf("[%d] Received value %d\n", tid, ret)
						}
					default:
						// Use default value set above.
						time.Sleep(time.Millisecond)
						awaitingInput = true
					}
				}
				t[dstOffset] = ret
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
			case 9:
				base += operands[1]
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
	return output, done, &awaitingInput
}
