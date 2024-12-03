package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"
)

var inputFile = flag.String("inputFile", "inputs/day03.input", "Relative file path to use as input.")

type row []int

var pattern = regexp.MustCompile(`mul\(([0-9]+),([0-9]+)\)`)
var enable = regexp.MustCompile(`do\(\)`)
var disable = regexp.MustCompile(`don't\(\)`)

type valid struct {
	pos, val int
}

func main() {
	flag.Parse()
	bytes, err := ioutil.ReadFile(*inputFile)
	if err != nil {
		return
	}
	contents := string(bytes)
	split := strings.Split(contents, "\n")

	// part A
	var sum int
	for _, s := range split[:len(split)-1] {
		results := pattern.FindAllStringSubmatch(s, -1)
		for _, r := range results {
			a, _ := strconv.Atoi(r[1])
			b, _ := strconv.Atoi(r[2])
			sum += a * b
		}
	}
	fmt.Println(sum)

	// part B
	sum = 0
	enabled := true
	for _, s := range split[:len(split)-1] {
		enables := enable.FindAllStringIndex(s, -1)
		disables := disable.FindAllStringIndex(s, -1)
		results := pattern.FindAllStringSubmatchIndex(s, -1)
		var ops []valid
		for _, r := range results {
			a, _ := strconv.Atoi(s[r[2]:r[3]])
			b, _ := strconv.Atoi(s[r[4]:r[5]])
			ops = append(ops, valid{r[0], a * b})
		}
		var eIdx, dIdx, oIdx int
		for i := 0; i < len(s); i++ {
			if eIdx < len(enables) && enables[eIdx][0] == i {
				enabled = true
				eIdx += 1
			}
			if dIdx < len(disables) && disables[dIdx][0] == i {
				enabled = false
				dIdx += 1
			}
			if oIdx < len(ops) && ops[oIdx].pos == i {
				if enabled {
					sum += ops[oIdx].val
				}
				oIdx += 1
			}
		}
	}
	fmt.Println(sum)
}
