package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

var inputFile = flag.String("inputFile", "inputs/day01.input", "Relative file path to use as input.")
var partB = flag.Bool("partB", false, "Whether to use the Part B logic.")

func main() {
	flag.Parse()
	bytes, err := ioutil.ReadFile(*inputFile)
	if err != nil {
		return
	}
	contents := string(bytes)
	sum := 0
	split := strings.Split(contents, "\n")
	for _, s := range split {
		if s == "" {
			continue
		}
		n, err := strconv.Atoi(s)
		if err != nil {
			fmt.Printf("Failed to parse %s\n", s)
			break
		}
		fuel := (n / 3) - 2
		if fuel > 0 {
			sum += fuel
		}
		if *partB {
			for fuel >= 1 {
				fuel = (fuel / 3) - 2
				if fuel > 0 {
					sum += fuel
				}
			}
		}
	}
	fmt.Println(sum)
}
