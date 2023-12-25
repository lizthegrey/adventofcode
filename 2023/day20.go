package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
)

var inputFile = flag.String("inputFile", "inputs/day20.input", "Relative file path to use as input.")

type Pulse bool

const Hi Pulse = true
const Lo Pulse = false

type Signal struct {
	Src string
	Val Pulse
	Dst string
}

type Module interface {
	Receive(string, Pulse) *Pulse
	Outs() []string
}

type Conjunction struct {
	Memory  map[string]Pulse
	Outputs []string
}

func (c *Conjunction) Receive(src string, val Pulse) *Pulse {
	c.Memory[src] = val
	ret := Lo
	for _, v := range c.Memory {
		if v == Lo {
			ret = Hi
		}
	}
	return &ret
}

func (c *Conjunction) Outs() []string {
	return c.Outputs
}

type FlipFlop struct {
	State   Pulse
	Outputs []string
}

func (ff *FlipFlop) Receive(src string, val Pulse) *Pulse {
	if val == Hi {
		return nil
	}
	ff.State = !ff.State
	ret := ff.State
	return &ret
}

func (ff *FlipFlop) Outs() []string {
	return ff.Outputs
}

func main() {
	flag.Parse()
	bytes, err := ioutil.ReadFile(*inputFile)
	if err != nil {
		return
	}
	contents := string(bytes)
	split := strings.Split(contents, "\n")

	var broadcast []string
	modules := make(map[string]Module)
	for _, s := range split[:len(split)-1] {
		parts := strings.Split(s, " -> ")
		name := parts[0][1:]
		outputs := strings.Split(parts[1], ", ")
		var module Module
		switch parts[0][0] {
		case 'b':
			broadcast = outputs
			continue
		case '&':
			module = &Conjunction{Memory: make(map[string]Pulse), Outputs: outputs}
		case '%':
			module = &FlipFlop{Outputs: outputs}
		default:
			log.Fatalf("Unknown input %s", s)
		}
		modules[name] = module
	}

	for conjName, m := range modules {
		if c, ok := m.(*Conjunction); ok {
			for other, n := range modules {
				for _, out := range n.Outs() {
					if out == conjName {
						c.Memory[other] = false
					}
				}
			}
		}
	}

	var los, his, seen int
	lcm := 1
	var rmMem map[string]Pulse
	var partB string
outer:
	for k, v := range modules {
		for _, out := range v.Outs() {
			if out == "rx" {
				partB = k
				break outer
			}
		}
	}
	rm, bigInput := modules[partB].(*Conjunction)
	if bigInput {
		rmMem = rm.Memory
	}
	for round := 1; ; round++ {
		var q []Signal
		for _, dst := range broadcast {
			q = append(q, Signal{Src: "broadcast", Dst: dst, Val: Lo})
		}
		// Account for the button and the broadcast fan-out.
		los += len(broadcast) + 1
		for len(q) > 0 {
			var nextQ []Signal
			for _, v := range q {
				cur := v.Dst
				m, ok := modules[cur]
				if !ok {
					continue
				}
				val := m.Receive(v.Src, v.Val)
				if val == nil {
					continue
				}
				if *val == Lo {
					los += len(m.Outs())
				} else {
					his += len(m.Outs())
				}
				for _, dst := range m.Outs() {
					if _, ok := rmMem[dst]; ok && *val == Lo {
						seen++
						lcm *= round
						if seen == len(rmMem) {
							fmt.Println(lcm)
							return
						}
					}
					send := Signal{Src: cur, Dst: dst, Val: *val}
					nextQ = append(nextQ, send)
				}
			}
			q = nextQ
		}
		if round == 1000 {
			fmt.Println(uint64(los) * uint64(his))
			if !bigInput {
				return
			}
		}
	}
}
