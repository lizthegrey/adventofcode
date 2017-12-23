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
	Id        int
	Program   []Instruction
	Registers map[byte]int
	SendQ     chan<- int
	RecvQ     <-chan int
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

	m := Machine{0, 0, make([]Instruction, len(instructions)), make(map[byte]int), nil, nil}
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
				f = func(m *Machine) {
					r := m.Registers
					i := &m.Ip
					fmt.Printf("Program %d sending value.\n", m.Id)
					m.SendQ <- r[firstOperand]
					*i++
				}
			} else {
				f = func(m *Machine) {
					i := &m.Ip
					fmt.Printf("Program %d sending value.\n", m.Id)
					m.SendQ <- lit
					*i++
				}
			}
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
		case "add":
			operands := strings.Split(inst[4:], " ")
			if lit, err := strconv.Atoi(operands[1]); err != nil {
				f = func(m *Machine) {
					r := m.Registers
					i := &m.Ip
					r[operands[0][0]] += r[secondOperand]
					*i++
				}
			} else {
				f = func(m *Machine) {
					r := m.Registers
					i := &m.Ip
					r[operands[0][0]] += lit
					*i++
				}
			}
		case "mul":
			operands := strings.Split(inst[4:], " ")
			if lit, err := strconv.Atoi(operands[1]); err != nil {
				f = func(m *Machine) {
					r := m.Registers
					i := &m.Ip
					r[operands[0][0]] *= r[secondOperand]
					*i++
				}
			} else {
				f = func(m *Machine) {
					r := m.Registers
					i := &m.Ip
					r[operands[0][0]] *= lit
					*i++
				}
			}
		case "mod":
			operands := strings.Split(inst[4:], " ")
			if lit, err := strconv.Atoi(operands[1]); err != nil {
				f = func(m *Machine) {
					r := m.Registers
					i := &m.Ip
					r[operands[0][0]] %= r[secondOperand]
					*i++
				}
			} else {
				f = func(m *Machine) {
					r := m.Registers
					i := &m.Ip
					r[operands[0][0]] %= lit
					*i++
				}
			}
		case "rcv":
			f = func(m *Machine) {
				r := m.Registers
				i := &m.Ip
				r[firstOperand] = <-m.RecvQ
				*i++
			}
		case "jgz":
			operands := strings.Split(inst[4:], " ")
			if comparator, err := strconv.Atoi(operands[0]); err != nil {
				// Just use the 0 letter of the 0th operand
				if lit, err := strconv.Atoi(operands[1]); err != nil {
					f = func(m *Machine) {
						r := m.Registers
						i := &m.Ip
						offset := r[operands[1][0]]
						if r[operands[0][0]] > 0 {
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
					f = func(m *Machine) {
						r := m.Registers
						i := &m.Ip
						offset := r[operands[1][0]]
						if comparator > 0 {
							*i += offset
						} else {
							*i++
						}
					}
				} else {
					f = func(m *Machine) {
						i := &m.Ip
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

	a := m.Clone()
	b := m.Clone()

	a.Id = 0
	b.Id = 1

	aToB := make(chan int, 1000000000)
	bToA := make(chan int, 1000000000)
	a.SendQ = aToB
	b.RecvQ = aToB
	a.RecvQ = bToA
	b.SendQ = bToA

	a.Registers['p'] = 0
	b.Registers['p'] = 1

	doneA := make(chan bool)
	doneB := make(chan bool)

	go func() {
		a.Run()
		doneA <- true
	}()
	go func() {
		b.Run()
		doneB <- true
	}()
	<-doneA
	<-doneB
}

func (m *Machine) Clone() Machine {
	r := make(map[byte]int)
	for k, v := range m.Registers {
		r[k] = v
	}
	return Machine{m.Ip, m.Id, m.Program, r, nil, nil}
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
