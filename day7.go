package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
)

type PendingOp struct {
	IVars []string
	Eval func(inputs []uint16)uint16
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	r := regexp.MustCompile("^(.+) -> ([a-z]+)$")

	simple := regexp.MustCompile("^([0-9]+|[a-z]+)$")
	binop := regexp.MustCompile("^([a-z]+|[0-9]+) (AND|OR) ([a-z]+)$")
	shift := regexp.MustCompile("^([a-z]+) (LSHIFT|RSHIFT) ([0-9]+)$")
	unary := regexp.MustCompile("^(NOT) ([a-z]+)$")

	var resolved map[string]uint16 = make(map[string]uint16)
	var pending map[string]PendingOp = make(map[string]PendingOp)

	// Part b
	resolved["b"] = 46065

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		parsed := r.FindStringSubmatch(line[0:len(line)-1])
		expr := parsed[1]
		ovar := parsed[2]

		// Part b
		if ovar == "b" {
			continue
		}

		if r := simple.FindStringSubmatch(expr); r != nil {
			if val, err := strconv.Atoi(r[1]); err == nil {
				resolved[ovar] = uint16(val)
			} else {
				pending[ovar] = PendingOp{[]string{r[1]}, func(inputs []uint16)uint16 {return inputs[0]}}
			}
		} else if r2 := binop.FindStringSubmatch(expr); r2 != nil {
			if val, err := strconv.Atoi(r2[1]); err == nil {
				lval := uint16(val)
				if r2[2] == "AND" {
					pending[ovar] = PendingOp{[]string{r2[3]}, func(inputs []uint16)uint16 {return lval & inputs[0]}}
				} else if r2[2] == "OR" {
					pending[ovar] = PendingOp{[]string{r2[3]}, func(inputs []uint16)uint16 {return lval | inputs[0]}}
				} else {
					fmt.Printf("Unable to parse expression: %s\n", expr)
				}
			} else {
				if r2[2] == "AND" {
					pending[ovar] = PendingOp{[]string{r2[1], r2[3]}, func(inputs []uint16)uint16 {return inputs[0] & inputs[1]}}
				} else if r2[2] == "OR" {
					pending[ovar] = PendingOp{[]string{r2[1], r2[3]}, func(inputs []uint16)uint16 {return inputs[0] | inputs[1]}}
				} else {
					fmt.Printf("Unable to parse expression: %s\n", expr)
				}
			}
		} else if r3 := shift.FindStringSubmatch(expr); r3 != nil {
			offset, _ := strconv.Atoi(r3[3])
			if r3[2] == "RSHIFT" {
				pending[ovar] = PendingOp{[]string{r3[1]}, func(inputs []uint16)uint16 {return inputs[0] >> uint(offset)}}
			} else if r3[2] == "LSHIFT" {
				pending[ovar] = PendingOp{[]string{r3[1]}, func(inputs []uint16)uint16 {return inputs[0] << uint(offset)}}
			} else {
				fmt.Printf("Unable to parse expression: %s\n", expr)
			}
		} else if r4 := unary.FindStringSubmatch(expr); r4 != nil {
			if r4[1] == "NOT" {
				pending[ovar] = PendingOp{[]string{r4[2]}, func(inputs []uint16)uint16 {return ^inputs[0]}}
			} else {
				fmt.Printf("Unable to parse expression: %s\n", expr)
			}
		} else {
			fmt.Printf("Unable to parse expression: %s\n", expr)
			return
		}
	}
	for ; len(pending) != 0; {
		inner:
		for key := range pending {
			inputs := make([]uint16, len(pending[key].IVars))
			ivars := pending[key].IVars
			for i := range ivars {
				if val, ok := resolved[ivars[i]]; ok {
					inputs[i] = val
				} else {
					continue inner
				}
			}
			resolved[key] = pending[key].Eval(inputs)
			delete(pending, key)
		}
	}
	fmt.Println(resolved["a"])
}
