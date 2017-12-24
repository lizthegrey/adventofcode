package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
)

var trace = false

type Machine struct {
	A, B               uint
	InstructionPointer int
	Instructions       []Instruction
}

type Instruction func(ip *int)

func (m *Machine) Execute() bool {
	if m.InstructionPointer < 0 || m.InstructionPointer >= len(m.Instructions) {
		return false
	}
	i := m.Instructions[m.InstructionPointer]
	i(&m.InstructionPointer)
	return true
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	r := regexp.MustCompile("(hlf|tpl|inc|jmp|jie|jio) ([ab]|[+-][0-9]+)(?:, ([+-][0-9]+))?")

	// Part (a)
	// m := Machine{0, 0, 0, make([]Instruction, 0)}
	// Part (b)
	m := Machine{1, 0, 0, make([]Instruction, 0)}

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		result := r.FindStringSubmatch(line)
		switch result[1] {
		case "hlf":
			var reg *uint
			if result[2] == "a" {
				reg = &m.A
			} else if result[2] == "b" {
				reg = &m.B
			} else {
				fmt.Println("Unrecognized register")
				continue
			}
			i := func(ip *int) {
				*ip++
				*reg /= 2
			}
			m.Instructions = append(m.Instructions, i)
		case "tpl":
			var reg *uint
			if result[2] == "a" {
				reg = &m.A
			} else if result[2] == "b" {
				reg = &m.B
			} else {
				fmt.Println("Unrecognized register")
				continue
			}
			i := func(ip *int) {
				*ip++
				*reg *= 3
			}
			m.Instructions = append(m.Instructions, i)
		case "inc":
			var reg *uint
			if result[2] == "a" {
				reg = &m.A
			} else if result[2] == "b" {
				reg = &m.B
			} else {
				fmt.Println("Unrecognized register")
				continue
			}
			i := func(ip *int) {
				*ip++
				*reg += 1
			}
			m.Instructions = append(m.Instructions, i)
		case "jmp":
			offset := 0
			if result[2][0] == '+' {
				parsed, err := strconv.Atoi(result[2][1:])
				if err != nil {
					fmt.Println("Unrecognized jump offset")
					continue
				}
				offset = parsed
			} else if result[2][0] == '-' {
				parsed, err := strconv.Atoi(result[2][1:])
				if err != nil {
					fmt.Println("Unrecognized jump offset")
					continue
				}
				offset = -parsed
			} else {
				fmt.Println("Unrecognized jump offset")
				continue
			}
			i := func(ip *int) {
				*ip += offset
			}
			m.Instructions = append(m.Instructions, i)
		case "jie":
			var reg *uint
			if result[2] == "a" {
				reg = &m.A
			} else if result[2] == "b" {
				reg = &m.B
			} else {
				fmt.Println("Unrecognized register")
				continue
			}

			offset := 0
			if result[3][0] == '+' {
				parsed, err := strconv.Atoi(result[3][1:])
				if err != nil {
					fmt.Println("Unrecognized jump offset")
					continue
				}
				offset = parsed
			} else if result[3][0] == '-' {
				parsed, err := strconv.Atoi(result[3][1:])
				if err != nil {
					fmt.Println("Unrecognized jump offset")
					continue
				}
				offset = -parsed
			} else {
				fmt.Println("Unrecognized jump offset")
				continue
			}

			i := func(ip *int) {
				if *reg%2 == 0 {
					*ip += offset
				} else {
					*ip++
				}
			}
			m.Instructions = append(m.Instructions, i)
		case "jio":
			var reg *uint
			if result[2] == "a" {
				reg = &m.A
			} else if result[2] == "b" {
				reg = &m.B
			} else {
				fmt.Println("Unrecognized register")
				continue
			}

			offset := 0
			if result[3][0] == '+' {
				parsed, err := strconv.Atoi(result[3][1:])
				if err != nil {
					fmt.Println("Unrecognized jump offset")
					continue
				}
				offset = parsed
			} else if result[3][0] == '-' {
				parsed, err := strconv.Atoi(result[3][1:])
				if err != nil {
					fmt.Println("Unrecognized jump offset")
					continue
				}
				offset = -parsed
			} else {
				fmt.Println("Unrecognized jump offset")
				continue
			}

			i := func(ip *int) {
				if *reg == 1 {
					*ip += offset
				} else {
					*ip++
				}
			}
			m.Instructions = append(m.Instructions, i)
		default:
			fmt.Println("Unrecognized opcode")
		}
	}

	for m.Execute() {
		if trace {
			fmt.Println(m.InstructionPointer)
		}
	}

	fmt.Printf("Final value: %d\n", m.B)
}
