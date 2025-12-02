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
	contents := string(bytes)[:len(bytes)-1]

	var countA, countB int
	for _, r := range strings.Split(contents, ",") {
		parts := strings.Split(r, "-")
		lower, _ := strconv.Atoi(parts[0])
		upper, _ := strconv.Atoi(parts[1])
		for i := lower; i <= upper; i++ {
			str := strconv.Itoa(i)
		outer:
			for repeats := 2; repeats <= len(str); repeats++ {
				if len(str)%repeats != 0 {
					continue
				}
				incr := len(str) / repeats
				pattern := str[:incr]
				for j := incr; j+incr <= len(str); j += incr {
					if pattern != str[j:j+incr] {
						continue outer
					}
				}
				if repeats == 2 {
					countA += i
				}
				countB += i
				break
			}
		}
	}
	fmt.Println(countA)
	fmt.Println(countB)
}
