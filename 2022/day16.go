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

func (e edges) findBest(ms, os moveSeq, rooms board, scoring func(moveSeq) int) int {
	// Allocate enough space in the array.
	if cap(ms.moves) < len(e)-1 {
		buf := make([]string, len(ms.moves), len(e)-1)
		copy(buf, ms.moves)
		ms.moves = buf
	}

	last := ms.moves[len(ms.moves)-1]
	var best int

	seen := make(map[string]bool)
	for _, pos := range ms.moves {
		seen[pos] = true
	}
	for _, pos := range os.moves {
		// Don't duplicate any work our other sequence did.
		seen[pos] = true
	}
	leaf := true
	oldLen := len(ms.moves)
	for next := range e[last] {
		proposedTurns := ms.elapsed + e[last][next] + 1
		if seen[next] || proposedTurns > maxTurns {
			// Don't propose backtracking or taking longer than 30 turns.
			continue
		}
		leaf = false
		// Reslice.
		ms.moves = ms.moves[0 : oldLen+1]
		ms.moves[oldLen] = next
		score := e.findBest(moveSeq{
			moves:   ms.moves,
			elapsed: proposedTurns,
			score:   ms.score + (maxTurns-proposedTurns)*rooms[next].flow,
		}, os, rooms, scoring)
		if score > best {
			best = score
		}
	}
	// Restore previous slicing.
	ms.moves = ms.moves[0:oldLen]
	if leaf {
		if scoring != nil {
			return scoring(ms)
		}
		return ms.score
	}
	return best
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
	initial := moveSeq{moves: []string{firstRoom}}
	fmt.Println(weights.findBest(initial, moveSeq{}, rooms, nil))

	// Part B
	initial = moveSeq{moves: []string{firstRoom}, elapsed: 4}
	memo := make(map[string]int)
	fmt.Println(weights.findBest(initial, moveSeq{}, rooms, func(ms moveSeq) int {
		key := ms.toMemoKey()
		eleScore, ok := memo[key]
		if !ok {
			// Do the sub-problem with many fewer nodes. This could be memoized
			// because it doesn't care what the order of nodes visited by me is,
			// only which it should consider off limits.
			ele := moveSeq{moves: []string{firstRoom}, elapsed: 4}
			eleScore = weights.findBest(ele, moveSeq{}, rooms, nil)
			memo[key] = eleScore
		}
		return ms.score + eleScore
	}))
}
