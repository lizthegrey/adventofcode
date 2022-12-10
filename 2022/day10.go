package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

var inputFile = flag.String("inputFile", "inputs/day10.input", "Relative file path to use as input.")

func main() {
	flag.Parse()
	bytes, err := ioutil.ReadFile(*inputFile)
	if err != nil {
		return
	}
	contents := string(bytes)
	split := strings.Split(contents, "\n")

	var strength int
	var pixels [240]bool
	x := 1
	pc := 1
	for _, s := range split[:len(split)-1] {
		if pc%40 == 20 {
			strength += x * pc
		}

		if x >= ((pc-1)%40)-1 && x <= ((pc-1)%40)+1 {
			pixels[pc-1] = true
		}

		instr := strings.Split(s, " ")
		switch instr[0] {
		case "addx":
			val, _ := strconv.Atoi(instr[1])
			pc++
			if pc%40 == 20 {
				strength += x * pc
			}
			if x >= ((pc-1)%40)-1 && x <= ((pc-1)%40)+1 {
				pixels[pc-1] = true
			}
			pc++
			x += val
		case "noop":
			pc++
		}
	}

	// part A
	fmt.Println(strength)

	// part B
	for i, v := range pixels {
		if v {
			fmt.Printf("#")
		} else {
			fmt.Printf(".")
		}
		if i%40 == 39 {
			fmt.Println()
		}
	}
}
