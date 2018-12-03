package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strconv"
)

var inputFile = flag.String("inputFile", "inputs/day03.input", "Relative file path to use as input.")

var reg = regexp.MustCompile("#([0-9]+) @ ([0-9]+),([0-9]+): ([0-9]+)x([0-9]+)")

type Claim struct {
	top, bottom, left, right int // inclusive
}

type Coord struct {
	X, Y int
}

func (c1 Claim) Contains(c2 Claim) *Claim {
	var result Claim
	if c1.top <= c2.top {
		result.top = c2.top
	} else {
		result.top = c1.top
	}
	if c1.bottom <= c2.bottom {
		result.bottom = c1.bottom
	} else {
		result.bottom = c2.bottom
	}
	if c1.left <= c2.left {
		result.left = c2.left
	} else {
		result.left = c1.left
	}
	if c1.right <= c2.right {
		result.right = c1.right
	} else {
		result.right = c2.right
	}
	if result.right < result.left || result.bottom < result.top {
		return nil
	}
	return &result
}

func main() {
	flag.Parse()
	f, err := os.Open(*inputFile)
	if err != nil {
		return
	}
	defer f.Close()

	claims := make(map[int]Claim)

	r := bufio.NewReader(f)
	for {
		l, err := r.ReadString('\n')
		if err != nil || len(l) == 0 {
			break
		}
		l = l[:len(l)-1]
		parsed := reg.FindStringSubmatch(l)
		claim, err := strconv.Atoi(parsed[1])
		offsetX, err := strconv.Atoi(parsed[2])
		offsetY, err := strconv.Atoi(parsed[3])
		sizeX, err := strconv.Atoi(parsed[4])
		sizeY, err := strconv.Atoi(parsed[5])
		claims[claim] = Claim{offsetY, offsetY + sizeY - 1, offsetX, offsetX + sizeX - 1}
	}

	conflicts := make(map[Claim]bool)
	conflicted := make(map[int]bool)
	for k1, v1 := range claims {
		for k2, v2 := range claims {
			if k1 == k2 {
				continue
			}
			conflict := v1.Contains(v2)
			if conflict != nil {
				conflicts[*conflict] = true
				conflicted[k1] = true
				conflicted[k2] = true
			}
		}
	}

	overlaps := make(map[Coord]bool)
	for k := range conflicts {
		for x := k.left; x <= k.right; x++ {
			for y := k.top; y <= k.bottom; y++ {
				overlaps[Coord{x, y}] = true
			}
		}
	}

	fmt.Printf("Result is %d\n", len(overlaps))
	for k := range claims {
		if !conflicted[k] {
			fmt.Printf("Not conflicted: %d\n", k)
		}
	}
}
