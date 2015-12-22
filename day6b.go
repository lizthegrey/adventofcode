package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
)

type Coord struct {
	x, y int
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	r := regexp.MustCompile("(turn off|toggle|turn on) ([0-9]+),([0-9]+) through ([0-9]+),([0-9]+)")
	var on map[Coord]int = make(map[Coord]int)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		parsed := r.FindStringSubmatch(line)
		command := parsed[1]
		x1, _ := strconv.Atoi(parsed[2])
		y1, _ := strconv.Atoi(parsed[3])
		x2, _ := strconv.Atoi(parsed[4])
		y2, _ := strconv.Atoi(parsed[5])
		if command == "turn off" {
			for x := x1; x <= x2; x++ {
				for y := y1; y <= y2; y++ {
					if on[Coord{x, y}] >= 1 {
						on[Coord{x,y}]--
					}
				}
			}
		} else if command == "turn on" {
			for x := x1; x <= x2; x++ {
				for y := y1; y <= y2; y++ {
					on[Coord{x, y}]++
				}
			}
		} else {
			for x := x1; x <= x2; x++ {
				for y := y1; y <= y2; y++ {
					on[Coord{x, y}] += 2
				}
			}
		}
	}
	sum := 0
	for key := range(on) {
		sum += on[key]
	}
	fmt.Println(sum)
}
