package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
)

func main() {
	toMatch := map[string]int{
		"children":    3,
		"cats":        7,
		"samoyeds":    2,
		"pomeranians": 3,
		"akitas":      0,
		"vizslas":     0,
		"goldfish":    5,
		"trees":       3,
		"cars":        2,
		"perfumes":    1,
	}

	reader := bufio.NewReader(os.Stdin)
	r := regexp.MustCompile("Sue ([0-9]+): ((?:[a-z]+: [0-9]+(?:, )?)+)")
	r2 := regexp.MustCompile("([a-z]+): ([0-9]+)(?:, )?")

readLoop:
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		parsed := r.FindStringSubmatch(line)
		sueNum, _ := strconv.Atoi(parsed[1])

		properties := make(map[string]int)
		props := r2.FindAllStringSubmatch(parsed[2], -1)
		for i := range props {
			item := props[i][1]
			value, _ := strconv.Atoi(props[i][2])
			properties[item] = value
		}

		for k := range toMatch {
			v, ok := properties[k]

			if k == "cats" || k == "trees" {
				if ok && v <= toMatch[k] {
					continue readLoop
				}
			} else if k == "pomeranians" || k == "goldfish" {
				if ok && v >= toMatch[k] {
					continue readLoop
				}
			} else {
				if ok && v != toMatch[k] {
					continue readLoop
				}
			}
		}

		fmt.Println(sueNum)
	}
}
