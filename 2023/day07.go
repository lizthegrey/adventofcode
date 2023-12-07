package main

import (
	"cmp"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"slices"
	"strconv"
	"strings"
)

var inputFile = flag.String("inputFile", "inputs/day07.input", "Relative file path to use as input.")

type Hand struct {
	Cards [5]uint8
	Bid   int
	score Score
}

type Score int

const (
	Unscored = iota
	HighCard
	OnePair
	TwoPair
	Triple
	FullHouse
	Quadruple
	Quintuple
)

const Wild uint8 = 1

func (h Hand) Score() Score {
	if h.score == Unscored {
		h.score = h.scoreHelper()
	}
	return h.score
}

func (h Hand) scoreHelper() Score {
	distinct := make(map[uint8]int)
	highCount := 0
	var highIdx uint8
	for _, v := range h.Cards {
		count := distinct[v] + 1
		distinct[v] = count
		if v != Wild && count > highCount {
			highCount = count
			highIdx = v
		}
	}

	// Assume we will always get an optimal outcome by adding all wilds to the
	// highest count of non-wild cards.
	distinct[highIdx] += distinct[Wild]
	distinct[Wild] = 0

	var triplePresent, doublePresent bool
	for _, v := range distinct {
		if v == 5 {
			return Quintuple
		}
		if v == 4 {
			return Quadruple
		}
		if v == 3 {
			triplePresent = true
		}
		if v == 2 {
			if doublePresent {
				return TwoPair
			}
			doublePresent = true
		}
	}
	if doublePresent && triplePresent {
		return FullHouse
	}
	if triplePresent {
		return Triple
	}
	if doublePresent {
		return OnePair
	}
	return HighCard
}

func sort(a, b Hand) int {
	ret := cmp.Compare(a.Score(), b.Score())
	if ret != 0 {
		return ret
	}
	for i := 0; i < 5; i++ {
		ret = cmp.Compare(a.Cards[i], b.Cards[i])
		if ret != 0 {
			return ret
		}
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

	var hands []Hand
	for _, s := range split[:len(split)-1] {
		var hand Hand
		for i := 0; i < 5; i++ {
			card := rune(s[i])
			var value uint8
			switch card {
			case 'A':
				value = 14
			case 'K':
				value = 13
			case 'Q':
				value = 12
			case 'J':
				value = 11
			case 'T':
				value = 10
			default:
				value = uint8(card - '0')
			}
			hand.Cards[i] = value
		}
		hand.Bid, err = strconv.Atoi(s[6:])
		if err != nil {
			log.Fatalf("Failed to parse line %s: %v", s, err)
		}
		hands = append(hands, hand)
	}

	slices.SortFunc(hands, sort)
	sum := 0
	for i, v := range hands {
		sum += (i + 1) * v.Bid
	}
	fmt.Println(sum)

	// Part B
	for i, v := range hands {
		for j, c := range v.Cards {
			if c == 11 {
				hands[i].Cards[j] = Wild
			}
		}
		hands[i].score = Unscored
	}

	slices.SortFunc(hands, sort)
	sum = 0
	for i, v := range hands {
		sum += (i + 1) * v.Bid
	}
	fmt.Println(sum)
}
