package main

import (
	"flag"
	"fmt"
	"github.com/lizthegrey/adventofcode/2019/intcode"
	"strconv"
	"strings"
)

var inputFile = flag.String("inputFile", "inputs/day17.input", "Relative file path to use as input.")

type Coord struct {
	X, Y int
}

func main() {
	flag.Parse()
	tape := intcode.ReadInput(*inputFile)
	if tape == nil {
		fmt.Println("Failed to parse input.")
		return
	}

	workingTape := tape.Copy()
	output, _ := workingTape.Process(nil)

	var robot Coord
	// up/right/down/left
	var robotDir int
	// Keeps the passable map.
	view := make(map[Coord]bool)

	bytes := make([]byte, 0)
	r := 0
	c := 0
	for n := range output {
		bytes = append(bytes, byte(n))
		loc := Coord{c, r}
		switch rune(n) {
		case '\n':
			r++
			c = 0
			continue
		case 'X':
			// Robot is tumbling.
			robot.X = -1
			robot.Y = -1
		case '.':
			// Nothing to do, empty space.
		case '^':
			robotDir = 0
			robot = loc
			view[loc] = true
		case '>':
			robotDir = 1
			robot = loc
			view[loc] = true
		case 'v':
			robotDir = 2
			robot = loc
			view[loc] = true
		case '<':
			robotDir = 3
			robot = loc
			view[loc] = true
		case '#':
			view[loc] = true
		default:
			fmt.Printf("Unrecognized character %v\n", n)
		}
		c++
	}
	fmt.Println()
	fmt.Println()
	str := string(bytes)
	fmt.Printf(str)

	fmt.Printf("Robot located at row %d, col %d and pointing %d\n", robot.Y, robot.X, robotDir)

	sum := 0
	for l := range view {
		neighbors := [4]Coord{
			// Up
			{l.X, l.Y - 1},
			// Right
			{l.X + 1, l.Y},
			// Down
			{l.X, l.Y + 1},
			// Left
			{l.X - 1, l.Y},
		}
		intersection := true
		for _, n := range neighbors {
			if !view[n] {
				intersection = false
			}
		}
		if intersection {
			sum += l.X * l.Y
		}
	}
	fmt.Println(sum)

	// Calculate the path, expressed as a set of "move forward 1", "turn" instructions.
	directions := make([]rune, 0)
	for {
		neighbors := [4]Coord{
			// Up
			{robot.X, robot.Y - 1},
			// Right
			{robot.X + 1, robot.Y},
			// Down
			{robot.X, robot.Y + 1},
			// Left
			{robot.X - 1, robot.Y},
		}
		if view[neighbors[robotDir]] {
			directions = append(directions, '1')
		} else if view[neighbors[(robotDir+1)%4]] {
			robotDir = (robotDir + 1) % 4
			directions = append(directions, 'R', '1')
		} else if view[neighbors[(robotDir+3)%4]] {
			robotDir = (robotDir + 3) % 4
			directions = append(directions, 'L', '1')
		} else {
			// Found end of path.
			break
		}
		robot = neighbors[robotDir]
	}

	dirs := string(directions)
	fmt.Println(dirs)

	for sMax := 5; sMax < len(dirs)-1; sMax++ {
		for tMin := 5; tMin < len(dirs)-1-sMax; tMin++ {
			start := string(directions[0:sMax])
			tail := string(directions[len(dirs)-tMin : len(dirs)])
			repl := strings.Replace(dirs, start, "A", 10)
			repl = strings.Replace(repl, tail, "C", 10)
			if len(repl) > 70 || len(start) > 100 || len(tail) > 100 {
				continue
			}
			fmt.Printf("%s, %s: %s\n", fold(start), fold(tail), fold(repl))
		}
	}

	// Then chomp off substring combinations that work.
	// Then feed it to the machine.

	input := make(chan int, 1)
	tape[0] = 2
	output, _ = tape.Process(input)
}

func fold(s string) string {
	var output []byte
	var last rune
	count := 1
	for _, d := range s {
		if d != last {
			if last == '1' {
				output = append(output, strconv.Itoa(count)...)
			} else {
				output = append(output, byte(last))
			}
			count = 1
			last = d
		} else {
			count++
		}
	}
	if last == '1' {
		output = append(output, strconv.Itoa(count)...)
	} else {
		output = append(output, byte(last))
	}
	return string(output)
}
