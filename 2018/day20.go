package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
)

var inputFile = flag.String("inputFile", "inputs/day20.input", "Relative file path to use as input.")
var minDoors = flag.Int("minDoors", 1000, "The minimum number of doors to count a room.")

type Coord struct {
	X, Y int
}

// Track for each item in the queue, how far the room is, so that we can
// track the distance for its children.
type BFSEntry struct {
	c Coord
	d int
}

// Coord -> NESW; missing adjacency = impassable.
// Less memory-efficient than an integer, but easier cognitively.
type Adjacencies map[Coord]string

// Function Add creates a door between a coordinate and the room in direction d from it.
func (a Adjacencies) Add(l Coord, d rune) {
	found := false
	for _, c := range a[l] {
		if c == d {
			// Wall already torn down.
			found = true
			break
		}
	}
	// Don't add the same room link twice.
	if !found {
		a[l] += string(d)
	}
}

// Function Neighbors returns a slice of coordinates of valid neighboring rooms.
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

// Function Move creates a hole in the wall in both directions and returns the new coordinate.
func (l Coord) Move(adj Adjacencies, d rune) Coord {
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
	return l
}

// Parse tree object tracking sibling and branch children for evaluation.
type PathSubstring struct {
	Payload     string
	Branches    []*PathSubstring // OR group underneath us.
	NextSibling *PathSubstring   // AND next to us e.g. linked list.
	Parent      *PathSubstring   // Pointer to node above. null if we're the root.
}

// Takes in a starting coordinatee, and punches holes starting from there that match the pattern.
func (p *PathSubstring) GeneratePaths(adj Adjacencies, l Coord) map[Coord]bool {
	// Do the unconditional moves from our payload.
	for _, c := range p.Payload {
		// Modifies l as we move along.
		l = l.Move(adj, c)
	}
	// Now we're at the end of our initial movement.

	// Deal with any sub-branches, starting from the position after initial movement.
	postBranchCoords := make(map[Coord]bool)
	if len(p.Branches) > 0 {
		for _, b := range p.Branches {
			for k := range b.GeneratePaths(adj, l) {
				postBranchCoords[k] = true
			}
		}
	} else {
		// If no sub-branches, just pass the initial movement coordinate.
		postBranchCoords[l] = true
	}

	// If no siblings, return where we wound up after initial and branches.
	if p.NextSibling == nil {
		return postBranchCoords
	}

	// Starting from each branch point, also check the result after siblings.
	finalCoords := make(map[Coord]bool)
	for b := range postBranchCoords {
		for k := range p.NextSibling.GeneratePaths(adj, b) {
			finalCoords[k] = true
		}
	}
	return finalCoords
}

func main() {
	flag.Parse()
	f, err := os.Open(*inputFile)
	if err != nil {
		return
	}
	defer f.Close()

	var raw string

	reader := bufio.NewReader(f)
	for {
		l, err := reader.ReadString('\n')
		if err != nil || len(l) == 0 {
			break
		}
		l = l[:len(l)-1]
		raw = l[1 : len(l)-1]
	}

	// Build up a list of paths by structuring the regex as a tree.
	current := &PathSubstring{}
	root := current
	for _, c := range raw {
		switch c {
		case 'N':
			fallthrough
		case 'E':
			fallthrough
		case 'S':
			fallthrough
		case 'W':
			// We're adding literals to a block.
			current.Payload += string(c)
		case '(':
			// Start a new child and put it under the current node.
			child := &PathSubstring{}
			child.Parent = current
			current.Branches = []*PathSubstring{child}
			current = child
		case ')':
			// Close off the current block, and then create a sibling to our parent.
			sibling := &PathSubstring{}
			sibling.Parent = current.Parent.Parent
			current.Parent.NextSibling = sibling
			current = sibling
		case '|':
			// Create a branch alternative.
			alternative := &PathSubstring{}
			alternative.Parent = current.Parent
			current.Parent.Branches = append(current.Parent.Branches, alternative)
			current = alternative
		}
	}

	// Track the holes we've madde in walls.
	adj := make(Adjacencies)

	// Start from 0,0 and pass the head of the parse tree.
	root.GeneratePaths(adj, Coord{0, 0})

	// Start our breadth first search for distances.
	seenShortestDistances := make(map[Coord]int)
	longest := 0
	q := []BFSEntry{{Coord{0, 0}, 0}}

	// Part B asks us how many long paths there are.
	longPaths := 0

	for len(q) > 0 {
		l := q[0]
		q = q[1:]

		if _, ok := seenShortestDistances[l.c]; ok {
			// Termination condition: don't add to the queue if we've seen a node.
			continue
		}

		neighbors := adj.Neighbors(l.c)
		qadditions := make([]BFSEntry, len(neighbors))
		seenShortestDistances[l.c] = l.d

		if l.d >= *minDoors {
			longPaths++
		}
		if l.d > longest {
			longest = l.d
		}

		for i, n := range neighbors {
			// Remember to increment the distance for our neighbors.
			qadditions[i] = BFSEntry{n, l.d + 1}
		}
		q = append(q, qadditions...)
	}

	fmt.Printf("Furthest room %d away; %d rooms %d doors away or more.\n", longest, longPaths, *minDoors)
}
