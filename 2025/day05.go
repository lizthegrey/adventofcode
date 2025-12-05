package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

var inputFile = flag.String("inputFile", "inputs/day05.input", "Relative file path to use as input.")

type Range struct {
	Lower, Upper int
}

func main() {
	flag.Parse()
	bytes, err := ioutil.ReadFile(*inputFile)
	if err != nil {
		return
	}
	contents := string(bytes)[:len(bytes)-1]

	var countA, countB int
	var fresh []Range
	var ingredients []int
	for _, s := range strings.Split(contents, "\n") {
		if len(s) == 0 {
			continue
		}
		parts := strings.Split(s, "-")
		if len(parts) == 1 {
			num, _ := strconv.Atoi(parts[0])
			ingredients = append(ingredients, num)
		}
		if len(parts) == 2 {
			var r Range
			r.Lower, _ = strconv.Atoi(parts[0])
			r.Upper, _ = strconv.Atoi(parts[1])
			fresh = append(fresh, r)
		}
	}

	changed := true
	for changed {
		var next []Range
		changed = false
		for i, r := range fresh {
			for j, o := range fresh {
				if i == j {
					continue
				}
				if r.Lower > o.Upper || r.Upper < o.Lower {
					// no potential overlap
					continue
				}
				if r.Lower >= o.Lower && r.Upper <= o.Upper {
					// Contained entirely.
					changed = true
					break
				}
				if r.Lower <= o.Lower && r.Upper >= o.Lower && r.Upper <= o.Upper {
					fresh[j].Lower = r.Lower
					changed = true
					break
				}
				if r.Upper >= o.Upper && r.Lower >= o.Lower && r.Lower <= o.Upper {
					fresh[j].Upper = r.Upper
					changed = true
					break
				}
			}
			if changed {
				for k, v := range fresh {
					if k == i {
						continue
					}
					next = append(next, v)
				}
				fresh = next
				break
			}
		}
	}
	for _, i := range ingredients {
		for _, r := range fresh {
			if i >= r.Lower && i <= r.Upper {
				countA++
				break
			}
		}
	}
	for _, r := range fresh {
		countB += 1 + r.Upper - r.Lower
	}
	fmt.Println(countA)
	fmt.Println(countB)
}
