package main

import (
	"asm"
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var inputFile = flag.String("inputFile", "inputs/day16.input", "Relative file path to use as input.")
var verbose = flag.Bool("verbose", false, "Whether to print verbose output.")

var regDiagram = regexp.MustCompile(".*:[ ]+\\[(\\d+), (\\d+), (\\d+), (\\d+)\\]")

type Instruction [4]int

type OpRegistry map[int]asm.Op
type Ops []asm.Op

func main() {
	flag.Parse()
	f, err := os.Open(*inputFile)
	if err != nil {
		return
	}
	defer f.Close()

	ops := make(Ops, 0)
	opN := make([]string, 0)
	for k, v := range asm.AllOps {
		opN = append(opN, k)
		ops = append(ops, v)
	}

	reader := bufio.NewReader(f)
	spaces := 0

	lines := make([]string, 0)
	instructions := make([]Instruction, 0)
	for {
		l, err := reader.ReadString('\n')
		if err != nil || len(l) == 0 {
			break
		}
		l = l[:len(l)-1]

		if spaces < 3 {
			// This is the instruction figuring out code.
			if len(l) == 0 {
				spaces++
			} else {
				spaces = 0
			}
			lines = append(lines, l)
		} else {
			rawInstr := strings.Split(l, " ")
			var instr Instruction
			for i := 0; i < 4; i++ {
				instr[i], _ = strconv.Atoi(rawInstr[i])
			}
			instructions = append(instructions, instr)
		}
	}

	// Map observed opcodes to possible ops[i] indices.
	candidates := make(map[int]map[int]bool)

	multimatches := 0
	for ln := 0; ln+4 < len(lines); {
		// Parse the before.
		before := regDiagram.FindStringSubmatch(lines[ln])[1:5]
		var rb asm.Registers
		for i := 0; i < 4; i++ {
			rb[i], _ = strconv.Atoi(before[i])
		}

		// Parse the instruction.
		ln++
		rawInstr := strings.Split(lines[ln], " ")
		var instr Instruction
		for i := 0; i < 4; i++ {
			instr[i], _ = strconv.Atoi(rawInstr[i])
		}

		// Parse the after.
		ln++
		after := regDiagram.FindStringSubmatch(lines[ln])[1:5]
		var ra asm.Registers
		for i := 0; i < 4; i++ {
			ra[i], _ = strconv.Atoi(after[i])
		}

		// Skip the line after.
		ln++
		ln++

		// Run the instruction against before for each possible op.
		possible := make(map[int]bool)
		for i, v := range ops {
			result := v(rb, instr[1], instr[2], instr[3])
			// See which match to the after.
			if result == ra {
				possible[i] = true
				if *verbose {
					fmt.Printf("matched instruction %s\n", opN[i])
				}
			}
		}
		// Increment the count of matching instructions if 3 or more all match.
		if len(possible) >= 3 {
			multimatches++
		}
		if candidates[instr[0]] == nil {
			candidates[instr[0]] = possible
		} else {
			for k := range candidates[instr[0]] {
				if !possible[k] {
					delete(candidates[instr[0]], k)
				}
			}
			if len(candidates[instr[0]]) == 0 {
				fmt.Println("Created impossibility.")
			}
		}
	}

	// Use process of elimination knowing that each opcode goes to exactly one op.
	uniques := make(map[int]bool)
	for len(uniques) != len(ops) {
		for _, v := range candidates {
			if len(v) == 1 {
				for idx := range v {
					uniques[idx] = true
				}
				continue
			}
			toClear := make([]int, 0)
			for idx := range v {
				if uniques[idx] {
					// Someone else has taken us; clear.
					toClear = append(toClear, idx)
				}
			}
			for _, idx := range toClear {
				delete(v, idx)
			}
		}
	}

	// Every opcode should have one index that matches.
	opFunctions := make(map[int]asm.Op)
	for k, v := range candidates {
		for idx := range v {
			opFunctions[k] = ops[idx]
		}
	}

	var r asm.Registers
	for _, i := range instructions {
		r = opFunctions[i[0]](r, i[1], i[2], i[3])
	}

	fmt.Printf("Multimatching instruction count: %d\n", multimatches)
	fmt.Printf("Final value of register 0: %d\n", r[0])
}
