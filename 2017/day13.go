package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

var inputFile = flag.String("inputFile", "inputs/day13.input", "Relative file path to use as input.")

type Scanner struct {
	Pos        int
	Range      int
	MovingDown bool
}
type Board struct {
	PacketDepth int
	Scanners    map[int]*Scanner
}

func main() {
	flag.Parse()

	bytes, err := ioutil.ReadFile(*inputFile)
	if err != nil {
		fmt.Printf("Could not open file %s because %v.\n", *inputFile, err)
	}
	contents := string(bytes)
	lines := strings.Split(contents, "\n")
	b := Board{-2, make(map[int]*Scanner)}

	maxDepth := 0
	for _, l := range lines {
		if len(l) == 0 {
			break
		}
		parts := strings.Split(l, ": ")

		d, err := strconv.Atoi(parts[0])
		if err != nil {
			fmt.Printf("Failed to parse depth in '%s'\n", l)
			return
		}

		r, err := strconv.Atoi(parts[1])
		if err != nil {
			fmt.Printf("Failed to parse range in '%s'\n", l)
			return
		}
		b.Scanners[d] = &Scanner{0, r, true}
		if d > maxDepth {
			maxDepth = d
		}
	}
	b.Enter()
	totalSeverity := 0
	for i := 0; i <= maxDepth; i++ {
		totalSeverity += b.Tick()
	}
	fmt.Printf("Total severity: %d\n", totalSeverity)
}

func (b *Board) Enter() {
	b.PacketDepth = -1
}

func (b *Board) Tick() int {
	severity := 0
	if b.PacketDepth >= -1 {
		b.PacketDepth++
	}
	if s, found := b.Scanners[b.PacketDepth]; found {
		if s.Pos == 0 {
			severity = b.PacketDepth * s.Range
		}
	}
	for _, v := range b.Scanners {
		if v.MovingDown {
			v.Pos++
			if v.Pos == v.Range-1 {
				v.MovingDown = false
			}
		} else {
			v.Pos--
			if v.Pos == 0 {
				v.MovingDown = true
			}
		}
	}
	return severity
}
