package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

var inputFile = flag.String("inputFile", "inputs/day05.input", "Relative file path to use as input.")
var partB = flag.Bool("partB", true, "Whether to use part B logic.")

type Inst int

func main() {
	flag.Parse()

	bytes, err := ioutil.ReadFile(*inputFile)
	if err != nil {
		fmt.Printf("Could not open file %s because %v.\n", *inputFile, err)
	}
	contents := string(bytes)
	lines := strings.Split(contents, "\n")
	l := make([]Inst, len(lines)-1)

	for i, line := range lines {
		if len(line) == 0 {
			break
		}
		n, err := strconv.Atoi(line)
		if err != nil {
			fmt.Printf("Could not parse %s because %v.\n", line, err)
			return
		}
		l[i] = Inst(n)
	}

	steps := 0
	for pos := 0; pos >= 0 && pos < len(l); steps++ {
		offset := l[pos]
		l[pos] += Inst(1)
		pos += int(offset)
	}

	fmt.Printf("Result is %d\n", steps)
}
