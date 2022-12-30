package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strings"
)

var inputFile = flag.String("inputFile", "inputs/day25.input", "Relative file path to use as input.")

func main() {
	flag.Parse()
	bytes, err := ioutil.ReadFile(*inputFile)
	if err != nil {
		return
	}
	contents := string(bytes)
	split := strings.Split(contents, "\n")

	// there's only a part A. part B is completing all the other puzzles.
	var sum int
	for _, s := range split[:len(split)-1] {
		var num int
		place := 1
		for i := len(s) - 1; i >= 0; i -= 1 {
			switch s[i] {
			case '=':
				num -= 2 * place
			case '-':
				num -= place
			case '0':
				// Nothing to add or remove.
			case '1':
				num += place
			case '2':
				num += 2 * place
			}
			place *= 5
		}
		sum += num
	}
	fmt.Println(sum)
}
