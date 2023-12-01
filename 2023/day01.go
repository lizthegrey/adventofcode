package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strings"
)

var inputFile = flag.String("inputFile", "inputs/day01.input", "Relative file path to use as input.")
var partB = flag.Bool("partB", false, "Whether to use part B logic.")

var digits = map[string]int{
	"zero":  0,
	"one":   1,
	"two":   2,
	"three": 3,
	"four":  4,
	"five":  5,
	"six":   6,
	"seven": 7,
	"eight": 8,
	"nine":  9,
}

func main() {
	flag.Parse()
	bytes, err := ioutil.ReadFile(*inputFile)
	if err != nil {
		return
	}
	contents := string(bytes)
	split := strings.Split(contents, "\n")

	sum := 0
	for _, s := range split[:len(split)-1] {
		first := -1
		last := -1
		for i := 0; i < len(s); i++ {
			if c := s[i]; c >= '0' && c <= '9' {
				digit := int(c - '0')
				if first == -1 {
					first = digit
				}
				last = digit
			} else if *partB {
				for str, digit := range digits {
					if i+len(str) <= len(s) && s[i:i+len(str)] == str {
						if first == -1 {
							first = digit
						}
						last = digit
						break
					}
				}
			}
		}
		sum += first*10 + last
	}
	fmt.Println(sum)
}
