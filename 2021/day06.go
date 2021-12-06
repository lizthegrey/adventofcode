package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

var inputFile = flag.String("inputFile", "inputs/day06.input", "Relative file path to use as input.")

type FishClock map[int]int

func main() {
	flag.Parse()
	bytes, err := ioutil.ReadFile(*inputFile)
	if err != nil {
		return
	}
	contents := string(bytes)
	split := strings.Split(strings.Split(contents, "\n")[0], ",")

	fish := make(FishClock)
	for _, s := range split {
		n, _ := strconv.Atoi(s)
		fish[n]++
	}
	fmt.Println(fish.iter(80))
	fmt.Println(fish.iter(256))
}

func (fish FishClock) iter(maxDays int) int {
	for day := 1; day <= maxDays; day++ {
		oldFish := fish
		fish = make(FishClock)
		for k, v := range oldFish {
			if k == 0 {
				fish[6] += v
				fish[8] += v
				continue
			}
			fish[k-1] += v
		}
	}

	count := 0
	for _, v := range fish {
		count += v
	}
	return count
}
