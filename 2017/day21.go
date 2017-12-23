package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strings"
)

var inputFile = flag.String("inputFile", "inputs/day21.input", "Relative file path to use as input.")
var cycles = flag.Int("cycles", 5, "The number of cycles to simulate.")

type Pattern2 [2][2]bool
type Pattern3 [3][3]bool
type Pattern4 [4][4]bool

func main() {
	flag.Parse()

	bytes, err := ioutil.ReadFile(*inputFile)
	if err != nil {
		fmt.Printf("Could not open file %s because %v.\n", *inputFile, err)
	}
	contents := string(bytes[:len(bytes)-1])
	lines := strings.Split(contents, "\n")

	book2 := make(map[Pattern2]Pattern3)
	book3 := make(map[Pattern3]Pattern4)

	// Convert input file into patterns, and add to map.
	for _, l := range lines {
		switch len(l) {
		case 20:
			p2, err := ParseP2(l[0:5])
			if err != nil {
				fmt.Println(err)
			}
			p3, err := ParseP3(l[9:20])
			if err != nil {
				fmt.Println(err)
			}
			book2[p2] = p3
		case 34:
			p3, err := ParseP3(l[0:11])
			if err != nil {
				fmt.Println(err)
			}
			p4, err := ParseP4(l[15:34])
			if err != nil {
				fmt.Println(err)
			}
			book3[p3] = p4
		default:
			fmt.Printf("Failed to parse line %s", l)
		}
	}

	// Generate all possible patterns using flipping and rotation.
	for len(book2) != 1<<4 {
		for k, v := range book2 {
			book2[k.FlipX()] = v
			book2[k.FlipY()] = v
			book2[k.Rotate()] = v
		}
	}

	for len(book3) != 1<<9 {
		for k, v := range book3 {
			book3[k.FlipX()] = v
			book3[k.FlipY()] = v
			book3[k.Rotate()] = v
		}
	}

	seed := [][]bool{
		[]bool{false, true, false},
		[]bool{false, false, true},
		[]bool{true, true, true},
	}

	board := seed

	for i := 0; i < *cycles; i++ {
		if len(board)%2 == 0 {
			sw := make([][]Pattern3, len(board)/2)
			for r := 0; r < len(board)/2; r++ {
				sw[r] = make([]Pattern3, len(board)/2)
				for c := 0; c < len(board)/2; c++ {
					sw[r][c] = book2[Swatch2(board, r*2, c*2)]
				}
			}

			newBoard := make([][]bool, 3*len(board)/2)
			for r := 0; r < len(newBoard); r++ {
				newBoard[r] = make([]bool, 3*len(board)/2)
				for c := 0; c < len(newBoard[r]); c++ {
					newBoard[r][c] = sw[r/3][c/3][r%3][c%3]
				}
			}
			board = newBoard
		} else {
			sw := make([][]Pattern4, len(board)/3)
			for r := 0; r < len(board)/3; r++ {
				sw[r] = make([]Pattern4, len(board)/3)
				for c := 0; c < len(board)/3; c++ {
					sw[r][c] = book3[Swatch3(board, r*3, c*3)]
				}
			}

			newBoard := make([][]bool, 4*len(board)/3)
			for r := 0; r < len(newBoard); r++ {
				newBoard[r] = make([]bool, 4*len(board)/3)
				for c := 0; c < len(newBoard[r]); c++ {
					newBoard[r][c] = sw[r/4][c/4][r%4][c%4]
				}
			}
			board = newBoard
		}
	}

	on := 0
	for i := 0; i < len(board); i++ {
		for j := 0; j < len(board[i]); j++ {
			if board[i][j] {
				fmt.Printf("#")
				on++
			} else {
				fmt.Printf(".")
			}
		}
		fmt.Println()
	}
	fmt.Println(on)
}

func Swatch2(search [][]bool, r, c int) Pattern2 {
	var p Pattern2
	for i := 0; i < len(p); i++ {
		for j := 0; j < len(p[i]); j++ {
			p[i][j] = search[r+i][c+j]
		}
	}
	return p
}

func Swatch3(search [][]bool, r, c int) Pattern3 {
	var p Pattern3
	for i := 0; i < len(p); i++ {
		for j := 0; j < len(p[i]); j++ {
			p[i][j] = search[r+i][c+j]
		}
	}
	return p
}

func (p Pattern2) FlipX() Pattern2 {
	var ret Pattern2
	for i := 0; i < len(p); i++ {
		for j := 0; j < len(p[i]); j++ {
			ret[i][j] = p[i][len(p[i])-1-j]
		}
	}
	return ret
}

func (p Pattern2) FlipY() Pattern2 {
	var ret Pattern2
	for i := 0; i < len(p); i++ {
		for j := 0; j < len(p[i]); j++ {
			ret[i][j] = p[len(p)-1-i][j]
		}
	}
	return ret
}

// 1 2 | 3 1
// 3 4 | 4 2
func (p Pattern2) Rotate() Pattern2 {
	var ret Pattern2
	for i := 0; i < len(p); i++ {
		for j := 0; j < len(p[i]); j++ {
			ret[i][j] = p[len(p)-1-j][i]
		}
	}
	return ret
}

func (p Pattern3) FlipX() Pattern3 {
	var ret Pattern3
	for i := 0; i < len(p); i++ {
		for j := 0; j < len(p[i]); j++ {
			ret[i][j] = p[i][len(p[i])-1-j]
		}
	}
	return ret
}

func (p Pattern3) FlipY() Pattern3 {
	var ret Pattern3
	for i := 0; i < len(p); i++ {
		for j := 0; j < len(p[i]); j++ {
			ret[i][j] = p[len(p)-1-i][j]
		}
	}
	return ret
}

// 1 2 3 | 7 4 1
// 4 5 6 | 8 5 2
// 7 8 9 | 9 6 3
func (p Pattern3) Rotate() Pattern3 {
	var ret Pattern3
	for i := 0; i < len(p); i++ {
		for j := 0; j < len(p[i]); j++ {
			ret[i][j] = p[len(p)-1-j][i]
		}
	}
	return ret
}

func ParseP2(s string) (Pattern2, error) {
	var p2 Pattern2
	if len(s) != 5 {
		return p2, fmt.Errorf("Failed to parse pattern %s", s)
	}
	for i := 0; i < 5; i++ {
		if i%3 == 2 {
			continue
		}
		on := s[i] == '#'
		p2[i/3][i%3] = on
	}
	return p2, nil
}

func ParseP3(s string) (Pattern3, error) {
	var p3 Pattern3
	if len(s) != 11 {
		return p3, fmt.Errorf("Failed to parse pattern %s", s)
	}
	for i := 0; i < 11; i++ {
		if i%4 == 3 {
			continue
		}
		on := s[i] == '#'
		p3[i/4][i%4] = on
	}
	return p3, nil
}

func ParseP4(s string) (Pattern4, error) {
	var p4 Pattern4
	if len(s) != 19 {
		return p4, fmt.Errorf("Failed to parse pattern %s", s)
	}
	for i := 0; i < 19; i++ {
		if i%5 == 4 {
			continue
		}
		on := s[i] == '#'
		p4[i/5][i%5] = on
	}
	return p4, nil
}
