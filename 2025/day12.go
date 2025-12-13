package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

var inputFile = flag.String("inputFile", "inputs/day12.input", "Relative file path to use as input.")

func main() {
	flag.Parse()
	bytes, err := ioutil.ReadFile(*inputFile)
	if err != nil {
		return
	}
	contents := string(bytes)
	split := strings.Split(contents, "\n")

	var shapeSize [6]int
	var accum int
	for i, s := range split[:30] {
		for _, v := range s {
			if v == '#' {
				accum++
			}
		}
		if i % 5 == 4 {
			shapeSize[i/5] = accum
			accum = 0
		}
	}

	var countA int
	for _, s := range split[30:len(split)-1] {
		parts := strings.Split(s, ": ")
		size := strings.Split(parts[0], "x")
		sizeX, _ := strconv.Atoi(size[0])
		sizeY, _ := strconv.Atoi(size[1])
		total := sizeX * sizeY
		shapes := strings.Split(parts[1], " ")
		var minSize int
		for i, v := range shapes {
			count, _ := strconv.Atoi(v)
			minSize += shapeSize[i]*count
		}
		if minSize <= total {
			countA++
		}
	}

	fmt.Println(countA)
}
