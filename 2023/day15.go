package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
)

var inputFile = flag.String("inputFile", "inputs/day15.input", "Relative file path to use as input.")

type Item struct {
	Label string
	Num   int
}

func main() {
	flag.Parse()
	bytes, err := ioutil.ReadFile(*inputFile)
	if err != nil {
		return
	}
	contents := string(bytes[:len(bytes)-1])
	steps := strings.Split(contents, ",")

	var boxes [256][]Item

	var sum int
outer:
	for _, step := range steps {
		sum += hash(step)

		if step[len(step)-1] == '-' {
			// Remove from box.
			label := step[:len(step)-1]
			h := hash(label)
			for i, v := range boxes[h] {
				if v.Label == label {
					copy(boxes[h][i:], boxes[h][i+1:])
					boxes[h] = boxes[h][:len(boxes[h])-1]
					break
				}
			}
			continue outer
		}
		parts := strings.Split(step, "=")
		if len(parts) != 2 {
			log.Fatalf("Found invalid step %s", step)
		}
		label := parts[0]
		h := hash(label)
		num, err := strconv.Atoi(parts[1])
		if err != nil {
			log.Fatalf("Failed to parse value of %s: %v", step, err)
		}
		this := Item{label, num}
		for i, v := range boxes[h] {
			if v.Label == label {
				boxes[h][i] = this
				continue outer
			}
		}
		boxes[h] = append(boxes[h], this)
	}
	fmt.Println(sum)

	var power int
	for idx, arr := range boxes {
		for slot, v := range arr {
			power += (1 + idx) * (1 + slot) * v.Num
		}
	}
	fmt.Println(power)
}

func hash(step string) int {
	var value int
	for _, c := range step {
		value += int(c)
		value *= 17
		value %= 256
	}
	return value
}
