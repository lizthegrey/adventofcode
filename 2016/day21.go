package main

import (
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"strings"
)

func process(instrs []string, in []byte) string {
	output := in
	for _, l := range instrs {
		temp := make([]byte, len(output))
		copy(temp, output)
		split := strings.Split(l, " ")
		switch split[0] {
		case "rotate":
			steps := 0
			switch split[1] {
			case "left":
				l, err := strconv.Atoi(split[2])
				if err != nil {
					log.Fatalf("Bad instruction: %s", l)
				}
				steps = len(output) - l
			case "right":
				r, err := strconv.Atoi(split[2])
				if err != nil {
					log.Fatalf("Bad instruction: %s", l)
				}
				steps = r
			case "based":
				needle := split[6][0]
				for i, c := range output {
					if c == needle {
						steps = i
					}
				}
				if steps >= 4 {
					steps++
				}
				steps++
			default:
				log.Fatalf("Bad instruction: %s", l)
			}
			for i, c := range output {
				temp[(i + steps) % len(output)] = c
			}
		case "reverse":
			from, err := strconv.Atoi(split[2])
			if err != nil {
				log.Fatalf("Bad instruction: %s", l)
			}
			to, err := strconv.Atoi(split[4])
			if err != nil {
				log.Fatalf("Bad instruction: %s", l)
			}
			j := to
			for i := from; i <= to; i++ {
				temp[i] = output[j]
				j--
			}
		case "swap":
			switch split[1] {
			case "position":
				x, err := strconv.Atoi(split[2])
				if err != nil {
					log.Fatalf("Bad instruction: %s", l)
				}
				y, err := strconv.Atoi(split[5])
				if err != nil {
					log.Fatalf("Bad instruction: %s", l)
				}
				temp[x] = output[y]
				temp[y] = output[x]
			case "letter":
				x := split[2][0]
				y := split[5][0]
				for i, c := range output {
					if c == x {
						temp[i] = y
					} else if c == y {
						temp[i] = x
					} else {
						temp[i] = c
					}
				}
			default:
				log.Fatalf("Bad instruction: %s", l)
			}
		case "move":
			from, err := strconv.Atoi(split[2])
			if err != nil {
				log.Fatalf("Bad instruction: %s", l)
			}
			to, err := strconv.Atoi(split[5])
			if err != nil {
				log.Fatalf("Bad instruction: %s", l)
			}
			temp[to] = output[from]
			offset := 0
			for i, c := range output {
				if i == from {
					offset -= 1
					continue
				} else if i + offset == to {
					offset += 1
				}
				temp[i + offset] = c
			}
		default:
			log.Fatalf("Bad instruction: %s", l)
		}
		output = temp
	}
	return string(output)
}

func main() {
	fmt.Println(process([]string{
		"swap position 4 with position 0",
		"swap letter d with letter b",
		"reverse positions 0 through 4",
		"rotate left 1 step",
		"move position 1 to position 4",
		"move position 3 to position 0",
		"rotate based on position of letter b",
		"rotate based on position of letter d",
	}, []byte("abcde")))

	input := strings.Split(`move position 0 to position 3
rotate right 0 steps
rotate right 1 step
move position 1 to position 5
swap letter h with letter b
reverse positions 1 through 3
swap letter a with letter g
swap letter b with letter h
rotate based on position of letter c
swap letter d with letter c
rotate based on position of letter c
swap position 6 with position 5
rotate right 7 steps
swap letter b with letter h
move position 4 to position 3
swap position 1 with position 0
swap position 7 with position 5
move position 7 to position 1
swap letter c with letter a
move position 7 to position 5
rotate right 4 steps
swap position 0 with position 5
move position 3 to position 1
swap letter c with letter h
rotate based on position of letter d
reverse positions 0 through 2
rotate based on position of letter g
move position 6 to position 7
move position 2 to position 5
swap position 1 with position 0
swap letter f with letter c
rotate right 1 step
reverse positions 2 through 4
rotate left 1 step
rotate based on position of letter h
rotate right 1 step
rotate right 5 steps
swap position 6 with position 3
move position 0 to position 5
swap letter g with letter f
reverse positions 2 through 7
reverse positions 4 through 6
swap position 4 with position 1
move position 2 to position 1
move position 3 to position 1
swap letter b with letter a
rotate based on position of letter b
reverse positions 3 through 5
move position 0 to position 2
rotate based on position of letter b
reverse positions 4 through 5
rotate based on position of letter g
reverse positions 0 through 5
swap letter h with letter c
reverse positions 2 through 5
swap position 7 with position 5
swap letter g with letter d
swap letter d with letter e
move position 1 to position 2
move position 3 to position 2
swap letter d with letter g
swap position 3 with position 7
swap letter b with letter f
rotate right 3 steps
move position 5 to position 3
move position 1 to position 2
rotate based on position of letter b
rotate based on position of letter c
reverse positions 2 through 3
move position 2 to position 3
rotate right 1 step
move position 7 to position 0
rotate right 3 steps
move position 6 to position 3
rotate based on position of letter e
swap letter c with letter b
swap letter f with letter d
swap position 2 with position 5
swap letter f with letter g
rotate based on position of letter a
reverse positions 3 through 4
rotate left 7 steps
rotate left 6 steps
swap letter g with letter b
reverse positions 3 through 6
rotate right 6 steps
rotate based on position of letter c
rotate based on position of letter b
rotate left 1 step
reverse positions 3 through 7
swap letter f with letter g
swap position 4 with position 1
rotate based on position of letter d
move position 0 to position 4
swap position 7 with position 6
rotate right 6 steps
rotate based on position of letter e
move position 7 to position 3
rotate right 3 steps
swap position 1 with position 2`, "\n")
	fmt.Println(process(input, []byte("abcdefgh")))
	for {
		chars := []byte("abcdefgh")
		shuffled := rand.Perm(len(chars))
		test := make([]byte, len(chars))
		for i, j := range shuffled {
			test[j] = chars[i]
		}
		if process(input, test) == "fbgdceah" {
			fmt.Println(string(test))
			break
		}
	}
}
