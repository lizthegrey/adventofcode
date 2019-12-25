package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/lizthegrey/adventofcode/2019/intcode"
	"os"
	"strings"
	"sync"
)

var inputFile = flag.String("inputFile", "inputs/day25.input", "Relative file path to use as input.")

var disallow = []string{
	"photons",
	"infinite loop",
}

type Coord struct {
	X, Y int
}

func (c Coord) Move(dir string) Coord {
	ret := c
	switch dir {
	case "north":
		ret.Y++
	case "east":
		ret.X++
	case "south":
		ret.Y--
	case "west":
		ret.X--
	}
	return ret
}

func main() {
	flag.Parse()
	tape := intcode.ReadInput(*inputFile)
	if tape == nil {
		fmt.Println("Failed to parse input.")
		return
	}

	workingTape := tape.Copy()
	input := make(chan int, 1)
	output, done := workingTape.Process(input)
	lines := make(chan string, 50)
	driver := make(chan string)

	go func() {
		line := make([]byte, 0)
		for c := range output {
			if c > 255 {
				fmt.Println(c)
				return
			}
			fmt.Printf("%c", c)

			if c == '\n' {
				lines <- string(line)
				line = make([]byte, 0)
			} else {
				line = append(line, byte(c))
			}
		}
		fmt.Println()
	}()

	go func() {
		reader := bufio.NewReader(os.Stdin)
		for {
			line, err := reader.ReadString('\n')
			if line == "\n" || err != nil {
				return
			}
			for _, r := range line {
				input <- int(r)
			}
		}
	}()

	go func() {
		for l := range driver {
			fmt.Printf("[Automated]: %s\n", l)
			for _, r := range l {
				input <- int(r)
			}
			input <- int('\n')
		}
	}()

	// Always pick up items unless on disallow
	// Stop and wait for human if we encounter unexpected problems.
	go func() {
		loc := Coord{0, 0}
		rooms := make(map[Coord]string)
		var items, outbound []string
		var q Queue
		parseMode := 0
		for l := range lines {
			if l == "" {
				parseMode = 0
				continue
			}
			if l == "You can't go that way." {
				fmt.Println("Disabling automatic system.")
				return
			}
			if l[0] == '=' {
				rooms[loc] = l[3 : len(l)-3]
				items = nil
				outbound = nil
				parseMode = 0
			} else if l[0] == '-' {
				// This is a list.
				// Doors here lead:
				// - north
				//
				// Items here:
				// - cake
				if parseMode == 1 {
					outbound = append(outbound, l[2:])
				} else if parseMode == 2 {
					items = append(items, l[2:])
				} else {
					fmt.Println("Didn't know what to do with a list outside parse mode.")
				}
			} else if l[len(l)-1] == ':' {
				// This is a menu ("Doors here lead:" or "Items here:").
				keyword := strings.Split(l, " ")[0]
				switch keyword {
				case "Doors":
					parseMode = 1
				case "Items":
					parseMode = 2
				default:
					fmt.Println("Unknown menu type.")
				}
			} else if l[len(l)-1] == '?' {
				// Program finally is waiting on our input. Decide what to do next.
				if len(items) > 0 {
					item := items[0]
					skip := false
					for _, v := range disallow {
						if v == item {
							skip = true
							break
						}
					}
					if !skip {
						driver <- fmt.Sprintf("take %s", items[0])
					}
					if len(items) > 1 {
						items = items[1:]
					} else {
						items = nil
					}
				} else if len(outbound) > 0 {
					for _, d := range outbound {
						// TODO rip out and replace with a normal edge to node name system.
						var dst Coord
						switch d {
						case "north":
							dst = Coord{loc.X, loc.Y + 1}
						case "east":
							dst = Coord{loc.X + 1, loc.Y}
						case "south":
							dst = Coord{loc.X, loc.Y - 1}
						case "west":
							dst = Coord{loc.X - 1, loc.Y}
						}
						if _, ok := rooms[dst]; !ok {
							// Throw it onto our exploration queue.
							fmt.Printf("Adding new room %v\n", dst)
							notifyRoom(&q, dst)
						}
					}
					// Pick a direction to move.
					next := nextDirection(&q, loc, rooms)
					if next == "dance" {
						done <- true
					}
					loc = loc.Move(next)
					driver <- next
				} else {
					fmt.Println("Got into somewhere we can't leave")
				}
			} else {
				// This is presumed to be flavor text or a response to command.
				// Just let the program print.
			}
		}
	}()
	<-done
}

type Queue struct {
	mtx  sync.Mutex
	list []Coord
	dst  Coord
}

func notifyRoom(q *Queue, nr Coord) {
	q.mtx.Lock()
	q.list = append(q.list, nr)
	q.mtx.Unlock()
}

func nextDirection(q *Queue, loc Coord, rooms map[Coord]string) string {
	// Keep track of what's unexplored, and backtrack to try to find unexplored nodes.
	// Use a DFS in order to avoid repeated backtracking.
	q.mtx.Lock()
	defer q.mtx.Unlock()

	if loc != q.dst {
		// We are in the middle of pathing somewhere...
		// TODO: read from the path cache.
		return bfs(loc, q.dst, rooms)[0]
	}
	// We've reached our destination and need a new destination.
	if len(q.list) == 0 {
		fmt.Println("Nothing left to explore. Result:")
		return "dance"
	}

	var toExplore Coord
	for {
		toExplore = q.list[0]
		if _, ok := rooms[toExplore]; ok {
			// Short circuit.
			q.list = q.list[1:]
			continue
		}
		break
	}
	if len(q.list) > 1 {
		q.list = q.list[1:]
	} else {
		q.list = nil
	}
	q.dst = toExplore
	// TODO: write to the path cache.
	return bfs(loc, q.dst, rooms)[0]
}

func bfs(src, dst Coord, rooms map[Coord]string) []string {
	// Perform a breadth-first search.
	shortest := map[Coord][]string{
		src: []string{},
	}
	worklist := []Coord{src}
	for {
		w := worklist[0]
		for _, dir := range []string{"north", "east", "south", "west"} {
			moved := w.Move(dir)
			if _, ok := rooms[moved]; !ok && moved != dst {
				// Allow ourselves to pass on unknown ground only for the last move.
				continue
			}
			if _, ok := shortest[moved]; ok {
				continue
			} else {
				directions := make([]string, len(shortest[w])+1)
				copy(directions, shortest[w])
				directions[len(shortest[w])] = dir
				shortest[moved] = directions
				if moved == dst {
					return shortest[moved]
				}
				worklist = append(worklist, moved)
			}
		}
		worklist = worklist[1:]
	}
}
