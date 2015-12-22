package main

import (
	"bufio"
	"fmt"
	"os"
)

type Coord struct {
	x, y int64
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	line, _ := reader.ReadString('\n')
	var x1, y1, x2, y2 int64
	var visited map[Coord]int = make(map[Coord]int)
	visited[Coord{0, 0}] = 1
	isRobot := false
	for i := range line {
		if i == len(line)-1 {
			break
		}
		var x, y int64
		if isRobot {
			x = x2
			y = y2
		} else {
			x = x1
			y = y1
		}
		if line[i] == '^' {
			y += 1
		} else if line[i] == '>' {
			x += 1
		} else if line[i] == 'v' {
			y -= 1
		} else if line[i] == '<' {
			x -= 1
		} else {
			fmt.Println("Parse error")
			break
		}
		visited[Coord{x, y}] += 1
		if isRobot {
			x2 = x
			y2 = y
		} else {
			x1 = x
			y1 = y
		}
		isRobot = !isRobot
	}
	fmt.Println(len(visited))
}
