package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
)

type Reindeer struct {
	Speed int
	Endurance int
	Rest int
	Distance int
	Points int
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	r := regexp.MustCompile("([a-zA-Z]+) can fly ([0-9]+) km/s for ([0-9]+) seconds, but then must rest for ([0-9]+) seconds.")
	time := 2503

	reindeer := make(map[string]*Reindeer)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		parsed := r.FindStringSubmatch(line)
		name := parsed[1]
		speed, _ := strconv.Atoi(parsed[2])
		endurance, _ := strconv.Atoi(parsed[3])
		rest, _ := strconv.Atoi(parsed[4])

		reindeer[name] = &Reindeer{speed, endurance, rest, 0, 0}
	}

	for t := 1; t <= time; t++ {
		furthest := 0

		for i := range reindeer {
			r := reindeer[i]
			if (t-1) % (r.Endurance + r.Rest) < r.Endurance {
				r.Distance += r.Speed
			}
			if furthest < r.Distance {
				furthest = r.Distance
			}
		}
		for i := range reindeer {
			r := reindeer[i]
			if r.Distance == furthest {
				r.Points++
			}
		}
	}

	mostPoints := 0
	for i := range reindeer {
		r := reindeer[i]
		if r.Points > mostPoints {
			mostPoints = r.Points
		}
	}
	fmt.Println(mostPoints)
}
