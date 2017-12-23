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
	MaxDepth    int
	Scanners    map[int]Scanner
}

func main() {
	flag.Parse()

	bytes, err := ioutil.ReadFile(*inputFile)
	if err != nil {
		fmt.Printf("Could not open file %s because %v.\n", *inputFile, err)
	}
	contents := string(bytes)

	b := SetupBoard(contents)

	b = b.Enter()
	totalSeverity := 0
	for b.PacketDepth <= b.MaxDepth {
		severity := 0
		b, _, severity = b.Tick()
		totalSeverity += severity
	}
	fmt.Printf("Total severity on naive entry: %d\n", totalSeverity)

	// Ensure we don't have to replay the state up to the tick before we entered every time.
	memo := SetupBoard(contents)

	outer:
	for toDelay := 1;; toDelay++ {
		memo, _, _ = memo.Tick()

		b := memo.Enter()
		for b.PacketDepth <= b.MaxDepth {
			caught := false
			b, caught, _ = b.Tick()
			if caught {
				// We were caught; this delay won't work.
				continue outer
			}
		}
		fmt.Printf("Got through with delay %d.\n", toDelay)
		break
	}
}

func SetupBoard(contents string) Board {
	lines := strings.Split(contents, "\n")
	b := Board{-2, 0, make(map[int]Scanner)}

	for _, l := range lines {
		if len(l) == 0 {
			break
		}
		parts := strings.Split(l, ": ")

		d, err := strconv.Atoi(parts[0])
		if err != nil {
			fmt.Printf("Failed to parse depth in '%s'\n", l)
			return b
		}

		r, err := strconv.Atoi(parts[1])
		if err != nil {
			fmt.Printf("Failed to parse range in '%s'\n", l)
			return b
		}
		b.Scanners[d] = Scanner{0, r, true}
		if d > b.MaxDepth {
			b.MaxDepth = d
		}
	}
	return b
}

func (b Board) Enter() Board {
	ret := b
	ret.Scanners = make(map[int]Scanner)
	for k,v := range b.Scanners {
		ret.Scanners[k] = v
	}
	ret.PacketDepth = -1
	return ret
}

func (b Board) Tick() (Board, bool, int) {
	ret := b

	caught := false
	severity := 0
	if ret.PacketDepth >= -1 {
		ret.PacketDepth++
	}
	if s, found := ret.Scanners[ret.PacketDepth]; found {
		if s.Pos == 0 {
			caught = true
			severity = ret.PacketDepth * s.Range
		}
	}
	for k, v := range ret.Scanners {
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
		ret.Scanners[k] = v
	}
	return ret, caught, severity
}
