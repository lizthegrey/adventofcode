package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
)

var inputFile = flag.String("inputFile", "inputs/day09.input", "Relative file path to use as input.")

func main() {
	flag.Parse()
	bytes, err := ioutil.ReadFile(*inputFile)
	if err != nil {
		return
	}
	contents := string(bytes)
	split := strings.Split(contents, "\n")

	var values [][]int
	for _, s := range split[:len(split)-1] {
		var series []int
		elems := strings.Split(s, " ")
		for _, str := range elems {
			v, err := strconv.Atoi(str)
			if err != nil {
				log.Fatalf("Failed parsing %s: %v", str, err)
			}
			series = append(series, v)
		}
		values = append(values, series)
	}

	var sumA, sumB int
	for _, series := range values {
		var first, last []int
		for {
			allZeroes := true
			var next []int
			first = append(first, series[0])
			last = append(last, series[len(series)-1])

			for i := 1; i < len(series); i++ {
				delta := series[i] - series[i-1]
				next = append(next, delta)
				if delta != 0 {
					allZeroes = false
				}
			}
			if allZeroes {
				break
			}
			series = next
		}
		for i := len(last) - 2; i >= 0; i-- {
			first[i] -= first[i+1]
			last[i] += last[i+1]
		}
		sumA += last[0]
		sumB += first[0]
	}
	fmt.Println(sumA)
	fmt.Println(sumB)
}
