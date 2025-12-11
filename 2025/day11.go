package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strings"
)

var inputFile = flag.String("inputFile", "inputs/day11.input", "Relative file path to use as input.")

// for some reason coincidentally there are 576 different devices, which is 64*9.
// this cannot possibly be a coincidence.
type CacheKey [9]uint64

func (c CacheKey) Used(d Device) bool {
	for i := range c {
		if d.Key[i]&c[i] > 0 {
			return true
		}
	}
	return false
}

func (c CacheKey) MarkUsed(d Device) CacheKey {
	for i := range c {
		c[i] |= d.Key[i]
	}
	return c
}

type Device struct {
	Name  string
	Conns []string
	Key   CacheKey
}

type State struct {
	Visited CacheKey
	Cur     string
}

func main() {
	flag.Parse()
	bytes, err := ioutil.ReadFile(*inputFile)
	if err != nil {
		return
	}
	contents := string(bytes)
	split := strings.Split(contents, "\n")

	devices := make(map[string]Device)
	for i, s := range split[:len(split)-1] {
		parts := strings.Split(s, ": ")

		var d Device
		d.Name = parts[0]
		for _, v := range strings.Split(parts[1], " ") {
			d.Conns = append(d.Conns, v)
		}
		d.Key[i/64] = 1 << (i % 64)
		devices[d.Name] = d
	}
	memo := make(map[State]uint64)
	fmt.Println(compute(devices, memo, State{devices["you"].Key, "you"}, "out"))
	clear(memo)
	fmt.Println(compute(devices, memo, State{devices["svr"].Key, "svr"}, "out"))
}

func compute(devices map[string]Device, memo map[State]uint64, state State, fin string) uint64 {
	var ret uint64
	if state.Cur == fin {
		// Mark invalid entries bad
		if state.Visited.Used(devices["svr"]) && !(state.Visited.Used(devices["fft"]) && state.Visited.Used(devices["dac"])) {
			return 0
		}
		return 1
	}
	if ret, ok := memo[state]; ok {
		return ret
	}
	for _, v := range devices[state.Cur].Conns {
		if state.Visited.Used(devices[v]) {
			continue
		}
		next := state
		next.Cur = v
		next.Visited = next.Visited.MarkUsed(devices[v])
		ret += compute(devices, memo, next, fin)
	}

	memo[state] = ret
	return ret
}
