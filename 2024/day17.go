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

	regs := make([]int64, 3)
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

	out := compute(regs, ops)
	if !slices.Equal(out, fast(regs[0])) {
		panic("fast algo is wrong")
	}
	for i, v := range out {
		if i > 0 {
			fmt.Printf(",")
		}
		fmt.Printf("%d", v)
	}
	fmt.Println()

	i := int64(35184370000000)

	// interesting clusters 35186300000000 35188450000000 35189890000000 35190160000000 35190560000000
	// [2 4 1 7 7 5 1 7 4 1 3 0 4 0 0 1]
	// we are looking for:
	// [2 4 1 7 7 5 1 7 4 6 0 3 5 5 3 0]
	// at 35193640000000 it looks more interesting
	// [2 4 1 7 7 5 1 7 4 4 0 1 4 0 0 1]
	// [2 4 1 7 7 5 1 7 4 5 0 1 4 0 0 1]
	// 35194710000000 they finally mostly agree!
	// [2 4 1 7 7 5 1 7 4 6 1 1 4 0 0 1]
	// no match by 35200000000000
	for {
		out = fast(i)
		if i%10000000 == 0 {
			if len(ops) > len(out) {
				fmt.Println("keep going up")
			}
			fmt.Println(i)
		}
		if slices.Equal(ops[0:9], out[0:9]) {
			fmt.Println(out)
			fmt.Println(ops)
		}
		if slices.Equal(ops, out) {
			fmt.Println(i)
			break
		}
		i++
	}
}

func compute(r, o []int64) []int64 {
	regs := slices.Clone(r)
	ops := slices.Clone(o)
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

func combo(oper int64, regs []int64) int64 {
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

func div(oper int64, regs []int64) int64 {
	return regs[0] / (1 << combo(oper, regs))
}

func fast(a int64) []int64 {
	var out []int64
	var b, c int64
	for {
		// 2,4
		b = a % 8

		// 1,7
		b ^= 7

		// 7,5
		c = a / (1 << b)

		// 1,7
		b ^= 7

		// 4,6
		b ^= c

		// 0,3
		a /= 8

		// 5,5
		out = append(out, b%8)

		// 3,0
		if a == 0 {
			break
		}
	}
	return out
}
