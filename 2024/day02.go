package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

var inputFile = flag.String("inputFile", "inputs/day02.input", "Relative file path to use as input.")

type row []int

func main() {
	flag.Parse()
	bytes, err := ioutil.ReadFile(*inputFile)
	if err != nil {
		return
	}
	contents := string(bytes)
	split := strings.Split(contents, "\n")

	var rows []row
	for _, s := range split[:len(split)-1] {
		var r row
		for _, num := range strings.Split(s, " ") {
			n, _ := strconv.Atoi(num)
			r = append(r, n)
		}
		rows = append(rows, r)
	}
	var safe, recovered int
outer:
	for _, r := range rows {
		if check(r) {
			safe++
			recovered++
			continue
		}
		for skip := range r {
			var copied row
			for i, v := range r {
				if i == skip {
					continue
				}
				copied = append(copied, v)
			}
			if check(copied) {
				recovered++
				continue outer
			}
		}
	}
	fmt.Println(safe)
	fmt.Println(recovered)
}

func check(r row) bool {
	var prev int
	var increasing bool
	for i, v := range r {
		if i == 0 {
			prev = v
			continue
		}
		delta := v - prev
		if i == 1 {
			if v < prev {
				increasing = false
			} else if v > prev {
				increasing = true
			} else {
				return false
			}
		} else {
			if increasing && delta < 0 || !increasing && delta > 0 {
				return false
			}
		}
		if delta == 0 {
			return false
		}
		if delta < 0 {
			delta *= -1
		}
		if delta > 3 {
			return false
		}
		prev = v
	}
	return true
}
