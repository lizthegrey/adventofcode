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

var inputFile = flag.String("inputFile", "inputs/day21.input", "Relative file path to use as input.")

var tr = otel.Tracer("day21")

type Player struct {
	Position int
	Score    int
}

type State struct {
	Players    [2]Player
	NextPlayer int
}

type OutcomeCounts [2]uint64

//                  Some other branching path [7,6]
//                  / weight 1
// State A -> State B  [1*3+7*1,6*1] -3-> Player 0 wins [1,0]
//            /
//         State C
type PathMap map[State]map[State]int
type Outcomes map[State]OutcomeCounts

type DeterministicDie struct {
	Value int
}

func (d *DeterministicDie) Roll() int {
	ret := d.Value
	d.Value++
	if d.Value == 101 {
		d.Value = 1
	}
	return ret
}

func (p Player) Advance(roll int) Player {
	p.Position += roll - 1
	p.Position = (p.Position % 10) + 1
	p.Score += p.Position
	return p
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

	p1, _ := strconv.Atoi(strings.Split(split[0], " ")[4])
	p2, _ := strconv.Atoi(strings.Split(split[1], " ")[4])
	player1 := Player{
		Score:    0,
		Position: p1,
	}
	player2 := Player{
		Score:    0,
		Position: p2,
	}

	var rolls, losingScore int
	d := DeterministicDie{1}
	for {
		player1 = player1.Advance(d.Roll() + d.Roll() + d.Roll())
		rolls += 3
		if player1.Score >= 1000 {
			losingScore = player2.Score
			break
		}
		player2 = player2.Advance(d.Roll() + d.Roll() + d.Roll())
		rolls += 3
		if player2.Score >= 1000 {
			losingScore = player1.Score
			break
		}
	}
	fmt.Println(rolls * losingScore)

	initial := State{
		Players:    [2]Player{{0, p1}, {0, p2}},
		NextPlayer: 0,
	}
	pm := make(PathMap)
	known := make(Outcomes)
	toVisit := map[State]bool{
		initial: true,
	}

	dieRollProbs := make(map[int]int)
	for x := 1; x <= 3; x++ {
		for y := 1; y <= 3; y++ {
			for z := 1; z <= 3; z++ {
				dieRollProbs[x+y+z]++
			}
		}
	}

	for len(toVisit) != 0 {
		for previous := range toVisit {
			delete(toVisit, previous)
			if _, ok := pm[previous]; ok {
				// Don't re-visit nodes already visited.
				continue
			}

			pm[previous] = make(map[State]int)
			player := previous.Players[previous.NextPlayer]
			for roll, weight := range dieRollProbs {
				newPlayer := player.Advance(roll)
				newState := previous
				newState.Players[previous.NextPlayer] = newPlayer
				newState.NextPlayer = 1 - previous.NextPlayer

				// Record the edge weight.
				pm[previous][newState] = weight

				// Score the board.
				if newPlayer.Score >= 21 {
					// Record the victory with this leaf node.
					var outcome OutcomeCounts
					outcome[previous.NextPlayer] = uint64(1)
					known[newState] = outcome
					continue
				}
				// Otherwise, we need to record that we need to visit this state.
				toVisit[newState] = true
			}
		}
	}

	for {
		if _, ok := known[initial]; ok {
			break
		}
	outer:
		for pos, edges := range pm {
			if _, ok := known[pos]; ok {
				continue
			}
			var outcomes OutcomeCounts
			for target, weight := range edges {
				if score, ok := known[target]; ok {
					outcomes[0] += score[0] * uint64(weight)
					outcomes[1] += score[1] * uint64(weight)
				} else {
					// We've been asked to evaluated a parent node without its leaves.
					continue outer
				}
			}
			known[pos] = outcomes
		}
	}
	if known[initial][0] > known[initial][1] {
		fmt.Println(known[initial][0])
	} else {
		fmt.Println(known[initial][1])
	}
}
