package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

var inputFile = flag.String("inputFile", "inputs/day16.input", "Relative file path to use as input.")
var debug = flag.Bool("debug", false, "Whether to print debug output along the way.")

type Constraint struct {
	Low, High int
}
type Fields map[string][]Constraint
type Ticket []int
type FieldOrder map[string]int

func main() {
	flag.Parse()
	bytes, err := ioutil.ReadFile(*inputFile)
	if err != nil {
		return
	}
	contents := string(bytes)
	split := strings.Split(contents, "\n")
	split = split[:len(split)-1]
	i := 0
	fields := make(Fields)
	for {
		s := split[i]
		if s == "" {
			break
		}
		parts := strings.Split(s, ": ")
		field := parts[0]
		values := strings.Split(parts[1], " or ")
		for _, v := range values {
			limits := strings.Split(v, "-")
			var c Constraint
			c.Low, err = strconv.Atoi(limits[0])
			if err != nil {
				fmt.Printf("Failed to parse %s\n", s)
				return
			}
			c.High, err = strconv.Atoi(limits[1])
			if err != nil {
				fmt.Printf("Failed to parse %s\n", s)
				return
			}
			fields[field] = append(fields[field], c)
		}
		i++
	}
	if *debug {
		fmt.Println(fields)
	}

	i++
	if split[i] != "your ticket:" {
		fmt.Println("Failed to find my ticket.")
	}
	i++
	my := parseTicket(split[i])
	i += 2
	if split[i] != "nearby tickets:" {
		fmt.Println("Failed to find nearby tickets.")
	}
	i++

	var nearby []Ticket
	for ; i < len(split); i++ {
		nearby = append(nearby, parseTicket(split[i]))
	}
	if *debug {
		fmt.Println(my)
		fmt.Println(nearby)
	}

	var invalid int
	var validTickets []Ticket
	for _, t := range nearby {
		validTicket := true
		for _, v := range t {
			validField := false
			for _, f := range fields {
				if valueValid(v, f) {
					validField = true
					break
				}
			}
			if !validField {
				invalid += v
				validTicket = false
			}
		}
		if validTicket {
			validTickets = append(validTickets, t)
		}
	}
	fmt.Println(invalid)

	// part B
	if *debug {
		fmt.Println(len(validTickets))
	}

	valuesAtPos := make(map[int][]int)
	for _, t := range validTickets {
		for p, v := range t {
			valuesAtPos[p] = append(valuesAtPos[p], v)
		}
	}

	// Loop through each position looking for candidate field names.
	candidateMappings := make(map[string]map[int]bool)
	reverseMappings := make(map[int]map[string]bool)
	for pos, vs := range valuesAtPos {
	outer:
		for name, f := range fields {
			for _, v := range vs {
				if !valueValid(v, f) {
					continue outer
				}
			}
			// This field passes validation at that pos for all tickets.
			if candidateMappings[name] == nil {
				candidateMappings[name] = make(map[int]bool)
			}
			candidateMappings[name][pos] = true
			if reverseMappings[pos] == nil {
				reverseMappings[pos] = make(map[string]bool)
			}
			reverseMappings[pos][name] = true
		}
	}

	if *debug {
		fmt.Println(candidateMappings)
	}

	fo := make(FieldOrder)
	for {
		changed := false
		for name, mappings := range candidateMappings {
			if len(mappings) != 1 {
				continue
			}

			pos := -1
			for k := range mappings {
				pos = k
			}
			fo[name] = pos
			if *debug {
				fmt.Printf("Mapped pos %d to %s\n", pos, name)
			}
			delete(candidateMappings, name)
			for _, remainingMapping := range candidateMappings {
				delete(remainingMapping, pos)
			}
			delete(reverseMappings, pos)
			for _, remainingMapping := range reverseMappings {
				delete(remainingMapping, name)
			}
			changed = true
			break
		}

		if !changed {
			break
		}
	}

	if *debug {
		fmt.Println(fo)
	}

	// Print the desired output.
	product := 1
	for k, idx := range fo {
		if len(k) >= 9 && k[0:9] == "departure" {
			product *= my[idx]
		}
	}
	fmt.Println(product)
}

func valueValid(v int, constraints []Constraint) bool {
	for _, c := range constraints {
		if c.Low <= v && c.High >= v {
			return true
		}
	}
	return false
}

func parseTicket(line string) Ticket {
	var t Ticket
	for _, v := range strings.Split(line, ",") {
		value, err := strconv.Atoi(v)
		if err != nil {
			fmt.Printf("Failed to parse %s\n", line)
			return nil
		}
		t = append(t, value)
	}
	return t
}
