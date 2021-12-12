package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strings"
)

var inputFile = flag.String("inputFile", "inputs/day12.input", "Relative file path to use as input.")

type Exits map[string][]string
type Path []string
type PathWithRepeat struct {
	Visited      Path
	RepeatedNode string
}

// Extend copies the current Path to not share the slice, and appends next node.
func (p Path) Extend(next string) Path {
	pathCopy := make(Path, len(p), len(p)+1)
	copy(pathCopy, p)
	pathCopy = append(pathCopy, next)
	return pathCopy
}

// mayVisit returns whether we may visit a node (and false if it has
// already been visited too many times). It also returns whether the
// visit is our one allowed repeat visit of a little node.
func (p PathWithRepeat) mayVisit(n string) (bool, bool) {
	// Big caves may always be visited.
	if n[0] >= 'A' && n[0] <= 'Z' {
		return true, false
	}
	// If we've already used our repeated node and we're there again, no 3rd visit.
	// Also special-cases start, which we're never allowed to re-visit.
	if n == p.RepeatedNode || n == "start" {
		return false, false
	}
	// Check all of the visited nodes in our path for duplication.
	for _, e := range p.Visited {
		if e == n {
			if p.RepeatedNode == "" {
				// This is okay to visit, but must be flagged as our repeat so
				// we don't visit a third time, or visit another node a second time.
				return true, true
			} else {
				// This is not okay to visit because we've already used our repeat
				// on another small cave in this path.
				return false, false
			}
		}
	}
	return true, false
}

func main() {
	flag.Parse()
	bytes, err := ioutil.ReadFile(*inputFile)
	if err != nil {
		return
	}
	contents := string(bytes)
	split := strings.Split(contents, "\n")
	split = split[:len(split)-1]

	exits := make(Exits)
	for _, s := range split {
		parts := strings.Split(s, "-")
		exits[parts[0]] = append(exits[parts[0]], parts[1])
		exits[parts[1]] = append(exits[parts[1]], parts[0])
	}
	fmt.Println(exits.Search([]PathWithRepeat{
		{Visited: Path{"start"}, RepeatedNode: "invalid"},
	}))
	fmt.Println(exits.Search([]PathWithRepeat{
		{Visited: Path{"start"}, RepeatedNode: ""},
	}))
}

// Search performs a breadth-first search non-recursively using a worklist to
// identify the number of unique paths that can be traversed.
// Paths end and are counted if they reach "end", and paths are only extended
// and put onto worklist for their neighboring reachable nodes.
func (exits Exits) Search(queue []PathWithRepeat) int {
	paths := 0
	for len(queue) != 0 {
		var nextQueue []PathWithRepeat
		for _, item := range queue {
			path := item.Visited
			last := path[len(path)-1]
			if last == "end" {
				// This is a unique path that has reached the end.
				paths++
				continue
			}
			for _, n := range exits[last] {
				mayVisit, isRepeat := item.mayVisit(n)
				if !mayVisit {
					continue
				}
				nextItem := PathWithRepeat{
					Visited:      path.Extend(n),
					RepeatedNode: item.RepeatedNode,
				}
				if isRepeat {
					nextItem.RepeatedNode = n
				}
				nextQueue = append(nextQueue, nextItem)
			}
		}
		queue = nextQueue
	}
	return paths
}
