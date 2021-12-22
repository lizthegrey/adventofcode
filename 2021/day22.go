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

var inputFile = flag.String("inputFile", "inputs/day22.input", "Relative file path to use as input.")

var tr = otel.Tracer("day22")

type Coord3 [3]int

const X = 0
const Y = 1
const Z = 2

type Prism struct {
	Lower, Upper Coord3
}
type Instruction struct {
	ToToggle        Prism
	ToggleDirection bool
}
type PrismSet []Instruction

func (p Prism) Volume() uint64 {
	volume := uint64(1)
	for d := X; d <= Z; d++ {
		volume *= uint64(p.Upper[d] - p.Lower[d] + 1)
	}
	return volume
}

func (p Prism) Overlap(o Prism) *Prism {
	// In 2d (for each dimension), we have 4 possibilities:
	// 1: [-------]   (-------)
	// 2: [----(^^]-----------)
	// 3: (---[^^^^^^^^^^^^^]-)
	// 4: [----(^^^^^^^^^^^^)-]
	// 5: (-------[^^^^)------]
	// 6: (-------)   [-------]
	var ret Prism
	for d := X; d <= Z; d++ {
		if p.Lower[d] > o.Upper[d] {
			// No overlap, case 1
			return nil
		} else if p.Upper[d] < o.Lower[d] {
			// No overlap, case 6
			return nil
		} else if p.Upper[d] <= o.Upper[d] && p.Lower[d] >= o.Lower[d] {
			// Full overlap, case 4
			ret.Upper[d] = p.Upper[d]
			ret.Lower[d] = p.Lower[d]
		} else if o.Upper[d] <= p.Upper[d] && o.Lower[d] >= p.Lower[d] {
			// Full overlap, case 3
			ret.Upper[d] = o.Upper[d]
			ret.Lower[d] = o.Lower[d]
		} else if o.Lower[d] <= p.Lower[d] && o.Upper[d] <= p.Upper[d] {
			// Partial overlap, case 2
			ret.Upper[d] = o.Upper[d]
			ret.Lower[d] = p.Lower[d]
		} else if p.Lower[d] <= o.Lower[d] && p.Upper[d] <= o.Upper[d] {
			// Partial overlap, case 5
			ret.Upper[d] = p.Upper[d]
			ret.Lower[d] = o.Lower[d]
		} else {
			// There is a bug.
			fmt.Println("Unaccounted overlap case.")
			return nil
		}
	}
	return &ret
}

func (p Prism) Subtract(overlap Prism) []Prism {
	// Takes a prism that's currently set on in its entirety.
	// Takes a chunk out of it (already pre-verified to be an exact subset).
	// Returns all the remaining chunks.
	var ret []Prism
	if overlap == p {
		// We're deleting in our entirety.
		return ret
	}
	// We need to construct the prisms after cutting out the overlap.
	// There are 6 potential pieces we need to construct:
	// The slab below (ZBottom)
	// Where the Z coordinates overlap, we have (top-down, X and Y only):
	// [-----YUpper---------]
	// [Left] Overlap [Right]
	// [-----YLower---------]
	// The slab above (ZAbove)
	if overlap.Lower[Z] != p.Lower[Z] {
		// Construct the bottom slab.
		lower := p.Lower
		upper := p.Upper
		upper[Z] = overlap.Lower[Z] - 1
		ret = append(ret, Prism{lower, upper})
	}
	// Construct the Y-lower slab, bounded at top and bottom Z by overlap Zs
	if overlap.Lower[Y] != p.Lower[Y] {
		lower := p.Lower
		upper := p.Upper
		lower[Z] = overlap.Lower[Z]
		upper[Z] = overlap.Upper[Z]
		upper[Y] = overlap.Lower[Y] - 1
		ret = append(ret, Prism{lower, upper})
	}

	// Construct the left and right slabs.
	if overlap.Lower[X] != p.Lower[X] {
		lower := overlap.Lower
		upper := overlap.Upper
		lower[X] = p.Lower[X]
		upper[X] = overlap.Lower[X] - 1
		ret = append(ret, Prism{lower, upper})
	}

	if overlap.Upper[X] != p.Upper[X] {
		lower := overlap.Lower
		upper := overlap.Upper
		lower[X] = overlap.Upper[X] + 1
		upper[X] = p.Upper[X]
		ret = append(ret, Prism{lower, upper})
	}

	// Construct the Y-upper slab, bounded at top and bottom Z by overlap Zs
	if overlap.Upper[Y] != p.Upper[Y] {
		lower := p.Lower
		upper := p.Upper
		lower[Z] = overlap.Lower[Z]
		upper[Z] = overlap.Upper[Z]
		lower[Y] = overlap.Upper[Y] + 1
		ret = append(ret, Prism{lower, upper})
	}
	if overlap.Upper[Z] != p.Upper[Z] {
		// Construct the ZAbove slab.
		lower := p.Lower
		upper := p.Upper
		lower[Z] = overlap.Upper[Z] + 1
		ret = append(ret, Prism{lower, upper})
	}
	return ret
}

func (s PrismSet) NumberOn() uint64 {
	// prismsOn contains a non-overlapping set of prisms that are contiguously on.
	prismsOn := make(map[Prism]bool)
	for _, instr := range s {
		var toAdd []Prism
		p := instr.ToToggle
		if instr.ToggleDirection {
			toAdd = append(toAdd, p)
		}
		for on := range prismsOn {
			overlap := p.Overlap(on)
			if overlap == nil {
				continue
			}
			// We have found an overlap, which means we need to explode it.
			delete(prismsOn, on)
			// We need to take the chunk out of the pre-existing prism
			// to make enough room for ourselves to either add or subtract.
			toAdd = append(toAdd, on.Subtract(*overlap)...)
		}
		for _, v := range toAdd {
			prismsOn[v] = true
		}
	}

	var volume uint64
	for p := range prismsOn {
		volume += p.Volume()
	}
	return volume
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

	// Populate an initial clear of -100k to +100k, since no input exceeds.
	var prisms PrismSet
	for i, s := range split {
		parts := strings.Split(s, " ")
		coords := strings.Split(parts[1], ",")
		var lower, upper Coord3

		for p, v := range coords {
			innerParts := strings.Split(v, "..")
			lower[p], _ = strconv.Atoi(innerParts[0][2:])
			upper[p], _ = strconv.Atoi(innerParts[1])
		}
		if i == 20 {
			fmt.Println(prisms.NumberOn())
		}
		prisms = append(prisms, Instruction{Prism{lower, upper}, parts[0] == "on"})
	}
	fmt.Println(prisms.NumberOn())
}
