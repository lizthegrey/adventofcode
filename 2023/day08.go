package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"math/big"
	"strings"
)

var inputFile = flag.String("inputFile", "inputs/day08.input", "Relative file path to use as input.")

func main() {
	flag.Parse()
	bytes, err := ioutil.ReadFile(*inputFile)
	if err != nil {
		return
	}
	contents := string(bytes)
	split := strings.Split(contents, "\n")

	instrs := split[0]
	paths := make(map[string][2]string)
	for _, s := range split[2 : len(split)-1] {
		cur := s[0:3]
		left := s[7:10]
		right := s[12:15]
		paths[cur] = [2]string{left, right}
	}
	var steps int
	for loc := "AAA"; loc != "ZZZ"; steps++ {
		dir := instrs[steps%len(instrs)]
		switch dir {
		case 'L':
			loc = paths[loc][0]
		case 'R':
			loc = paths[loc][1]
		}
	}
	fmt.Println(steps)

	var locs []string
	for k := range paths {
		if k[2] == 'A' {
			locs = append(locs, k)
		}
	}
	var loop []int
	for _, loc := range locs {
		var steps int
		for loc[2] != 'Z' {
			dir := instrs[steps%len(instrs)]
			switch dir {
			case 'L':
				loc = paths[loc][0]
			case 'R':
				loc = paths[loc][1]
			}
			steps++
		}
		loop = append(loop, steps)
	}
	result := uint64(1)
	for _, v := range loop {
		result = lcm(result, uint64(v))
	}
	fmt.Println(result)
}

func lcm(a, b uint64) uint64 {
	var x, y, product, ret, gcd big.Int
	x.SetUint64(a)
	y.SetUint64(b)
	product.Mul(&x, &y)
	gcd.GCD(nil, nil, &x, &y)
	ret.Div(&product, &gcd)
	return ret.Uint64()
}
