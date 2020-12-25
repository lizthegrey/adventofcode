package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strconv"
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
	split = split[:len(split)-1]
	card, err := strconv.Atoi(split[0])
	if err != nil {
		fmt.Printf("Failed to parse %s\n", split[0])
		return
	}
	door, err := strconv.Atoi(split[1])
	if err != nil {
		fmt.Printf("Failed to parse %s\n", split[1])
		return
	}

	doorLoop := Loop(7, door)
	value := 1
	for i := 1; i <= doorLoop; i++ {
		value *= card
		value %= 20201227
	}
	fmt.Println(value)
}

func Loop(subject, target int) int {
	value := 1
	for i := 1; i < 1000000000; i++ {
		value *= subject
		value %= 20201227
		if value == target {
			return i
		}
	}
	return -1
}
