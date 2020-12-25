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

	dLoop := Loop(7, door)
	cValue := 1
	for i := 1; i <= dLoop; i++ {
		cValue *= card
		cValue %= 20201227
	}
	cLoop := Loop(7, card)
	dValue := 1
	for i := 1; i <= cLoop; i++ {
		dValue *= door
		dValue %= 20201227
	}
	if dValue != cValue {
		fmt.Printf("%d != %d", dValue, cValue)
	}
	fmt.Println(dValue)
}

func Loop(subject, target int) int {
	value := 1
	for i := 1; ; i++ {
		value *= subject
		value %= 20201227
		if value == target {
			return i
		}
	}
	return -1
}
