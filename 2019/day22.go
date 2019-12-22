package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

var inputFile = flag.String("inputFile", "inputs/day22.input", "Relative file path to use as input.")

type Deck [10007]int

func main() {
	flag.Parse()
	bytes, err := ioutil.ReadFile(*inputFile)
	if err != nil {
		return
	}

	var d Deck
	for i := range d {
		d[i] = i
	}

	contents := string(bytes)
	split := strings.Split(contents, "\n")
	for _, s := range split {
		if s == "" {
			continue
		}
		tokens := strings.Split(s, " ")
		var newDeck Deck
		switch tokens[0] {
		case "cut":
			count, err := strconv.Atoi(tokens[1])
			if err != nil {
				fmt.Printf("Failed to parse cut count %s\n", tokens[1])
			}
			if count > 0 {
				copy(newDeck[len(d)-count:], d[0:count])
				copy(newDeck[0:len(d)-count], d[count:])
			} else {
				copy(newDeck[0:-count], d[len(d)+count:])
				copy(newDeck[-count:], d[0:len(d)+count])
			}
		case "deal":
			switch tokens[1] {
			case "into":
				for i, v := range d {
					newDeck[len(d)-i-1] = v
				}
			case "with":
				inc, err := strconv.Atoi(tokens[3])
				if err != nil {
					fmt.Printf("Failed to parse increment %s\n", tokens[4])
					return
				}
				for i, v := range d {
					newDeck[(inc*i)%len(d)] = v
				}
			}
		}
		d = newDeck
	}
	for i, v := range d {
		if v == 2019 {
			fmt.Println(i)
			return
		}
	}
}
