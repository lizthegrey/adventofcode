package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

var inputFile = flag.String("inputFile", "inputs/day08.input", "Relative file path to use as input.")

type Instruction func(*VM)

type VM struct {
	Program        []Instruction
	Accumulator    int
	Executed       map[int]bool
	ProgramCounter int
}

// Cycle returns whether we repeated, and whether to continue execution.
func (v *VM) Cycle(terminateOnRepeat bool) (bool, bool) {
	pc := v.ProgramCounter
	if v.Executed[pc] && terminateOnRepeat {
		return true, false
	}
	v.Executed[pc] = true
	v.ProgramCounter++
	if v.ProgramCounter == len(v.Program) {
		return false, false
	}
	v.Program[pc](v)

	return false, true
}

func main() {
	flag.Parse()
	bytes, err := ioutil.ReadFile(*inputFile)
	if err != nil {
		return
	}
	contents := string(bytes)
	split := strings.Split(contents, "\n")
	split = split[:len(split)-1]
	vm := parse(split)
	for {
		_, cont := vm.Cycle(true)
		if !cont {
			break
		}
	}
	fmt.Println(vm.Accumulator)

outer:
	for i := range split {
		instrs := make([]string, len(split))
		copy(instrs, split)
		parts := strings.Split(instrs[i], " ")
		substitution := []byte(instrs[i])
		switch parts[0] {
		case "nop":
			copy(substitution, "jmp")
		case "jmp":
			copy(substitution, "nop")
		default:
			// Not eligible to be flipped.
			continue outer
		}
		instrs[i] = string(substitution)
		vm = parse(instrs)
		for {
			cycle, cont := vm.Cycle(true)
			if !cont {
				if cycle {
					continue outer
				} else {
					// I'm done.
					fmt.Println(vm.Accumulator)
					break outer
				}
			}
		}
	}
}

// parse converts an slice of assembly strings into an initialized VM.
func parse(split []string) *VM {
	var vm VM
	for _, s := range split {
		parts := strings.Split(s, " ")
		if len(parts) != 2 {
			fmt.Printf("Failed to parse: %s\n", s)
			return nil
		}
		sign := parts[1][0]
		var value int
		var err error
		switch sign {
		case '-':
			value, err = strconv.Atoi(parts[1][1:])
			value *= -1
			if err != nil {
				fmt.Printf("Failed to parse %s\n", s)
				return nil
			}
		case '+':
			value, err = strconv.Atoi(parts[1][1:])
			if err != nil {
				fmt.Printf("Failed to parse %s\n", s)
				return nil
			}
		default:
			fmt.Printf("Failed to parse %s\n", s)
			return nil
		}
		var instr func(*VM)
		switch parts[0] {
		case "nop":
			instr = func(vm *VM) {}
		case "acc":
			instr = func(vm *VM) {
				vm.Accumulator += value
			}
		case "jmp":
			instr = func(vm *VM) {
				vm.ProgramCounter += value - 1
			}
		}
		vm.Program = append(vm.Program, instr)
	}
	vm.Executed = make(map[int]bool)
	return &vm
}
