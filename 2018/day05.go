package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

var inputFile = flag.String("inputFile", "inputs/day05.input", "Relative file path to use as input.")
var partB = flag.Bool("partB", false, "Whether to use the Part B logic.")

func main() {
	flag.Parse()
	f, err := os.Open(*inputFile)
	if err != nil {
		return
	}
	defer f.Close()

	r := bufio.NewReader(f)
	l, err := r.ReadString('\n')
	l = l[:len(l)-1]

	if !*partB {
		unmodified := react(l)
		result := len(unmodified)
		fmt.Printf("Result is %d\n", result)
	} else {
		shortest := len(l)
		for c := 'A'; c <= 'Z'; c++ {
			second := c + 32
			modified := strings.Replace(l, string([]rune{second}), "", -1)
			modified = strings.Replace(modified, string([]rune{c}), "", -1)
			candidate := len(react(modified))
			if candidate < shortest {
				shortest = candidate
			}
		}
		fmt.Printf("Result is %d\n", shortest)
	}
}

func react(l string) string {
	current := l
	for {
		transformed := false
		prev := ' '
		for i, c := range current {
			if c-prev == 32 || prev-c == 32 {
				transformed = true
				current = current[0:i-1] + current[i+1:]
				break
			}
			prev = c
		}
		if !transformed {
			break
		}
	}
	return current
}
