package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var inputFile = flag.String("inputFile", "inputs/day19.input", "Relative file path to use as input.")
var partB = flag.Bool("partB", false, "Whether to use the Part B logic.")

type Registers [6]int
type Instruction struct {
	F        Op
	Operands [3]int
}

type Op func(Registers, int, int, int) Registers

func Addr(r Registers, operA, operB, targetReg int) Registers {
	result := r
	result[targetReg] = r[operA] + r[operB]
	return result
}
func Addi(r Registers, operA, operB, targetReg int) Registers {
	result := r
	result[targetReg] = r[operA] + operB
	return result
}
func Mulr(r Registers, operA, operB, targetReg int) Registers {
	result := r
	result[targetReg] = r[operA] * r[operB]
	return result
}
func Muli(r Registers, operA, operB, targetReg int) Registers {
	result := r
	result[targetReg] = r[operA] * operB
	return result
}
func Banr(r Registers, operA, operB, targetReg int) Registers {
	result := r
	result[targetReg] = r[operA] & r[operB]
	return result
}
func Bani(r Registers, operA, operB, targetReg int) Registers {
	result := r
	result[targetReg] = r[operA] & operB
	return result
}
func Borr(r Registers, operA, operB, targetReg int) Registers {
	result := r
	result[targetReg] = r[operA] | r[operB]
	return result
}
func Bori(r Registers, operA, operB, targetReg int) Registers {
	result := r
	result[targetReg] = r[operA] | operB
	return result
}
func Setr(r Registers, operA, operB, targetReg int) Registers {
	result := r
	result[targetReg] = r[operA]
	return result
}
func Seti(r Registers, operA, operB, targetReg int) Registers {
	result := r
	result[targetReg] = operA
	return result
}
func Gtir(r Registers, operA, operB, targetReg int) Registers {
	result := r
	if operA > r[operB] {
		result[targetReg] = 1
	} else {
		result[targetReg] = 0
	}
	return result
}
func Gtri(r Registers, operA, operB, targetReg int) Registers {
	result := r
	if r[operA] > operB {
		result[targetReg] = 1
	} else {
		result[targetReg] = 0
	}
	return result
}
func Gtrr(r Registers, operA, operB, targetReg int) Registers {
	result := r
	if r[operA] > r[operB] {
		result[targetReg] = 1
	} else {
		result[targetReg] = 0
	}
	return result
}
func Eqir(r Registers, operA, operB, targetReg int) Registers {
	result := r
	if operA == r[operB] {
		result[targetReg] = 1
	} else {
		result[targetReg] = 0
	}
	return result
}
func Eqri(r Registers, operA, operB, targetReg int) Registers {
	result := r
	if r[operA] == operB {
		result[targetReg] = 1
	} else {
		result[targetReg] = 0
	}
	return result
}
func Eqrr(r Registers, operA, operB, targetReg int) Registers {
	result := r
	if r[operA] == r[operB] {
		result[targetReg] = 1
	} else {
		result[targetReg] = 0
	}
	return result
}

func main() {
	flag.Parse()
	f, err := os.Open(*inputFile)
	if err != nil {
		return
	}
	defer f.Close()

	ops := map[string]Op{
		"addr": Addr,
		"addi": Addi,
		"mulr": Mulr,
		"muli": Muli,
		"banr": Banr,
		"bani": Bani,
		"borr": Borr,
		"bori": Bori,
		"setr": Setr,
		"seti": Seti,
		"gtir": Gtir,
		"gtri": Gtri,
		"gtrr": Gtrr,
		"eqir": Eqir,
		"eqri": Eqri,
		"eqrr": Eqrr,
	}

	reader := bufio.NewReader(f)

	instructions := make([]Instruction, 0)
	ipreg := -1

	for idx := -1; ; idx++ {
		l, err := reader.ReadString('\n')
		if err != nil || len(l) == 0 {
			break
		}
		l = l[:len(l)-1]

		if idx == -1 {
			// Read the instruction pointer register number.
			ipreg, _ = strconv.Atoi(strings.Split(l, " ")[1])
			continue
		}

		rawInstr := strings.Split(l, " ")
		var instr Instruction
		instr.F = ops[rawInstr[0]]
		for i := 1; i < 4; i++ {
			instr.Operands[i-1], _ = strconv.Atoi(rawInstr[i])
		}
		instructions = append(instructions, instr)
	}

	var r Registers
	if *partB {
		r[0] = 1
	}
	for r[ipreg] >= 0 && r[ipreg] < len(instructions) {
		i := instructions[r[ipreg]]
		r = i.F(r, i.Operands[0], i.Operands[1], i.Operands[2])
		r[ipreg]++
	}

	fmt.Printf("Final value of register 0: %d\n", r[0])
}
