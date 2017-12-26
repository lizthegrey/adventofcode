package main

import (
	"flag"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var trace = flag.Bool("trace", false, "Print out each instruction as it's being executed.")
var eggs = flag.Int("eggs", 7, "The number of eggs to put in register A.")

type Machine struct {
	A, B, C, D         int
	InstructionPointer int
	Instructions       []Instruction
}

type Instruction struct {
	Op  string
	Src *int
	Val int
	Dst *int
	Off int
}

func (m *Machine) Execute() bool {
	if m.InstructionPointer < 0 || m.InstructionPointer >= len(m.Instructions) {
		return false
	}
	i := m.Instructions[m.InstructionPointer]
	switch i.Op {
	case "tgl":
		offset := *i.Dst + m.InstructionPointer
		if offset < 0 || offset >= len(m.Instructions) {
			m.InstructionPointer++
			return true
		}
		mod := &m.Instructions[offset]
		switch mod.Op {
		case "inc":
			mod.Op = "dec"
		case "tgl":
			fallthrough
		case "dec":
			mod.Op = "inc"
		case "jnz":
			mod.Op = "cpy"
		case "cpy":
			mod.Op = "jnz"
		}
		m.InstructionPointer++
	case "cpy":
		m.InstructionPointer++
		if i.Src != nil {
			*i.Dst = *i.Src
		} else if i.Dst != nil {
			*i.Dst = i.Val
		}
	case "inc":
		m.InstructionPointer++
		*i.Dst++
	case "dec":
		m.InstructionPointer++
		*i.Dst--
	case "jnz":
		if (i.Src != nil && *i.Src != 0) || (i.Src == nil && i.Val != 0) {
			if i.Dst != nil {
				m.InstructionPointer += *i.Dst
			} else {
				m.InstructionPointer += i.Off
			}
		} else {
			m.InstructionPointer++
		}
	}
	return true
}

func main() {
	flag.Parse()
	testInput := strings.Split(`cpy 2 a
tgl a
tgl a
tgl a
cpy 1 a
dec a
dec a`, "\n")
	execute(testInput, 0)

	input := strings.Split(`cpy a b
dec b
cpy a d
cpy 0 a
cpy b c
inc a
dec c
jnz c -2
dec d
jnz d -5
dec b
cpy b c
cpy c d
dec d
inc c
jnz d -2
tgl c
cpy -16 c
jnz 1 c
cpy 79 c
jnz 74 d
inc a
inc d
jnz d -2
inc c
jnz c -5`, "\n")
	execute(input, *eggs)
}

func execute(input []string, initA int) {
	r := regexp.MustCompile("(cpy|inc|dec|jnz|tgl) ([abcd]|[0-9-]+)(?: ([abcd]|[0-9-]+))?")

	var m Machine
	m = Machine{initA, 0, 0, 0, 0, make([]Instruction, 0)}

	for _, line := range input {
		result := r.FindStringSubmatch(line)
		switch result[1] {
		case "tgl":
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
			i := Instruction{"tgl", nil, 0, reg, 0}
			m.Instructions = append(m.Instructions, i)
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
			i := Instruction{"cpy", src, val, dst, 0}
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
			i := Instruction{"inc", nil, 0, reg, 0}
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
			i := Instruction{"dec", nil, 0, reg, 0}
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
			var off *int
			if result[3] == "a" {
				off = &m.A
			} else if result[3] == "b" {
				off = &m.B
			} else if result[3] == "c" {
				off = &m.C
			} else if result[3] == "d" {
				off = &m.D
			} else if result[3][0] == '-' {
				parsed, err := strconv.Atoi(result[3][1:])
				if err != nil {
					fmt.Printf("Unrecognized jump offset in line %s\n", line)
					continue
				}
				offset = -parsed
			} else {
				parsed, err := strconv.Atoi(result[3])
				if err != nil {
					fmt.Printf("Unrecognized jump offset in line %s\n", line)
					continue
				}
				offset = parsed
			}

			i := Instruction{"jnz", reg, val, off, offset}
			m.Instructions = append(m.Instructions, i)
		default:
			fmt.Println("Unrecognized opcode")
		}
	}

	for m.Execute() {
		if *trace {
			var op string
			if m.InstructionPointer >= 0 && m.InstructionPointer < len(m.Instructions) {
				op = m.Instructions[m.InstructionPointer].Op
			}
			fmt.Printf("IP: %d (%s) -- A: %d B: %d C: %d D: %d\n", m.InstructionPointer, op, m.A, m.B, m.C, m.D)
		}
	}

	fmt.Printf("Final value: %d\n", m.A)
}
