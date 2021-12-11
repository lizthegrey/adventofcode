package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strings"
)

var inputFile = flag.String("inputFile", "inputs/day11.input", "Relative file path to use as input.")

type OctoArray [10][10]int
type FlashArray [10][10]bool

func (o *OctoArray) increment(r, c int, flashed FlashArray) {
	if r < 0 || c < 0 || r >= 10 || c >= 10 || flashed[r][c] {
		return
	}
	o[r][c]++
}

func (o *OctoArray) iterate() int {
	for r, row := range o {
		for c, _ := range row {
			o[r][c]++
		}
	}
	flashes := 0
	var flashed FlashArray
	for {
		done := true
		for r, row := range o {
			for c, v := range row {
				if !flashed[r][c] && v > 9 {
					done = false
					o[r][c] = 0
					flashed[r][c] = true
					flashes++

					for neighborR := r - 1; neighborR <= r+1; neighborR++ {
						for neighborC := c - 1; neighborC <= c+1; neighborC++ {
							if neighborR == r && neighborC == c {
								continue
							}
							o.increment(neighborR, neighborC, flashed)
						}
					}
				}
			}
		}
		if done {
			break
		}
	}
	return flashes
}

func main() {
	flag.Parse()
	bytes, err := ioutil.ReadFile(*inputFile)
	if err != nil {
		return
	}
	contents := string(bytes)
	split := strings.Split(contents, "\n")
	split = split[:len(split)-1]

	var octopodes OctoArray
	for r, l := range split {
		for c, v := range l {
			octopodes[r][c] = int(v - '0')
		}
	}

	flashes := 0
	for i := 1; true; i++ {
		current := octopodes.iterate()
		flashes += current
		if i == 100 {
			fmt.Println(flashes)
		}
		if current == 100 {
			fmt.Println(i)
			break
		}
	}
}
