package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

var inputFile = flag.String("inputFile", "inputs/day02.input", "Relative file path to use as input.")

func main() {
	flag.Parse()
	bytes, err := ioutil.ReadFile(*inputFile)
	if err != nil {
		return
	}
	contents := string(bytes)
	split := strings.Split(contents, "\n")
	split = split[:len(split)-1]

	posX := 0
	posY := 0
	for _, s := range split {
		parts := strings.Split(s, " ")
		command := parts[0]
		num, _ := strconv.Atoi(parts[1])
		switch command {
		case "up":
			posY -= num
		case "down":
			posY += num
		case "forward":
			posX += num
		}
	}
	fmt.Println(posX * posY)

	posX = 0
	posY = 0
	aim := 0
	for _, s := range split {
		parts := strings.Split(s, " ")
		command := parts[0]
		num, _ := strconv.Atoi(parts[1])
		switch command {
		case "up":
			aim -= num
		case "down":
			aim += num
		case "forward":
			posX += num
			posY += num * aim
		}
	}
	fmt.Println(posX * posY)
}
