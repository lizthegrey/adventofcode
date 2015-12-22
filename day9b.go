package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
)

type Pair struct {
	x, y string
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	r := regexp.MustCompile("([a-zA-Z]+) to ([a-zA-Z]+) = ([0-9]+)")
	distances := make(map[Pair]uint32)
	places := make(map[string]bool)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		parsed := r.FindStringSubmatch(line)
		src := parsed[1]
		dst := parsed[2]
		distance, _ := strconv.Atoi(parsed[3])
		distances[Pair{src, dst}] = uint32(distance)
		distances[Pair{dst, src}] = uint32(distance)
		places[src] = true
		places[dst] = true
	}

	toVisit := make([]string, len(places))
	i := 0
	for k := range places {
		toVisit[i] = k
		i++
	}

	fmt.Println(visitAll(nil, toVisit, distances))
}

func visitAll(start *string, toVisit []string, distances map[Pair]uint32) uint32 {
	longest := uint32(0)
	for i := range toVisit {
		destination := toVisit[i]
		remainingPlaces := make([]string, len(toVisit)-1)
		copy(remainingPlaces[:i], toVisit[:i])
		copy(remainingPlaces[i:], toVisit[i+1:])

		total := uint32(0)
		if len(remainingPlaces) > 0 {
			total += visitAll(&destination, remainingPlaces, distances)
		}
		if start != nil {
			if distance, ok := distances[Pair{*start, destination}]; ok {
				total += distance
			} else {
				log.Fatalf("Failed to find city pair: %s, %s\n", *start, destination)
			}
		}
		if total > longest {
			longest = total
		}
	}
	return longest
}
