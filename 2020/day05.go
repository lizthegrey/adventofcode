package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

var inputFile = flag.String("inputFile", "inputs/day05.input", "Relative file path to use as input.")

func main() {
	flag.Parse()
	bytes, err := ioutil.ReadFile(*inputFile)
	if err != nil {
		return
	}
	contents := string(bytes)
	split := strings.Split(contents, "\n")
	split = split[:len(split)-1]
	highestSeat := 0
	var seen [1024]bool
	for _, s := range split {
		// Read bits from MSB to LSB
		if len(s) != (7 + 3) {
			fmt.Printf("Invalid boarding pass length: %s\n", s)
		}
		bits := strings.ReplaceAll(s, "B", "1")
		bits = strings.ReplaceAll(bits, "R", "1")
		bits = strings.ReplaceAll(bits, "F", "0")
		bits = strings.ReplaceAll(bits, "L", "0")
		val, err := strconv.ParseInt(bits, 2, 0)
		if err != nil {
			fmt.Printf("Failed to parse %s\n", bits)
		}
		value := int(val)
		seen[value] = true
		if value > highestSeat {
			highestSeat = value
		}
	}
	fmt.Println(highestSeat)
	init := true
	for i, found := range seen {
		if init && found {
			init = false
			continue
		}
		if !init && !found {
			fmt.Println(i)
			break
		}
	}
}
