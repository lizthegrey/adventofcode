package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
)

var inputFile = flag.String("inputFile", "inputs/day02.input", "Relative file path to use as input.")

func main() {
	flag.Parse()
	f, err := os.Open(*inputFile)
	if err != nil {
		return
	}
	defer f.Close()

	r := bufio.NewReader(f)
	triples := 0
	doubles := 0
	memo := make(map[int]map[string]bool)
	for {
		l, err := r.ReadString('\n')
		if err != nil || len(l) == 0 {
			break
		}
		l = l[:len(l)-1]

		for pos := range l {
			if memo[pos] == nil {
				memo[pos] = make(map[string]bool)
			}
			trunc := l[0:pos] + l[pos+1:]
			if memo[pos][trunc] {
				fmt.Printf("Shared characters: %s\n", trunc)
				break
			}
			memo[pos][trunc] = true
		}

		seen := make(map[rune]int)
		for _, c := range l {
			seen[c] += 1
		}
		seen3 := false
		seen2 := false
		for _, v := range seen {
			if v == 3 {
				seen3 = true
			} else if v == 2 {
				seen2 = true
			}
		}
		if seen3 {
			triples += 1
		}
		if seen2 {
			doubles += 1
		}
	}

	result := triples * doubles
	fmt.Printf("Result is %d\n", result)
}
