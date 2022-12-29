package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

var inputFile = flag.String("inputFile", "inputs/day21.input", "Relative file path to use as input.")

type channels map[string]chan int

func (c channels) get(key string) chan int {
	if ret, ok := c[key]; ok {
		return ret
	} else {
		ret = make(chan int, 100)
		c[key] = ret
		return ret
	}
}

func (c channels) start(key string, f func() int) {
	v := f()
	for i := 0; i < 100; i++ {
		c.get(key) <- v
	}
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
	// Cheat using goroutines.
	monkeys := make(channels)
	for _, s := range split[:len(split)-1] {
		parts := strings.Split(s, ": ")
		key := parts[0]

		// Pre-populate the map
		monkeys.get(key)

		command := strings.Split(parts[1], " ")
		if len(command) == 1 {
			val, _ := strconv.Atoi(command[0])
			go monkeys.start(key, func() int {
				return val
			})
		} else {
			x := command[0]
			op := command[1]
			y := command[2]

			// Pre-populate the map
			monkeys.get(x)
			monkeys.get(y)

			go monkeys.start(key, func() int {
				a := <-monkeys.get(x)
				b := <-monkeys.get(y)
				switch op {
				case "+":
					return a + b
				case "-":
					return a - b
				case "*":
					return a * b
				case "/":
					return a / b
				}
				return -1
			})
		}
	}
	fmt.Println(<-monkeys["root"])
}
