package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

var inputFile = flag.String("inputFile", "inputs/day22.input", "Relative file path to use as input.")

type deltaKey [4]int8

func main() {
	flag.Parse()
	bytes, err := ioutil.ReadFile(*inputFile)
	if err != nil {
		return
	}
	contents := string(bytes)
	split := strings.Split(contents, "\n")

	var nums []uint64
	for _, s := range split[:len(split)-1] {
		n, _ := strconv.Atoi(s)
		nums = append(nums, uint64(n))
	}

	var sum uint64
	profits := make(map[deltaKey]int)
	for _, v := range nums {
		var deltas, prices []int8
		previous := int8(v % 10)
		for range 2000 {
			v = next(v)
			price := int8(v % 10)
			prices = append(prices, price)
			delta := price - previous
			deltas = append(deltas, delta)
			previous = price
		}
		sum += v

		candidates := make(map[deltaKey]int)
		for j := range len(deltas) - 4 {
			key := deltaKey{deltas[j], deltas[j+1], deltas[j+2], deltas[j+3]}
			if _, ok := candidates[key]; ok {
				continue
			}
			value := int(prices[j+3])
			candidates[key] = value
		}
		for k, v := range candidates {
			profits[k] += v
		}
	}
	fmt.Println(sum)

	var best int
	for _, v := range profits {
		if v > best {
			best = v
		}
	}
	fmt.Println(best)
}

func next(a uint64) uint64 {
	a = mix(a, a<<6)
	a = prune(a)
	a = mix(a, a>>5)
	a = prune(a)
	a = mix(a, a<<11)
	a = prune(a)
	return a
}

func mix(a, b uint64) uint64 {
	return a ^ b
}
func prune(a uint64) uint64 {
	return a % 16777216
}
