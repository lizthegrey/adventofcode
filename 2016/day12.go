package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var trace = false

type Machine struct {
	A, B, C, D         int
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
	testInput := strings.Split(`cpy 41 a
inc a
inc a
dec a
jnz a 2
dec a`, "\n")
	execute(testInput, true)

	input := strings.Split(`cpy 1 a
cpy 1 b
cpy 26 d
jnz c 2
jnz 1 5
cpy 7 c
inc d
dec c
jnz c -2
cpy a c
inc a
dec b
jnz b -2
cpy c b
dec d
jnz d -6
cpy 16 c
cpy 17 d
inc a
dec d
jnz d -2
dec c
jnz c -5`, "\n")
	execute(input, true)
	execute(input, false)
}

func execute(input []string, partA bool) {
	r := regexp.MustCompile("(cpy|inc|dec|jnz) ([abcd]|[0-9-]+)(?: ([abcd]|[0-9-]+))?")

	var m Machine
	if partA {
		m = Machine{0, 0, 0, 0, 0, make([]Instruction, 0)}
	} else {
		m = Machine{0, 0, 1, 0, 0, make([]Instruction, 0)}
	}

	for _, line := range input {
		result := r.FindStringSubmatch(line)
		switch result[1] {
		case "cpy":
			var src *int
			var val int
			if result[2] == "a" {
				src = &m.A
			} else if result[2] == "b" {
				src = &m.B
			} else if result[2] == "c" {
				src = &m.C
			} else if result[2] == "d" {
				src = &m.D
			} else if result[2][0] == '-' {
				parsed, err := strconv.Atoi(result[2][1:])
				if err != nil {
					fmt.Println("Unrecognized int literal")
					continue
				}
				val = -parsed
			} else {
				parsed, err := strconv.Atoi(result[2])
				if err != nil {
					fmt.Println("Unrecognized int literal")
					continue
				}
				val = parsed
			}

			var dst *int
			if result[3] == "a" {
				dst = &m.A
			} else if result[3] == "b" {
				dst = &m.B
			} else if result[3] == "c" {
				dst = &m.C
			} else if result[3] == "d" {
				dst = &m.D
			} else {
				fmt.Println("Unrecognized register")
				continue
			}
			i := func(ip *int) {
				*ip++
				if src != nil {
					*dst = *src
				} else {
					*dst = val
				}
			}
			m.Instructions = append(m.Instructions, i)
		case "inc":
			var reg *int
			if result[2] == "a" {
				reg = &m.A
			} else if result[2] == "b" {
				reg = &m.B
			} else if result[2] == "c" {
				reg = &m.C
			} else if result[2] == "d" {
				reg = &m.D
			} else {
				fmt.Println("Unrecognized register")
				continue
			}
			i := func(ip *int) {
				*ip++
				*reg += 1
			}
			m.Instructions = append(m.Instructions, i)
		case "dec":
			var reg *int
			if result[2] == "a" {
				reg = &m.A
			} else if result[2] == "b" {
				reg = &m.B
			} else if result[2] == "c" {
				reg = &m.C
			} else if result[2] == "d" {
				reg = &m.D
			} else {
				fmt.Println("Unrecognized register")
				continue
			}
			i := func(ip *int) {
				*ip++
				*reg -= 1
			}
			m.Instructions = append(m.Instructions, i)
		case "jnz":
			var reg *int
			var val int
			if result[2] == "a" {
				reg = &m.A
			} else if result[2] == "b" {
				reg = &m.B
			} else if result[2] == "c" {
				reg = &m.C
			} else if result[2] == "d" {
				reg = &m.D
			} else if result[2][0] == '-' {
				parsed, err := strconv.Atoi(result[2][1:])
				if err != nil {
					fmt.Println("Unrecognized int literal")
					continue
				}
				val = -parsed
			} else {
				parsed, err := strconv.Atoi(result[2])
				if err != nil {
					fmt.Println("Unrecognized int literal")
					continue
				}
				val = parsed
			}

			offset := 0
			if result[3][0] == '-' {
				parsed, err := strconv.Atoi(result[3][1:])
				if err != nil {
					fmt.Println("Unrecognized jump offset")
					continue
				}
				offset = -parsed
			} else {
				parsed, err := strconv.Atoi(result[3])
				if err != nil {
					fmt.Println("Unrecognized jump offset")
					continue
				}
				offset = parsed
			}

			i := func(ip *int) {
				if (reg != nil && *reg != 0) || (reg == nil && val != 0) {
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
			fmt.Printf("IP: %d -- A: %d B: %d C: %d D: %d\n", m.InstructionPointer, m.A, m.B, m.C, m.D)
		}
	}

	fmt.Printf("Final value: %d\n", m.A)
}
