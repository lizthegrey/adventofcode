package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"sort"
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

type moveSeq struct {
	moves   []string
	elapsed int
	score   int
}

func (ms moveSeq) toMemoKey() string {
	tmp := make([]string, len(ms.moves))
	copy(tmp, ms.moves)
	sort.Strings(tmp)
	return strings.Join(tmp, ",")
}

func (e edges) moves(ms, os moveSeq, rooms board) []moveSeq {
	last := ms.moves[len(ms.moves)-1]
	ret := make([]moveSeq, 0, len(e[last])+1-len(ms.moves))

	seen := make(map[string]bool)
	for _, pos := range ms.moves {
		seen[pos] = true
	}
	for _, pos := range os.moves {
		// Don't duplicate any work our other sequence did.
		seen[pos] = true
	}
	for next := range e[last] {
		proposedTurns := ms.elapsed + e[last][next] + 1
		if seen[next] || proposedTurns > maxTurns {
			// Don't propose backtracking or taking longer than 30 turns.
			continue
		}
		sub := make([]string, len(ms.moves)+1)
		copy(sub, ms.moves)
		sub[len(ms.moves)] = next
		ret = append(ret, moveSeq{
			moves:   sub,
			elapsed: proposedTurns,
			score:   ms.score + (maxTurns-proposedTurns)*rooms[next].flow,
		})
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
	q := []moveSeq{{moves: []string{firstRoom}}}
	// TODO: this would benefit from being recursive/DFS instead to reuse the
	// same moves array prefix and not need to copy() it to avoid stomping it.
	for len(q) > 0 {
		head := q[0]
		q = q[1:]

		moves := weights.moves(head, moveSeq{}, rooms)
		// Only check leaf nodes, since score can always improve from adding moves if possible.
		if len(moves) == 0 {
			if head.score > best {
				best = head.score
			}
		}
		q = append(q, moves...)
	}
	fmt.Println(best)

	// Part B
	best = 0
	q = []moveSeq{{moves: []string{firstRoom}, elapsed: 4}}
	memo := make(map[string]int)
	for len(q) > 0 {
		head := q[0]
		q = q[1:]

		moves := weights.moves(head, moveSeq{}, rooms)
		// Only check leaf nodes, since score can always improve from adding moves if possible.
		if len(moves) == 0 {
			key := head.toMemoKey()
			eleScore, ok := memo[key]
			if !ok {
				// Do the sub-problem with many fewer nodes. This could be memoized
				// because it doesn't care what the order of nodes visited by me is,
				// only which it should consider off limits.
				eleQ := []moveSeq{{moves: []string{firstRoom}, elapsed: 4}}
				for len(eleQ) > 0 {
					eleHead := eleQ[0]
					eleQ = eleQ[1:]
					eleMoves := weights.moves(eleHead, head, rooms)
					if len(eleMoves) == 0 {
						if eleHead.score > eleScore {
							eleScore = eleHead.score
						}
					}
					eleQ = append(eleQ, eleMoves...)
				}
				memo[key] = eleScore
			}
			totalScore := head.score + eleScore
			if totalScore > best {
				best = totalScore
			}
		}
		q = append(q, moves...)
	}
	fmt.Println(best)
}
