package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strings"
)

var inputFile = flag.String("inputFile", "inputs/day20.input", "Relative file path to use as input.")

type Coord struct {
	R, C int
}

type Grid map[Coord]bool

func (g Grid) Count() int {
	var set int
	for _, v := range g {
		if v {
			set++
		}
	}
	return set
}

func (g Grid) Tick(key [512]bool, i int) Grid {
	ret := make(Grid)
	var defaultFill bool

	if key[0] && !key[511] {
		defaultFill = (i%2 == 0)
	} else {
		defaultFill = false
	}

	min := -2 * i
	max := 100 + 2*i
	for r := min; r < max; r++ {
		for c := min; c < max; c++ {
			var lookup int
			for rOffset := -1; rOffset <= 1; rOffset++ {
				for cOffset := -1; cOffset <= 1; cOffset++ {
					lookup = lookup << 1
					if val, ok := g[Coord{r + rOffset, c + cOffset}]; ok && val {
						lookup++
					} else if !ok && defaultFill {
						lookup++
					}
				}
			}
			ret[Coord{r, c}] = key[lookup]
		}
	}
	return ret
}

func main() {
	flag.Parse()

	bytes, err := ioutil.ReadFile(*inputFile)
	if err != nil {
		return
	}
	contents := string(bytes)
	split := strings.Split(contents, "\n")
	split = split[:len(split)-1]

	var key [512]bool
	for i, c := range split[0] {
		switch c {
		case '.':
			// Pass
		case '#':
			key[i] = true
		default:
			fmt.Println("invalid character.")
			return
		}
	}

	values := make(Grid)
	for r, l := range split[2:] {
		for c, v := range l {
			switch v {
			case '.':
				values[Coord{r, c}] = false
			case '#':
				values[Coord{r, c}] = true
			default:
				fmt.Println("invalid character.")
				return
			}
		}
	}
	for i := 1; i <= 50; i++ {
		values = values.Tick(key, i)
		if i == 2 {
			fmt.Println(values.Count())
		}
	}
	fmt.Println(values.Count())
}
