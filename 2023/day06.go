package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

var inputFile = flag.String("inputFile", "inputs/day06.input", "Relative file path to use as input.")

func main() {
	flag.Parse()
	bytes, err := ioutil.ReadFile(*inputFile)
	if err != nil {
		return
	}
	contents := string(bytes)
	split := strings.Split(contents, "\n")

	var times, distances []int
	for _, v := range strings.Split(split[0], " ") {
		if i, err := strconv.Atoi(v); err == nil {
			times = append(times, i)
		}
	}
	for _, v := range strings.Split(split[1], " ") {
		if i, err := strconv.Atoi(v); err == nil {
			distances = append(distances, i)
		}
	}

	product := 1
	for i := 0; i < len(times); i++ {
		var ways int
		time := times[i]
		distance := distances[i]
		for t := 0; t <= time; t++ {
			if (time-t)*t > distance {
				ways++
			}
		}
		product *= ways
	}

	fmt.Println(product)

	bigTime, _ := strconv.Atoi(strings.ReplaceAll(strings.Split(split[0], ": ")[1], " ", ""))
	bigDistance, _ := strconv.Atoi(strings.ReplaceAll(strings.Split(split[1], ": ")[1], " ", ""))
	var ways int
	for t := 0; t <= bigTime; t++ {
		if (bigTime-t)*t > bigDistance {
			ways++
		}
	}
	fmt.Println(ways)
}
