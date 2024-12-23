package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"slices"
	"strconv"
	"strings"
)

var inputFile = flag.String("inputFile", "inputs/day17.input", "Relative file path to use as input.")

func main() {
	flag.Parse()
	bytes, err := ioutil.ReadFile(*inputFile)
	if err != nil {
		return
	}
	contents := string(bytes)
	split := strings.Split(contents, "\n")

	var regs [3]int64
	var ops []int64
	for i, s := range split[:len(split)-1] {
		switch i {
		case 3:
			continue
		case 4:
			for j, v := range s[9:] {
				if j%2 == 1 {
					continue
				}
				ops = append(ops, int64(v-'0'))
			}
		default:
			regs[i], _ = strconv.ParseInt(s[12:], 10, 0)
		}
	}

	out := compute(ops, regs)
	for i, v := range out {
		if i > 0 {
			fmt.Printf(",")
		}
		fmt.Printf("%d", v)
	}
	fmt.Println()

	fmt.Println(*search(ops, regs, 0, 0))
}

func search(ops []int64, regs [3]int64, seed int64, digit int) *int64 {
	if digit == len(ops) {
		ret := seed
		return &ret
	}
	shift := (len(ops) - digit - 1) * 3
	for i := int64(0); i < (1 << 3); i++ {
		candidate := seed + (i << shift)
		regs[0] = candidate
		test := compute(ops, regs)
		if len(test) != len(ops) || test[len(ops)-digit-1] != ops[len(ops)-digit-1] {
			continue
		}
		ret := search(ops, regs, candidate, digit+1)
		if ret != nil {
			return ret
		}
	}
	return nil
}

func compute(o []int64, r [3]int64) []int64 {
	ops := slices.Clone(o)
	regs := r
	var out []int64
	for pc := 0; pc >= 0 && pc < len(ops); pc += 2 {
		oper := ops[pc+1]
		switch ops[pc] {
		case 0:
			regs[0] = div(oper, regs)
		case 1:
			regs[1] ^= oper
		case 2:
			regs[1] = combo(oper, regs) % 8
		case 3:
			if regs[0] == 0 {
				continue
			}
			pc = int(oper - 2)
		case 4:
			regs[1] ^= regs[2]
		case 5:
			out = append(out, combo(oper, regs)%8)
		case 6:
			regs[1] = div(oper, regs)
		case 7:
			regs[2] = div(oper, regs)
		default:
			panic("invalid opcode")
		}
	}
	return out
}

func combo(oper int64, regs [3]int64) int64 {
	switch oper {
	case 0:
		return 0
	case 1:
		return 1
	case 2:
		return 2
	case 3:
		return 3
	case 4:
		return regs[0]
	case 5:
		return regs[1]
	case 6:
		return regs[2]
	default:
		panic("invalid combo")
	}
}

func div(oper int64, regs [3]int64) int64 {
	return regs[0] / (1 << combo(oper, regs))
}
