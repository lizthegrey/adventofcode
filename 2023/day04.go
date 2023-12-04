package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
)

var inputFile = flag.String("inputFile", "inputs/day04.input", "Relative file path to use as input.")

func main() {
	flag.Parse()
	bytes, err := ioutil.ReadFile(*inputFile)
	if err != nil {
		return
	}
	contents := string(bytes)
	split := strings.Split(contents, "\n")

	memo := make(map[int][]int)
	sum := 0
	pending := make(map[int]int)
	for i, s := range split[:len(split)-1] {
		pending[i] = 1
		parts := strings.Split(s[10:], " | ")
		winners := make(map[int]bool)
		for _, v := range strings.Split(parts[0], " ") {
			if v == "" {
				continue
			}
			w, err := strconv.Atoi(v)
			if err != nil {
				log.Fatalf("Failed parsing %s: %v", v, err)
			}
			winners[w] = true
		}
		var won int
		for _, v := range strings.Split(parts[1], " ") {
			if v == "" {
				continue
			}
			w, err := strconv.Atoi(v)
			if err != nil {
				log.Fatalf("Failed parsing %s: %v", v, err)
			}
			if winners[w] {
				won++
				memo[i] = append(memo[i], i+won)
			}
		}
		var value int
		for j := 0; j < len(memo[i]); j++ {
			if value == 0 {
				value = 1
			} else {
				value *= 2
			}
		}
		sum += value
	}
	fmt.Println(sum)

	inventory := make(map[int]int)
	for len(pending) != 0 {
		next := make(map[int]int)
		for k, v := range pending {
			inventory[k] += v
			for _, a := range memo[k] {
				next[a] += v
			}
		}
		pending = next
	}
	var cards int
	for _, v := range inventory {
		cards += v
	}
	fmt.Println(cards)
}
