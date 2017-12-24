package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	r := regexp.MustCompile("([a-zA-Z]+) can fly ([0-9]+) km/s for ([0-9]+) seconds, but then must rest for ([0-9]+) seconds.")
	time := 2503

	furthest := 0
	bestName := ""
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

		cycle := endurance + rest
		distance := time / cycle * (speed * endurance)
		leftover := time % cycle
		if leftover >= endurance {
			distance += endurance * speed
		} else {
			distance += leftover * speed
		}

		if furthest < distance {
			furthest = distance
			bestName = name
		}
	}

	fmt.Println(furthest, bestName)
}
