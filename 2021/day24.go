package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/lizthegrey/adventofcode/2021/trace"
	"go.opentelemetry.io/otel"
	//"go.opentelemetry.io/otel/attribute"
)

var inputFile = flag.String("inputFile", "inputs/day24.input", "Relative file path to use as input.")
var candidate = flag.Uint64("candidate", 99999999999999, "Seed candidate to use.")
var debug = flag.Bool("debug", false, "Whether to print debug output.")

var tr = otel.Tracer("day24")

type Instruction func(*Computer)

type Computer struct {
	Inputs [14]int
	Instrs []Instruction
	Regs   [4]uint64
	PC     int
	InputC int
}

type Register int

const (
	Undef Register = iota - 1
	W
	X
	Y
	Z
)

func (c *Computer) Step() bool {
	c.Instrs[c.PC](c)
	c.PC++
	if c.PC >= len(c.Instrs) {
		return false
	}
	return true
}

func main() {
	flag.Parse()

	ctx := context.Background()
	hny, tp := trace.InitializeTracing(ctx)
	defer hny.Shutdown(ctx)
	defer tp.Shutdown(ctx)

	bytes, err := ioutil.ReadFile(*inputFile)
	if err != nil {
		return
	}
	contents := string(bytes)
	split := strings.Split(contents, "\n")
	split = split[:len(split)-1]

	var computer Computer
	for _, s := range split {
		parts := strings.Split(s, " ")
		var instr Instruction
		target := Register(parts[1][0] - byte('w'))
		var immediate int
		src := Register(Undef)
		if len(parts) > 2 {
			immediate, err = strconv.Atoi(parts[2])
			if err != nil {
				src = Register(parts[2][0] - byte('w'))
			}
		}
		switch parts[0] {
		case "inp":
			instr = func(c *Computer) {
				c.Regs[target] = uint64(c.Inputs[c.InputC])
				c.InputC++
			}
		case "add":
			instr = func(c *Computer) {
				if src != Undef {
					c.Regs[target] += c.Regs[src]
				} else {
					c.Regs[target] += uint64(immediate)
				}
			}
		case "mul":
			instr = func(c *Computer) {
				if src != Undef {
					c.Regs[target] *= c.Regs[src]
				} else {
					c.Regs[target] *= uint64(immediate)
				}
			}
		case "div":
			instr = func(c *Computer) {
				if src != Undef {
					c.Regs[target] /= c.Regs[src]
				} else {
					c.Regs[target] /= uint64(immediate)
				}
			}
		case "mod":
			instr = func(c *Computer) {
				if src != Undef {
					c.Regs[target] %= c.Regs[src]
				} else {
					c.Regs[target] %= uint64(immediate)
				}
			}
		case "eql":
			instr = func(c *Computer) {
				var equal bool
				if src != Undef {
					equal = c.Regs[target] == c.Regs[src]
				} else {
					equal = c.Regs[target] == uint64(immediate)
				}
				if equal {
					c.Regs[target] = 1
				} else {
					c.Regs[target] = 0
				}
			}
		}
		computer.Instrs = append(computer.Instrs, instr)
	}

	v := *candidate
	for i := 13; i >= 0; i-- {
		computer.Inputs[i] = int(v % 10)
		v /= 10
	}
	for computer.Step() {
		if *debug {
			fmt.Println(computer.Regs)
		}
	}
	success := computer.Regs[Z] == 0
	fmt.Println(success)
}
