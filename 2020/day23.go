package main

import (
	"container/ring"
	"flag"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

var inputFile = flag.String("inputFile", "inputs/day23.input", "Relative file path to use as input.")
var partB = flag.Bool("partB", false, "Whether to use part B logic.")

func main() {
	flag.Parse()
	bytes, err := ioutil.ReadFile(*inputFile)
	if err != nil {
		return
	}
	contents := string(bytes)
	input := strings.Split(contents, "\n")[0]

	cache := make(map[int]*ring.Ring)

	cups := ring.New(len(input))
	for _, c := range input {
		n, err := strconv.Atoi(string(c))
		if err != nil {
			fmt.Printf("Failed to parse %s\n", input)
			break
		}
		cups.Value = n
		cache[n] = cups
		cups = cups.Next()
	}
	if *partB {
		last := cups.Prev()
		added := ring.New(1000000 - len(input))
		for i := len(input) + 1; i <= 1000000; i++ {
			added.Value = i
			cache[i] = added
			added = added.Next()
		}
		last.Link(added)
	}

	ringSize := cups.Len()

	iters := 100
	if *partB {
		iters = 10000000
	}

	current := cups
	for i := 0; i < iters; i++ {
		removed := current.Unlink(3)
		dst := 1 + ((ringSize + current.Value.(int) - 2) % ringSize)
		inRemoved := make(map[int]bool)
		for n := 1; n <= 3; n++ {
			inRemoved[removed.Move(n).Value.(int)] = true
		}
		for inRemoved[dst] {
			dst = 1 + ((ringSize + dst - 2) % ringSize)
		}
		cache[dst].Link(removed)
		current = current.Next()
	}

	first := cache[1]
	if *partB {
		a := first.Move(1).Value.(int)
		b := first.Move(2).Value.(int)
		fmt.Println(a * b)
		return
	}

	// cups should be on the cup labeled 1.
	for i := 1; i < len(input); i++ {
		first = first.Next()
		fmt.Printf("%d", first.Value.(int))
	}
	fmt.Println()
}
