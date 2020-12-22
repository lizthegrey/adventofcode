package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

var inputFile = flag.String("inputFile", "inputs/day22.input", "Relative file path to use as input.")
var partB = flag.Bool("partB", false, "Whether to use part B logic.")

type Cards []int

func main() {
	flag.Parse()
	bytes, err := ioutil.ReadFile(*inputFile)
	if err != nil {
		return
	}
	contents := string(bytes)
	split := strings.Split(contents, "\n")
	split = split[:len(split)-1]

	var decks []Cards
	playerNo := 0
	for _, s := range split {
		if s == "" {
			playerNo++
			continue
		}
		if strings.HasPrefix(s, "Player ") {
			decks = append(decks, nil)
			continue
		}
		n, err := strconv.Atoi(s)
		if err != nil {
			fmt.Printf("Failed to parse %s\n", s)
			break
		}
		decks[playerNo] = append(decks[playerNo], n)
	}

	gameWinner := PlayGame(decks)

	var score int
	for i, c := range decks[gameWinner] {
		score += (len(decks[gameWinner]) - i) * c
	}
	fmt.Println(score)
}

func ModifyDeck(decks []Cards, drawn Cards, winner int) {
	var sorted Cards
	sorted = append(sorted, drawn[winner])
	for i := range decks {
		if i == winner {
			continue
		}
		sorted = append(sorted, drawn[i])
	}
	decks[winner] = append(decks[winner], sorted...)
}

func PlayGame(decks []Cards) int {
	var seen [][]Cards
	for {
	history:
		for _, previous := range seen {
			for playerNo, d := range decks {
				if len(d) != len(previous[playerNo]) {
					continue history
				}
				for i, elem := range d {
					if previous[playerNo][i] != elem {
						// This isn't an exact repeat.
						continue history
					}
				}
			}
			// If we got this far, we have an exact match.
			if *partB {
				return 0
			}
		}
		var cloned []Cards
		for _, d := range decks {
			newDeck := make(Cards, len(d))
			copy(newDeck, d)
			cloned = append(cloned, newDeck)
		}
		seen = append(seen, cloned)

		for i, d := range decks {
			if len(d) == 0 {
				return (i + 1) % 2
			}
		}
		var inPlay Cards
		for playerNo, d := range decks {
			inPlay = append(inPlay, d[0])
			decks[playerNo] = d[1:len(d)]
		}

		if *partB {
			eligibleForRecursion := true
			for playerNo, d := range decks {
				if len(d) < inPlay[playerNo] {
					eligibleForRecursion = false
					break
				}
			}
			if eligibleForRecursion {
				var clones []Cards
				for playerNo, d := range decks {
					newDeck := make(Cards, inPlay[playerNo])
					copy(newDeck, d[0:inPlay[playerNo]])
					clones = append(clones, newDeck)
				}
				winner := PlayGame(clones)
				ModifyDeck(decks, inPlay, winner)
				continue
			}
		}

		var winner int
		var highCard int
		for playerNo, c := range inPlay {
			if c > highCard {
				highCard = c
				winner = playerNo
			}
		}
		ModifyDeck(decks, inPlay, winner)
	}
	return -1
}
