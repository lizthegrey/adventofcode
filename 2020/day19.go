package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

var inputFile = flag.String("inputFile", "inputs/day19.input", "Relative file path to use as input.")
var partB = flag.Bool("partB", false, "Whether to use part B logic.")

type Rule struct {
	Literal  string
	Children [][]int
}

type Ruleset map[int]Rule

// Match returns the number of characters matched.
func (rs Ruleset) Match(ruleNo int, s string) []int {
	r := rs[ruleNo]
	if len(r.Children) == 0 {
		if len(s) < len(r.Literal) {
			return nil
		}
		if s[:len(r.Literal)] == r.Literal {
			return []int{len(r.Literal)}
		}
	}
	var matchedChars []int
	for _, c := range r.Children {
		potentialMatches := []int{0}
		for _, g := range c {
			var newPotentialMatches []int
			// Check each potential prefix.
			// 1, 2 || 3, 4.
			// Inside the loop I'm only dealing with 1->2,
			// but I need to check each potential length of 1
			// as starting points to evaluate 2.
			for _, m := range potentialMatches {
				matches := rs.Match(g, s[m:len(s)])
				if len(matches) == 0 {
					// This isn't a possible match.
					continue
				}
				for _, v := range matches {
					newPotentialMatches = append(newPotentialMatches, v+m)
				}
			}
			potentialMatches = newPotentialMatches
		}
		matchedChars = append(matchedChars, potentialMatches...)
	}
	return matchedChars
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

	rules := make(Ruleset)
	var messages []string
	doingMessages := false
	for _, s := range split {
		if s == "" {
			doingMessages = true
			continue
		}
		if doingMessages {
			messages = append(messages, s)
			continue
		}

		delim := strings.Index(s, ":")
		ruleNo, err := strconv.Atoi(s[0:delim])
		if err != nil {
			fmt.Printf("Failed to parse %s\n", s)
			break
		}
		rule := Rule{}
		if s[delim+2] == '"' {
			rule.Literal = s[delim+3 : len(s)-1]
		} else {
			matches := strings.Split(s[delim+2:len(s)], " | ")
			for _, m := range matches {
				elems := strings.Split(m, " ")
				var group []int
				for _, elem := range elems {
					child, err := strconv.Atoi(elem)
					if err != nil {
						fmt.Printf("Failed to parse %s\n", s)
						break
					}
					group = append(group, child)
				}
				rule.Children = append(rule.Children, group)
			}
		}
		rules[ruleNo] = rule
	}

	if *partB {
		rules[8] = Rule{Children: [][]int{{42}, {42, 8}}}
		rules[11] = Rule{Children: [][]int{{42, 31}, {42, 11, 31}}}
	}

	matches := 0
	for _, m := range messages {
		results := rules.Match(0, m)
		for _, r := range results {
			if r == len(m) {
				matches++
				break
			}
		}
	}
	fmt.Println(matches)
}
