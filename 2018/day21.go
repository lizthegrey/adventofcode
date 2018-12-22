package main

import (
	"asm"
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var inputFile = flag.String("inputFile", "inputs/day21.input", "Relative file path to use as input.")

func main() {
	flag.Parse()
	f, err := os.Open(*inputFile)
	if err != nil {
		return
	}
	defer f.Close()

	reader := bufio.NewReader(f)

	instructions := make([]asm.Instruction, 0)
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
		var instr asm.Instruction
		instr.F = asm.AllOps[rawInstr[0]]
		for i := 1; i < 4; i++ {
			instr.Operands[i-1], _ = strconv.Atoi(rawInstr[i])
		}
		instructions = append(instructions, instr)
	}

	var r asm.Registers
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

		instructions[r[ipreg]].Run(&r)
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
