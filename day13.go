package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
)

type Pair struct {
	x, y string
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	r := regexp.MustCompile("([a-zA-Z]+) would (gain|lose) ([0-9]+) happiness units by sitting next to ([a-zA-Z]+).")
	happiness := make(map[Pair]int)
	people := make(map[string]bool)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		parsed := r.FindStringSubmatch(line)
		x := parsed[1]
		y := parsed[4]
		happy, _ := strconv.Atoi(parsed[3])
		if parsed[2] == "lose" {
			happy = -happy
		}
		happiness[Pair{x, y}] = happy
		people[x] = true
		people[y] = true
	}
	people["me"] = true

	toVisit := make([]string, len(people))
	i := 0
	for k := range people {
		toVisit[i] = k
		i++
	}

	fmt.Println(visitAll(nil, nil, toVisit, happiness))
}

func visitAll(cur, last *string, toVisit []string, distances map[Pair]int) int {
	longest := 0
	for i := range toVisit {
		destination := toVisit[i]
		myLast := last
		if myLast == nil {
			myLast = &destination
		}
		remainingPlaces := make([]string, len(toVisit)-1)
		copy(remainingPlaces[:i], toVisit[:i])
		copy(remainingPlaces[i:], toVisit[i+1:])

		total := 0
		if len(remainingPlaces) > 0 {
			total += visitAll(&destination, myLast, remainingPlaces, distances)
		} else {
			total += distances[Pair{destination, *myLast}]
			total += distances[Pair{*myLast, destination}]
		}
		if cur != nil {
			total += distances[Pair{*cur, destination}]
			total += distances[Pair{destination, *cur}]
		}
		if total > longest {
			longest = total
			// fmt.Println(destination, remainingPlaces, total)
		}
	}
	return longest
}
