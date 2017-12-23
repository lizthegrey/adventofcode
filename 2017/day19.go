package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strings"
)

var inputFile = flag.String("inputFile", "inputs/day19.input", "Relative file path to use as input.")

type Loc struct {
	Passable, ChangeDir bool
	Letter              rune
}

func main() {
	flag.Parse()

	bytes, err := ioutil.ReadFile(*inputFile)
	if err != nil {
		fmt.Printf("Could not open file %s because %v.\n", *inputFile, err)
	}
	contents := string(bytes)
	lines := strings.Split(contents, "\n")

	height := len(lines)
	width := len(lines[0])

	board := make([][]Loc, height)
	for y := 0; y < height; y++ {
		board[y] = make([]Loc, width)
	}

	for y, l := range lines {
		if len(l) == 0 {
			break
		}
		for x, c := range l {
			switch c {
			case ' ':
				board[y][x] = Loc{false, false, ' '}
			case '+':
				board[y][x] = Loc{true, true, ' '}
			case '|':
				fallthrough
			case '-':
				board[y][x] = Loc{true, false, ' '}
			default:
				board[y][x] = Loc{true, false, c}
			}
		}
	}

	var packetX, packetY int
	// 0 = N, 1 = E, 2 = S, 3 = W
	packetDir := 2
	for x, c := range board[0] {
		if c.Passable {
			packetX = x
			break
		}
	}

	for {
		switch packetDir {
		case 0:
			packetY--
		case 1:
			packetX++
		case 2:
			packetY++
		case 3:
			packetX--
		}
		if board[packetY][packetX].ChangeDir {
			switch packetDir {
			case 0:
				fallthrough
			case 2:
				if board[packetY][packetX+1].Passable {
					packetDir = 1
				} else {
					packetDir = 3
				}
			case 1:
				fallthrough
			case 3:
				if board[packetY+1][packetX].Passable {
					packetDir = 2
				} else {
					packetDir = 0
				}
			}
		}

		if board[packetY][packetX].Letter != ' ' {
			fmt.Printf("%c", board[packetY][packetX].Letter)
		}
		if !board[packetY][packetX].Passable {
			break
		}
	}
	fmt.Println()
}
