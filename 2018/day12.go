package main

import (
	"bufio"
	"flag"
	"fmt"
	"math"
	"os"
	"strings"
)

var inputFile = flag.String("inputFile", "inputs/day12.input", "Relative file path to use as input.")
var generations = flag.Int("generations", 20, "Number of generations.")

type Pattern struct {
	Neg2, Neg1, Self, Pos1, Pos2 bool
}

func Display(state map[int]bool, min, max int) {
	fmt.Printf("%d  ", min)
	for i := min; i <= max; i++ {
		if state[i] {
			fmt.Printf("#")
		} else {
			fmt.Printf(".")
		}
	}
	fmt.Printf("  %d\n", max)
}

func main() {
	flag.Parse()
	f, err := os.Open(*inputFile)
	if err != nil {
		return
	}
	defer f.Close()

	pots := make(map[int]bool)
	patterns := make(map[Pattern]bool)

	r := bufio.NewReader(f)
	initial, _ := r.ReadString('\n')
	initial = initial[:len(initial)-1]
	initial = strings.Split(initial, " ")[2]
	for i, c := range initial {
		if c == '#' {
			pots[i] = true
		}
	}
	minIdxAlive := 0
	maxIdxAlive := len(initial)

	r.ReadString('\n')

	for {
		l, err := r.ReadString('\n')
		if err != nil || len(l) == 0 {
			break
		}
		l = l[:len(l)-1]
		parts := strings.Split(l, " ")
		pattern := parts[0]
		neg2 := pattern[0] == '#'
		neg1 := pattern[1] == '#'
		self := pattern[2] == '#'
		pos1 := pattern[3] == '#'
		pos2 := pattern[4] == '#'
		nextGen := parts[2] == "#"
		patterns[Pattern{neg2, neg1, self, pos1, pos2}] = nextGen
	}

	for gen := 0; gen < *generations; gen++ {
		newPots := make(map[int]bool)
		minIdxAliveOld := minIdxAlive
		maxIdxAliveOld := maxIdxAlive
		minIdxAlive = math.MaxInt32
		maxIdxAlive = math.MinInt32
		for i := minIdxAliveOld - 2; i <= maxIdxAliveOld+2; i++ {
			template := Pattern{pots[i-2], pots[i-1], pots[i], pots[i+1], pots[i+2]}
			newState := patterns[template]
			if newState {
				newPots[i] = true
			}
			if newState && i < minIdxAlive {
				minIdxAlive = i
			}
			if newState && i > maxIdxAlive {
				maxIdxAlive = i
			}
		}
		pots = newPots
	}

	Display(pots, minIdxAlive, maxIdxAlive)
	result := 0
	farFuture := 50000000000
	futureSum := 0
	for k, _ := range pots {
		futureSum += farFuture + k - *generations
		result += k
	}
	fmt.Printf("Result is %d; far future is %d\n", result, futureSum)
}
