package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strings"
)

var inputFile = flag.String("inputFile", "inputs/day05.input", "Relative file path to use as input.")

type Seat struct {
	Row, Column int
}

func (s *Seat) ID() int {
	return s.Row*8 + s.Column
}

func main() {
	flag.Parse()
	bytes, err := ioutil.ReadFile(*inputFile)
	if err != nil {
		return
	}
	contents := string(bytes)
	split := strings.Split(contents, "\n")
	split = split[:len(split)-1]
	seats := make([]Seat, len(split))
	highestSeat := 0
	for n, s := range split {
		// Read bits from MSB to LSB
		if len(s) != (7 + 3) {
			fmt.Printf("Invalid boarding pass length: %s\n", s)
		}
		row := 0
		for i := 6; i >= 0; i-- {
			if s[6-i] == 'B' {
				row += 1 << i
			}
		}

		col := 0
		for i := 2; i >= 0; i-- {
			if s[9-i] == 'R' {
				col += 1 << i
			}
		}
		seats[n] = Seat{row, col}
		if seats[n].ID() > highestSeat {
			highestSeat = seats[n].ID()
		}
	}
	fmt.Println(highestSeat)
	var seen [1024]bool
	for _, v := range seats {
		seen[v.ID()] = true
	}
	gap := false
	init := true
	for i, found := range seen {
		if init && found {
			init = false
			gap = false
			continue
		}
		if gap && found {
			fmt.Println(i - 1)
			break
		}
		gap = !found
	}
}
