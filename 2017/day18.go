package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

var inputFile = flag.String("inputFile", "inputs/day18.input", "Relative file path to use as input.")

type Machine struct {
	Ip        int
	Program   []Instruction
	Registers map[byte]int
}

type Instruction func(map[byte]int, *int)

const snd byte = 's'

func main() {
	flag.Parse()

	bytes, err := ioutil.ReadFile(*inputFile)
	if err != nil {
		fmt.Printf("Could not open file %s because %v.\n", *inputFile, err)
	}
	contents := string(bytes)
	instructions := strings.Split(contents[:len(contents)-1], "\n")

	m := Machine{0, make([]Instruction, len(instructions)), make(map[byte]int)}
	for idx, inst := range instructions {
		var f Instruction

		firstOperand := inst[4]
		var secondOperand byte
		if len(inst) >= 7 {
			secondOperand = inst[6]
		}

		switch inst[0:3] {
		case "snd":
			if lit, err := strconv.Atoi(inst[4:]); err != nil {
				f = func(r map[byte]int, i *int) {
					r[snd] = r[firstOperand]
					fmt.Printf("Played %d\n", r[snd])
					*i++
				}
			} else {
				f = func(r map[byte]int, i *int) {
					r[snd] = lit
					fmt.Printf("Played %d\n", r[snd])
					*i++
				}
			}
		case "set":
			operands := strings.Split(inst[4:], " ")
			if lit, err := strconv.Atoi(operands[1]); err != nil {
				f = func(r map[byte]int, i *int) {
					fmt.Printf("Setting '%c' to value of '%c'\n", operands[0][0], secondOperand)
					r[operands[0][0]] = r[secondOperand]
					*i++
				}
			} else {
				f = func(r map[byte]int, i *int) {
					fmt.Printf("Setting '%c' to %d\n", operands[0][0], lit)
					r[operands[0][0]] = lit
					*i++
				}
			}
		case "add":
			operands := strings.Split(inst[4:], " ")
			if lit, err := strconv.Atoi(operands[1]); err != nil {
				f = func(r map[byte]int, i *int) {
					r[operands[0][0]] += r[secondOperand]
					*i++
				}
			} else {
				f = func(r map[byte]int, i *int) {
					r[operands[0][0]] += lit
					*i++
				}
			}
		case "mul":
			operands := strings.Split(inst[4:], " ")
			if lit, err := strconv.Atoi(operands[1]); err != nil {
				f = func(r map[byte]int, i *int) {
					r[operands[0][0]] *= r[secondOperand]
					*i++
				}
			} else {
				f = func(r map[byte]int, i *int) {
					r[operands[0][0]] *= lit
					*i++
				}
			}
		case "mod":
			operands := strings.Split(inst[4:], " ")
			if lit, err := strconv.Atoi(operands[1]); err != nil {
				f = func(r map[byte]int, i *int) {
					r[operands[0][0]] %= r[secondOperand]
					*i++
				}
			} else {
				f = func(r map[byte]int, i *int) {
					r[operands[0][0]] %= lit
					*i++
				}
			}
		case "rcv":
			f = func(r map[byte]int, i *int) {
				if r[firstOperand] != 0 {
					r[firstOperand] = r[snd]
					fmt.Printf("Recovered %d\n", r[snd])
				}
				*i++
			}
		case "jgz":
			operands := strings.Split(inst[4:], " ")
			if comparator, err := strconv.Atoi(operands[0]); err != nil {
				// Just use the 0 letter of the 0th operand
				if lit, err := strconv.Atoi(operands[1]); err != nil {
					f = func(r map[byte]int, i *int) {
						offset := r[operands[1][0]]
						if r[operands[0][0]] > 0 {
							*i += offset
						} else {
							*i++
						}
					}
				} else {
					f = func(r map[byte]int, i *int) {
						offset := lit
						if r[operands[0][0]] > 0 {
							*i += offset
						} else {
							*i++
						}
					}
				}
			} else {
				// Use the parsed comparator literal value.
				if lit, err := strconv.Atoi(operands[1]); err != nil {
					f = func(r map[byte]int, i *int) {
						offset := r[operands[1][0]]
						if comparator > 0 {
							*i += offset
						} else {
							*i++
						}
					}
				} else {
					f = func(r map[byte]int, i *int) {
						offset := lit
						if comparator > 0 {
							*i += offset
						} else {
							*i++
						}
					}
				}
			}
		default:
			fmt.Printf("Unrecognized instruction %s\n", inst)
			return
		}
		m.Program[idx] = f
	}

	m.Run()
}

func (m *Machine) Run() {
	for {
		if m.Ip >= len(m.Program) || m.Ip < 0 {
			fmt.Printf("Out of bounds, terminating.\n")
			break
		}
		fmt.Printf("Executing IP: %d; 'a': %d, 'b': %d, 'p': %d\n", m.Ip, m.Registers['a'], m.Registers['b'], m.Registers['p'])
		m.Program[m.Ip](m.Registers, &m.Ip)
	}
}
