package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

var inputFile = flag.String("inputFile", "inputs/day07.input", "Relative file path to use as input.")

type CrabRave map[int]int

func main() {
	flag.Parse()
	bytes, err := ioutil.ReadFile(*inputFile)
	if err != nil {
		return
	}
	contents := string(bytes)
	split := strings.Split(strings.Split(contents, "\n")[0], ",")

	lowest := 10000000000
	highest := -1
	crabs := make(CrabRave)
	for _, s := range split {
		n, _ := strconv.Atoi(s)
		crabs[n]++
		if n > highest {
			highest = n
		}
		if n < lowest {
			lowest = n
		}
	}

	lowestFuel := len(split) * 1000000
	lowestFuelB := len(split) * 1000000
	for target := lowest; target <= highest; target++ {
		fuel := crabs.equalize(target)
		fuelB := crabs.equalizeB(target)
		if fuel < lowestFuel {
			lowestFuel = fuel
		}
		if fuelB < lowestFuelB {
			lowestFuelB = fuelB
		}
	}

	fmt.Println(lowestFuel)
	fmt.Println(lowestFuelB)
}

func (c CrabRave) equalize(pos int) int {
	fuel := 0
	for k, v := range c {
		if k > pos {
			fuel += v * (k - pos)
		} else {
			fuel += v * (pos - k)
		}
	}
	return fuel
}

func (c CrabRave) equalizeB(pos int) int {
	fuel := 0
	for k, v := range c {
		if k > pos {
			fuel += v * (k - pos) * (1 + k - pos) / 2
		} else {
			fuel += v * (pos - k) * (1 + pos - k) / 2
		}
	}
	return fuel
}
