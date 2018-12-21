package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var inputFile = flag.String("inputFile", "inputs/day21.input", "Relative file path to use as input.")

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
	cycles := 0
	winners := make(map[int]int)
	foundFirstWinner := false

	for r[ipreg] >= 0 && r[ipreg] < len(instructions) {
		cycles++
		if r[ipreg] == 28 {
			// This line tests for equality against our input.
			// Find out what value would have matched, then pretend we didn't match.
			if !foundFirstWinner {
				foundFirstWinner = true
				fmt.Printf("Smallest matching input is %d.\n", r[2])
			}
			if winners[r[2]] == 0 {
				winners[r[2]] = cycles
			} else {
				// We've repeated and can stop.
				break
			}
			r[ipreg] += 2
		}
		if r[ipreg] == 17 {
			// Simulate doing the expensive divide operation to run faster.
			r[5] /= 256
			cycles += 9 + 7*r[5]
			r[ipreg] = 8
		}

		i := instructions[r[ipreg]]
		r = i.F(r, i.Operands[0], i.Operands[1], i.Operands[2])
		r[ipreg]++
	}

	fmt.Printf("Found a loop after %d cycles; stopping.\n", cycles)

	winningInput := -1
	highest := 0
	for k, v := range winners {
		if v > highest {
			winningInput = k
			highest = v
		}
	}
	fmt.Printf("Largest matching input is %d.\n", winningInput)
}
