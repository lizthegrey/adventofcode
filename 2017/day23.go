package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

var inputFile = flag.String("inputFile", "inputs/day23.input", "Relative file path to use as input.")
var partB = flag.Bool("partB", false, "Use Part B logic.")

type Machine struct {
	Ip        int
	Id        int
	Program   []Instruction
	Registers map[byte]int
}

type Instruction func(*Machine)

func main() {
	flag.Parse()

	bytes, err := ioutil.ReadFile(*inputFile)
	if err != nil {
		fmt.Printf("Could not open file %s because %v.\n", *inputFile, err)
	}
	contents := string(bytes)
	instructions := strings.Split(contents[:len(contents)-1], "\n")

	m := Machine{0, 0, make([]Instruction, len(instructions)), make(map[byte]int)}
	for idx, inst := range instructions {
		var f Instruction

		var secondOperand byte
		if len(inst) >= 7 {
			secondOperand = inst[6]
		}

		switch inst[0:3] {
		case "set":
			operands := strings.Split(inst[4:], " ")
			if lit, err := strconv.Atoi(operands[1]); err != nil {
				f = func(m *Machine) {
					r := m.Registers
					i := &m.Ip
					r[operands[0][0]] = r[secondOperand]
					*i++
				}
			} else {
				f = func(m *Machine) {
					r := m.Registers
					i := &m.Ip
					r[operands[0][0]] = lit
					*i++
				}
			}
		case "sub":
			operands := strings.Split(inst[4:], " ")
			if lit, err := strconv.Atoi(operands[1]); err != nil {
				f = func(m *Machine) {
					r := m.Registers
					i := &m.Ip
					r[operands[0][0]] -= r[secondOperand]
					*i++
				}
			} else {
				f = func(m *Machine) {
					r := m.Registers
					i := &m.Ip
					r[operands[0][0]] -= lit
					*i++
				}
			}
		case "mul":
			operands := strings.Split(inst[4:], " ")
			if lit, err := strconv.Atoi(operands[1]); err != nil {
				f = func(m *Machine) {
					r := m.Registers
					i := &m.Ip
					r['m']++
					r[operands[0][0]] *= r[secondOperand]
					*i++
				}
			} else {
				f = func(m *Machine) {
					r := m.Registers
					i := &m.Ip
					r['m']++
					r[operands[0][0]] *= lit
					*i++
				}
			}
		case "jnz":
			operands := strings.Split(inst[4:], " ")
			if comparator, err := strconv.Atoi(operands[0]); err != nil {
				// Just use the 0 letter of the 0th operand
				if lit, err := strconv.Atoi(operands[1]); err != nil {
					f = func(m *Machine) {
						r := m.Registers
						i := &m.Ip
						offset := r[operands[1][0]]
						if r[operands[0][0]] != 0 {
							*i += offset
						} else {
							*i++
						}
					}
				} else {
					f = func(m *Machine) {
						r := m.Registers
						i := &m.Ip
						offset := lit
						if r[operands[0][0]] != 0 {
							*i += offset
						} else {
							*i++
						}
					}
				}
			} else {
				// Use the parsed comparator literal value.
				if lit, err := strconv.Atoi(operands[1]); err != nil {
					f = func(m *Machine) {
						r := m.Registers
						i := &m.Ip
						offset := r[operands[1][0]]
						if comparator != 0 {
							*i += offset
						} else {
							*i++
						}
					}
				} else {
					f = func(m *Machine) {
						i := &m.Ip
						offset := lit
						if comparator != 0 {
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

	if *partB {
		m.Registers['a'] = 1
	}

	m.Run()
	fmt.Printf("Mul was called %d times.\n", m.Registers['m'])
}

func (m *Machine) Run() {
	for {
		if m.Ip >= len(m.Program) || m.Ip < 0 {
			fmt.Printf("Out of bounds, terminating.\n")
			break
		}
		m.Program[m.Ip](m)
	}
}
