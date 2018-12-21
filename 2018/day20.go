package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"
)

var inputFile = flag.String("inputFile", "inputs/day20.input", "Relative file path to use as input.")
var verbose = flag.Bool("verbose", false, "Whether to print verbose debug output.")
var partB = flag.Bool("partB", false, "Whether to use the Part B logic.")

type Coord struct {
	X, Y int
}

type BFSEntry struct {
	c Coord
	d int
}

// Coord -> NESW; missing adjacency = impassable.
type Adjacencies map[Coord]string

func (a Adjacencies) Add(l Coord, d rune) {
	found := false
	for _, c := range a[l] {
		if c == d {
			// Wall already torn down.
			found = true
			break
		}
	}
	if !found {
		a[l] += string(d)
	}
}

func (a Adjacencies) Neighbors(c Coord) []Coord {
	ret := make([]Coord, 0)
	for _, d := range a[c] {
		l := c
		switch d {
		case 'N':
			l.Y--
		case 'E':
			l.X++
		case 'S':
			l.Y++
		case 'W':
			l.X--
		}
		ret = append(ret, l)
	}
	return ret
}

type PathSubstring struct {
	Payload     string
	Branches    []*PathSubstring // OR group underneath us.
	NextSibling *PathSubstring   // AND next to us e.g. linked list.
	Parent      *PathSubstring   // Pointer to node above. null if we're the root.
}

func (p *PathSubstring) GeneratePaths() []string {
	var branches []string

	// Deal with any sub-branches
	if len(p.Branches) != 0 {
		for _, b := range p.Branches {
			branches = append(branches, b.GeneratePaths()...)
		}
	} else {
		branches = append(branches, "")
	}

	// Deal with any siblings.
	suffixes := []string{""}
	if p.NextSibling != nil {
		suffixes = p.NextSibling.GeneratePaths()
	}

	results := make([]string, 0)
	for _, br := range branches {
		for _, post := range suffixes {
			results = append(results, p.Payload+br+post)
		}
	}
	return results
}

func main() {
	flag.Parse()
	f, err := os.Open(*inputFile)
	if err != nil {
		return
	}
	defer f.Close()

	var r *regexp.Regexp
	var raw string

	reader := bufio.NewReader(f)
	for {
		l, err := reader.ReadString('\n')
		if err != nil || len(l) == 0 {
			break
		}
		l = l[:len(l)-1]
		r = regexp.MustCompile(l)
		// Trim ^ and $
		raw = l[1 : len(l)-1]
	}

	// Build up a list of paths by structuring the regex as a tree.
	current := &PathSubstring{}
	root := current
	for _, c := range raw {
		if *verbose {
			fmt.Printf("%c, %v\n", c, current)
		}

		switch c {
		case 'N':
			fallthrough
		case 'E':
			fallthrough
		case 'S':
			fallthrough
		case 'W':
			current.Payload += string(c)
		case '(':
			child := &PathSubstring{}
			child.Parent = current
			current.Branches = []*PathSubstring{child}
			current = child
		case ')':
			sibling := &PathSubstring{}
			sibling.Parent = current.Parent.Parent
			current.Parent.NextSibling = sibling
			current = sibling
		case '|':
			sibling := &PathSubstring{}
			sibling.Parent = current.Parent
			current.Parent.Branches = append(current.Parent.Branches, sibling)
			current = sibling
		}
	}

	paths := root.GeneratePaths()

	for _, p := range paths {
		// Verify that each path matches the regexp as a final test.
		if !r.MatchString(p) {
			fmt.Printf("Failed to match generated route %s.\n", p)
			return
		}
		if *verbose {
			fmt.Println(p)
		}
	}

	adj := make(Adjacencies)
	for _, p := range paths {
		l := Coord{0, 0}
		// Walk through the path, tearing down walls and putting in doors as we go.
		for _, d := range p {
			adj.Add(l, d)
			switch d {
			case 'N':
				l.Y--
				adj.Add(l, 'S')
			case 'E':
				l.X++
				adj.Add(l, 'W')
			case 'S':
				l.Y++
				adj.Add(l, 'N')
			case 'W':
				l.X--
				adj.Add(l, 'E')
			}
		}
	}

	// Start our breadth first search.
	seenShortestDistances := make(map[Coord]int)
	longest := 0
	q := []BFSEntry{{Coord{0, 0}, 0}}

	for len(q) > 0 {
		l := q[0]
		q = q[1:]

		if _, ok := seenShortestDistances[l.c]; ok {
			continue
		}

		neighbors := adj.Neighbors(l.c)
		qadditions := make([]BFSEntry, len(neighbors))
		seenShortestDistances[l.c] = l.d
		if l.d > longest {
			longest = l.d
		}

		for i, n := range neighbors {
			qadditions[i] = BFSEntry{n, l.d + 1}
		}
		q = append(q, qadditions...)
	}
	fmt.Println(longest)
}
