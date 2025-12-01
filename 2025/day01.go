package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

var inputFile = flag.String("inputFile", "inputs/day01.input", "Relative file path to use as input.")

func main() {
	flag.Parse()
	bytes, err := ioutil.ReadFile(*inputFile)
	if err != nil {
		return
	}
	contents := string(bytes)
	split := strings.Split(contents, "\n")

	dial := 50
	var countA, countB int
	for _, s := range split[:len(split)-1] {
		dir := s[0]
		num := s[1:]
		i, _ := strconv.Atoi(num)
		if i >= 100 {
			countB += i / 100
		}
		i %= 100
		start := dial
		switch dir {
		case 'L':
			dial -= i
		case 'R':
			dial += i
		}
		if dial > 100 || (start != 0 && dial < 0) {
			countB++
		}
		dial %= 100
		if dial == 0 {
			countA++
			countB++
		}
		if dial < 0 {
			dial += 100
		}
	}
	fmt.Println(countA)
	fmt.Println(countB)
}
