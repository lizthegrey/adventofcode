package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

var inputFile = flag.String("inputFile", "inputs/day06.input", "Relative file path to use as input.")

func main() {
	flag.Parse()
	bytes, err := ioutil.ReadFile(*inputFile)
	if err != nil {
		return
	}
	contents := string(bytes)
	split := strings.Split(contents, "\n")

	var countA, countB uint64
	var stacks [][]int
	for i, s := range split[:len(split)-1] {
		strs := strings.Fields(s)
		if i == len(split)-2 {
			for j, v := range strs {
				switch v {
				case "+":
					sol := 0
					for k := range stacks {
						sol += stacks[k][j]
					}
					countA += uint64(sol)
				case "*":
					sol := 1
					for k := range stacks {
						sol *= stacks[k][j]
					}
					countA += uint64(sol)
				}
			}
			break
		}

		stack := make([]int, len(strs))
		for j, v := range strs {
			num, _ := strconv.Atoi(v)
			stack[j] = num
		}
		stacks = append(stacks, stack)
	}
	fmt.Println(countA)

	var last uint64
	var accum func(uint64)
	for col := 0; col < len(split[len(split)-2]); col++ {
		op := split[len(split)-2][col]
		switch op {
		case '*':
			countB += last
			last = 1
			accum = func(n uint64) {
				last *= n
			}
		case '+':
			countB += last
			last = 0
			accum = func(n uint64) {
				last += n
			}
		case ' ':
		default:
			panic("invalid operation")
		}
		// Perform accumulation on everything above us.
		var n uint64
		var digits int
		for row := len(split) - 3; row >= 0; row-- {
			digit := split[row][col]
			if digit < '0' || digit > '9' {
				// skip ' ' characters.
				continue
			}
			digits++
			place := uint64(1)
			for _ = range digits - 1 {
				place *= 10
			}
			n += place * uint64(split[row][col]-'0')
		}
		if digits > 0 {
			accum(n)
		}
	}
	countB += last
	fmt.Println(countB)
}
