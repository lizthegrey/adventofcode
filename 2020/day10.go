package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"sort"
	"strconv"
	"strings"
)

var inputFile = flag.String("inputFile", "inputs/day10.input", "Relative file path to use as input.")

func main() {
	flag.Parse()
	bytes, err := ioutil.ReadFile(*inputFile)
	if err != nil {
		return
	}
	contents := string(bytes)
	split := strings.Split(contents, "\n")
	split = split[:len(split)-1]
	voltages := make([]int, len(split))
	for i, s := range split {
		n, err := strconv.Atoi(s)
		if err != nil {
			fmt.Printf("Failed to parse %s\n", s)
			break
		}
		voltages[i] = n
	}
	sort.Ints(voltages)
	var prev, ones, threes int

	exists := make(map[int]int)
	for i, v := range voltages {
		exists[v] = i
		diff := v - prev
		switch diff {
		case 1:
			ones++
		case 3:
			threes++
		}
		prev = v
	}
	threes++
	fmt.Println(ones * threes)

	// (1)   (1)   (1)    ways(N-1) + ways(N)  w(N-2)+w(N-1)+w(N) etc
	// N+3 -> N -> N-1 -> N-2      ->          N-3
	// ............111
	// ........... 321
	// ...... ....6631
	// 0..............
	// (except in 3d, not in 2d)
	ways := make([]int, len(voltages))
	ways[len(voltages)-1] = 1
	for i := len(voltages) - 2; i >= 0; i-- {
		sum := 0
		for diff := 1; diff <= 3; diff++ {
			if pos, ok := exists[voltages[i]+diff]; ok {
				sum += ways[pos]
			}
		}
		ways[i] = sum
	}
	ret := 0
	for v := 1; v <= 3; v++ {
		if pos, ok := exists[v]; ok {
			ret += ways[pos]
		}
	}
	fmt.Println(ret)
}
