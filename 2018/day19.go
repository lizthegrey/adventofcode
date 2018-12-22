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

var inputFile = flag.String("inputFile", "inputs/day19.input", "Relative file path to use as input.")
var partB = flag.Bool("partB", false, "Whether to use the Part B logic.")

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
	if *partB {
		r[0] = 1
	}
	for r[ipreg] >= 0 && r[ipreg] < len(instructions) {
		instructions[r[ipreg]].Run(&r)
		r[ipreg]++
	}

	fmt.Printf("Final value of register 0: %d\n", r[0])
}
