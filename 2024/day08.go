package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strings"
)

var inputFile = flag.String("inputFile", "inputs/day08.input", "Relative file path to use as input.")

type coord struct {
	r, c int
}

func (c coord) add(o coord) coord {
	return coord{int(c.r + o.r), int(c.c + o.c)}
}

func (c coord) sub(o coord) coord {
	return coord{int(c.r - o.r), int(c.c - o.c)}
}

func (c coord) bounds(o coord) bool {
	return c.r >= 0 && c.c >= 0 && c.r <= o.r && c.c <= o.c
}

func main() {
	flag.Parse()
	bytes, err := ioutil.ReadFile(*inputFile)
	if err != nil {
		return
	}
	contents := string(bytes)
	split := strings.Split(contents, "\n")

	byType := make(map[rune][]coord)
	var maxR, maxC int
	for r, s := range split[:len(split)-1] {
		maxR = r
		maxC = len(s) - 1
		for c, v := range s {
			if v == '.' {
				continue
			}
			loc := coord{r, c}
			byType[v] = append(byType[v], loc)
		}
	}

	limits := coord{maxR, maxC}
	loci := make(map[coord]bool)
	lociInf := make(map[coord]bool)
	for _, list := range byType {
		for i, x := range list {
			for j, y := range list {
				if i == j {
					break
				}
				lociInf[x] = true
				lociInf[y] = true
				delta := y.sub(x)
				first := true
				for pos := y.add(delta); pos.bounds(limits); pos = pos.add(delta) {
					if first {
						loci[pos] = true
						first = false
					}
					lociInf[pos] = true
				}
				delta = x.sub(y)
				first = true
				for pos := x.add(delta); pos.bounds(limits); pos = pos.add(delta) {
					if first {
						loci[pos] = true
						first = false
					}
					lociInf[pos] = true
				}
			}
		}
	}
	fmt.Println(len(loci))
	fmt.Println(len(lociInf))
}
