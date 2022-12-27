package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

var inputFile = flag.String("inputFile", "inputs/day16.input", "Relative file path to use as input.")

const maxTurns = 30
const firstRoom = "AA"

type room struct {
	name  string
	flow  int
	neigh []string
}

type board map[string]room

type edges map[string]map[string]int

type moveSeq []string

func (e edges) score(b board, ms moveSeq) int {
	var rate, sum, elapsed int
	prev := ms[0]
	for i := 1; i < len(ms); i++ {
		cur := ms[i]
		// Include the time to turn the lever.
		turns := e[prev][cur] + 1
		sum += rate * turns
		elapsed += turns
		rate += b[cur].flow
		prev = cur
	}
	// After all moves, keep the water gushing for the total 30 turns.
	sum += rate * (maxTurns - elapsed)
	return sum
}

func (e edges) moves(ms, os moveSeq) []moveSeq {
	last := ms[len(ms)-1]
	ret := make([]moveSeq, 0, len(e[last])+1-len(ms))

	var elapsed int
	seen := make(map[string]bool)
	prev := ms[0]
	for i := 1; i < len(ms); i++ {
		cur := ms[i]
		seen[cur] = true
		turns := e[prev][cur] + 1
		elapsed += turns
		prev = cur
	}
	for _, v := range os {
		// Don't duplicate any work our other sequence did.
		seen[v] = true
	}
	for next := range e[last] {
		if seen[next] || elapsed+e[last][next]+1 > maxTurns {
			// Don't propose backtracking or taking longer than 30 turns.
			continue
		}
		sub := make(moveSeq, len(ms)+1)
		copy(sub, ms)
		sub[len(ms)] = next
		ret = append(ret, sub)
	}
	return ret
}

func main() {
	flag.Parse()
	bytes, err := ioutil.ReadFile(*inputFile)
	if err != nil {
		return
	}
	contents := string(bytes)
	split := strings.Split(contents, "\n")

	rooms := make(board)
	weights := make(edges)
	for _, s := range split[:len(split)-1] {
		parts := strings.SplitN(s, " ", 10)
		name := parts[1]
		flow, _ := strconv.Atoi(parts[4][5 : len(parts[4])-1])
		neighbours := strings.Split(parts[9], ", ")
		rooms[name] = room{
			name:  name,
			flow:  flow,
			neigh: neighbours,
		}
		if name == firstRoom || flow > 0 {
			weights[name] = make(map[string]int)
		}
	}

	// Pre-process by generating shortest distances using BFS from "AA" plus
	// all operative valves, to all operative valves. All weights are equal
	// so we don't need F-W or Dijkstra's.
	for start := range weights {
		q := []string{start}
		var iters int
		for len(q) > 0 {
			var next []string
			for _, v := range q {
				if weights[start][v] != 0 {
					continue
				}
				weights[start][v] = iters
				next = append(next, rooms[v].neigh...)
			}
			q = next
			iters++
		}

		// Finally, delete all unneeded records. once we've finished our traversal.
		for dst := range weights[start] {
			if weights[dst] == nil || dst == firstRoom {
				delete(weights[start], dst)
			}
		}
	}

	// Part A
	var best int
	q := []moveSeq{{firstRoom}}
	for len(q) > 0 {
		head := q[0]
		q = q[1:]

		moves := weights.moves(head, nil)
		// Only check leaf nodes, since score can always improve from adding moves if possible.
		if len(moves) == 0 {
			score := weights.score(rooms, head)
			if score > best {
				best = score
			}
		}
		q = append(q, moves...)
	}

	fmt.Println(best)

	// Part B
	// Now we have two actors, not just one. We'll add 4 to every weight starting from AA
	// to avoid needing to change any of the other code while accounting for teaching time.
	for k := range weights[firstRoom] {
		weights[firstRoom][k] += 4
	}

	best = 0
	q = []moveSeq{{firstRoom}}
	for len(q) > 0 {
		head := q[0]
		q = q[1:]

		moves := weights.moves(head, nil)
		// Only check leaf nodes, since score can always improve from adding moves if possible.
		if len(moves) == 0 {
			score := weights.score(rooms, head)

			// Do the sub-problem with many fewer nodes. This could be memoized
			// because it doesn't care what the order of nodes visited by me is,
			// only which it should consider off limits.
			eleQ := []moveSeq{{firstRoom}}
			for len(eleQ) > 0 {
				eleHead := eleQ[0]
				eleQ = eleQ[1:]

				eleMoves := weights.moves(eleHead, head)
				if len(eleMoves) == 0 {
					// Both elephant and I have visited as many nodes as possible.
					// They won't overlap the rooms visited, so just straight add.
					totalScore := score + weights.score(rooms, eleHead)
					if totalScore > best {
						best = totalScore
					}
				}
				eleQ = append(eleQ, eleMoves...)
			}
		}
		q = append(q, moves...)
	}

	fmt.Println(best)
}
