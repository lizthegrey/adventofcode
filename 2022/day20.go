package main

import (
	"container/ring"
	"flag"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

var inputFile = flag.String("inputFile", "inputs/day20.input", "Relative file path to use as input.")

func main() {
	flag.Parse()
	bytes, err := ioutil.ReadFile(*inputFile)
	if err != nil {
		return
	}
	contents := string(bytes)
	input := strings.Split(contents, "\n")
	input = input[:len(input)-1]
	ringSize := len(input)

	valIndex := make(map[int]*ring.Ring)
	seqIndex := make([]*ring.Ring, 0, ringSize)

	message := ring.New(ringSize)
	for _, c := range input {
		n, err := strconv.Atoi(string(c))
		if err != nil {
			fmt.Printf("Failed to parse %s\n", input)
			break
		}
		message.Value = n
		valIndex[n] = message
		seqIndex = append(seqIndex, message)
		message = message.Next()
	}

	// Loop once around through the _original_ sequence.
	for _, cur := range seqIndex {
		val := cur.Value.(int)
		if val == 0 {
			continue
		}
		prev := cur.Move(-1)
		prev.Unlink(1)
		prev.Move(val).Link(cur)
	}

	var total int
	cur := valIndex[0]
	for i := 1; i <= 3*1000; i++ {
		cur = cur.Next()
		if i%1000 == 0 {
			total += cur.Value.(int)
		}
	}
	fmt.Println(total)
}
