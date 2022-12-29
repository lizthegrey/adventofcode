package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

var inputFile = flag.String("inputFile", "inputs/day19.input", "Relative file path to use as input.")
var debug = flag.Bool("debug", false, "Whether to print debug output & progress.")

type resource int

const (
	none resource = iota - 1
	ore
	clay
	obsidian
	geode
)

func toResource(s string) resource {
	switch s {
	case "ore":
		return ore
	case "clay":
		return clay
	case "obsidian":
		return obsidian
	case "geode":
		return geode
	default:
		fmt.Printf("Invalid type of robot: %s\n", s)
		return none
	}
}

type inventory [4]int8
type recipe [4]inventory
type moves []resource

type state struct {
	bots inventory
	raw  inventory
}

func (i inventory) hash() int32 {
	return int32(i[0])<<0 + int32(i[1])<<4 + int32(i[2])<<8 + int32(i[3])<<12
}

func (s state) hash() int32 {
	return s.raw.hash() + s.bots.hash()<<16
}

func (s state) step(build resource, bp recipe) *state {
	next := s
	if build != none {
		for r, q := range bp[build] {
			if q == 0 {
				continue
			}
			proposed := next.raw[r] - q
			if proposed < 0 {
				return nil
			}
			next.raw[r] = proposed
		}
	}
	for r, q := range s.bots {
		if q == 0 {
			continue
		}
		next.raw[r] += q
	}
	if build != none {
		next.bots[build] += 1
	}
	return &next
}

// At each of the 24 turns, we can choose to let resources pile up, or build one specific bot.
// We also want to disqualify building nothing on a turn when we can build all types of bots;
// there's nothing further to save up for and we know that's suboptimal.
func (s state) children(bp recipe) []state {
	candidates := make([]state, 0, 3)
	// Propose in order of best resource first.
	for res := geode; res > none; res-- {
		next := s.step(res, bp)
		if next != nil {
			candidates = append(candidates, *next)
		}
	}
	if len(candidates) < 4 {
		// Always safe to do nothing.
		candidates = append(candidates, *s.step(none, bp))
	}
	return candidates
}

func main() {
	flag.Parse()
	bytes, err := ioutil.ReadFile(*inputFile)
	if err != nil {
		return
	}
	contents := string(bytes)
	split := strings.Split(contents, "\n")

	var recipes []recipe
	for _, s := range split[:len(split)-1] {
		parts := strings.Split(s, " Each ")
		var bp recipe
		for _, r := range parts[1:] {
			rp := strings.Split(r[:len(r)-1], " robot costs ")
			product := toResource(rp[0])
			for _, v := range strings.Split(rp[1], " and ") {
				vp := strings.Split(v, " ")
				quantity, _ := strconv.Atoi(vp[0])
				ingredient := toResource(vp[1])
				bp[product][ingredient] = int8(quantity)
			}
		}
		recipes = append(recipes, bp)
	}

	// part A
	var total int
	for i, bp := range recipes {
		best := bp.findBest(24)
		if *debug {
			fmt.Printf("Blueprint %d: best score %d\n", i+1, best)
		}
		total += (i + 1) * int(best)
	}
	fmt.Println(total)

	// part B
	product := 1
	for i := 0; i < 3; i++ {
		best := recipes[i].findBest(32)
		if *debug {
			fmt.Printf("Blueprint %d: best score %d\n", i+1, best)
		}
		product *= int(best)
	}
	fmt.Println(product)
}

func (bp recipe) findBest(maxTime int8) int8 {
	var best int8
	visited := make(map[int32]bool)
	q := []state{{
		bots: inventory{1, 0, 0, 0},
	}}
	q2 := []int8{0}
	for len(q) > 0 {
		head := q[0]
		turns := q2[0]
		q = q[1:]
		q2 = q2[1:]

		// We've already evaluated this position.
		if hash := head.hash(); visited[hash] {
			continue
		} else {
			visited[hash] = true
		}

		var children []state
		if turns == maxTime-3 {
			// Score it and terminate the tree. On the second and third to last rounds the only thing that can make
			// a difference is building a geode miner.
			for j := 0; j < 2; j++ {
				if builtGeode := head.step(geode, bp); builtGeode != nil {
					head = *builtGeode
				} else {
					head = *head.step(none, bp)
				}
			}
			// Doesn't matter what we build on the last round, it won't start mining in time.
			head = *head.step(none, bp)
			if head.raw[geode] > best {
				best = head.raw[geode]
			}
			continue
		} else if turns >= maxTime-4 {
			// Build a geode bot as our first (and only) choice; if we can't build a geode bot,
			// don't bother building anything other than dependencies for geodes.
			if builtGeode := head.step(geode, bp); builtGeode != nil {
				children = append(children, *builtGeode)
			} else {
				for res := obsidian; res > none; res-- {
					if bp[obsidian][res] > 0 {
						if proposed := head.step(res, bp); proposed != nil {
							children = append(children, *proposed)
						}
					}
					// Also always propose doing nothing.
					children = append(children, *head.step(none, bp))
				}
			}
		} else {
			children = head.children(bp)
		}
		q = append(q, children...)
		for i := 0; i < len(children); i++ {
			q2 = append(q2, turns+1)
		}
	}
	return best
}
