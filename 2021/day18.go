package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/lizthegrey/adventofcode/2021/trace"
	"go.opentelemetry.io/otel"
	// "go.opentelemetry.io/otel/attribute"
)

var inputFile = flag.String("inputFile", "inputs/day18.input", "Relative file path to use as input.")

var tr = otel.Tracer("day18")

type Pair struct {
	Left, Right   *int
	LeftP, RightP *Pair
	Depth         int
	Parent        *Pair
}

func Sum(n, m *Pair) *Pair {
	ret := Pair{
		Left:   nil,
		Right:  nil,
		LeftP:  n.ExpandDepth(),
		RightP: m.ExpandDepth(),
		Parent: nil,
	}
	ret.LeftP.Parent = &ret
	ret.RightP.Parent = &ret

	ret.Reduce()
	return &ret
}

func (p Pair) Magnitude() int {
	sum := 0
	if p.Left != nil {
		sum += 3 * *p.Left
	} else {
		sum += 3 * p.LeftP.Magnitude()
	}
	if p.Right != nil {
		sum += 2 * *p.Right
	} else {
		sum += 2 * p.RightP.Magnitude()
	}
	return sum
}

func (p *Pair) ExpandDepth() *Pair {
	p.Depth++
	if p.LeftP != nil {
		p.LeftP.ExpandDepth()
	}
	if p.RightP != nil {
		p.RightP.ExpandDepth()
	}
	return p
}

func (p *Pair) Reduce() {
	for {
		actionPerformed := p.Explode()
		if actionPerformed {
			continue
		}
		actionPerformed = p.Split()
		if !actionPerformed {
			// We are finished.
			break
		}
		// Otherwise loop around again.
	}
}

func (p *Pair) StashUpLeft(value int) {
	parent := p.Parent
	if parent == nil {
		return
	}
	if parent.RightP == p {
		if parent.Left != nil {
			*parent.Left += value
		} else {
			parent.LeftP.StashDownLeft(value)
		}
	} else if parent.LeftP == p {
		parent.StashUpLeft(value)
	} else {
		fmt.Println("Couldn't find self in parent to stash.")
	}
}

func (p *Pair) StashDownLeft(value int) {
	// StashDownLeft always starts from the right side, and should always find somewhere to land.
	if p.Right != nil {
		*p.Right += value
	} else {
		p.RightP.StashDownLeft(value)
	}
}

func (p *Pair) StashUpRight(value int) {
	parent := p.Parent
	if parent == nil {
		return
	}
	if parent.LeftP == p {
		if parent.Right != nil {
			*parent.Right += value
		} else {
			parent.RightP.StashDownRight(value)
		}
	} else if parent.RightP == p {
		parent.StashUpRight(value)
	} else {
		fmt.Println("Couldn't find self in parent to stash.")
	}
}

func (p *Pair) StashDownRight(value int) {
	// StashDownRight always starts from the left side, and should always find somewhere to land.
	if p.Left != nil {
		*p.Left += value
	} else {
		p.LeftP.StashDownRight(value)
	}
}

func (p *Pair) Explode() bool {
	if p.Depth >= 4 && p.Left != nil && p.Right != nil {
		// We explode ourself.
		p.StashUpLeft(*p.Left)
		p.StashUpRight(*p.Right)

		// Finally erase ourselves from our parent.
		if p.Parent.LeftP == p {
			p.Parent.LeftP = nil
			zero := 0
			p.Parent.Left = &zero
		} else if p.Parent.RightP == p {
			p.Parent.RightP = nil
			zero := 0
			p.Parent.Right = &zero
		} else {
			fmt.Println("Couldn't find self in parent to delete.")
			return false
		}
		return true
	}
	if p.LeftP != nil {
		if changed := p.LeftP.Explode(); changed {
			return changed
		}
	}
	if p.RightP != nil {
		if changed := p.RightP.Explode(); changed {
			return changed
		}
	}
	return false
}

func (p *Pair) Split() bool {
	// Read from left to right, first acting on ourself, then on children.
	if p.Left != nil && *p.Left > 9 {
		left := *p.Left / 2
		right := *p.Left - left
		p.LeftP = &Pair{
			Left:   &left,
			Right:  &right,
			Parent: p,
			Depth:  p.Depth + 1,
		}
		p.Left = nil
		return true
	}
	if p.LeftP != nil {
		if ret := p.LeftP.Split(); ret {
			return ret
		}
	}
	if p.Right != nil && *p.Right > 9 {
		left := *p.Right / 2
		right := *p.Right - left
		p.RightP = &Pair{
			Left:   &left,
			Right:  &right,
			Parent: p,
			Depth:  p.Depth + 1,
		}
		p.Right = nil
		return true
	}
	if p.RightP != nil {
		if ret := p.RightP.Split(); ret {
			return ret
		}
	}
	return false
}

func (p Pair) Print() {
	fmt.Printf("[")
	if p.Left != nil {
		fmt.Printf("%d", *p.Left)
	} else {
		p.LeftP.Print()
	}
	fmt.Printf(",")
	if p.Right != nil {
		fmt.Printf("%d", *p.Right)
	} else {
		p.RightP.Print()
	}
	fmt.Printf("]")
}

// Parse recursively parses input strings into Pairs.
func Parse(s string, depth int, parent *Pair) *Pair {
	if s[0] != byte('[') || s[len(s)-1] != byte(']') {
		fmt.Printf("Invalid input %s\n", s)
		return nil
	}
	var ret Pair
	ret.Depth = depth
	ret.Parent = parent
	// Either we have a plain number, or we have another Pair on the left.
	var positionAfterComma int
	if s[1] == byte('[') {
		// We need to read matching parens in order to find the substring to recurse.
		// We know s[1] is '[', so start the bracket count at 1.
		bracketCount := 1
		for i := 2; bracketCount != 0; i++ {
			switch s[i] {
			case byte('['):
				bracketCount++
			case byte(']'):
				bracketCount--
			}
			positionAfterComma = i + 2
		}
		ret.LeftP = Parse(s[1:positionAfterComma-1], depth+1, &ret)
	} else {
		// This is just a plain number. They are always single digit.
		left := int(s[1] - byte('0'))
		ret.Left = &left
		if s[2] != byte(',') {
			fmt.Printf("Invalid input %s\n", s)
			return nil
		}
		// Proceed to parsing the right part.
		positionAfterComma = 3
	}
	if s[positionAfterComma] == byte('[') {
		ret.RightP = Parse(s[positionAfterComma:len(s)-1], depth+1, &ret)
	} else {
		// This is just a plain number, single digit.
		right := int(s[positionAfterComma] - byte('0'))
		ret.Right = &right
		if positionAfterComma != len(s)-2 {
			fmt.Printf("Invalid input %s\n", s)
			return nil
		}
	}
	return &ret
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

	var numbers []*Pair
	for _, s := range split {
		numbers = append(numbers, Parse(s, 0, nil))
	}

	sum := numbers[0]
	for i, n := range numbers {
		if i == 0 {
			continue
		}
		sum = Sum(sum, n)
	}
	fmt.Println(sum.Magnitude())

	var maxMagnitude int
	for i, x := range split {
		for j, y := range split {
			if i == j {
				continue
			}
			magnitude := Sum(Parse(x, 0, nil), Parse(y, 0, nil)).Magnitude()
			if magnitude > maxMagnitude {
				maxMagnitude = magnitude
			}
		}
	}
	fmt.Println(maxMagnitude)
}
