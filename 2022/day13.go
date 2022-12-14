package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"sort"
	"strconv"
	"strings"
)

var inputFile = flag.String("inputFile", "inputs/day13.input", "Relative file path to use as input.")

type elem struct {
	children []elem
	leaf     int
}

func (l elem) display() {
	if l.children == nil {
		fmt.Printf("%d", l.leaf)
		return
	}
	fmt.Printf("[")
	for i, c := range l.children {
		if i > 0 {
			fmt.Printf(",")
		}
		c.display()
	}
	fmt.Printf("]")
}

func (l elem) compare(r elem) *bool {
	var ret bool
	// Convert the case of mismatched input types to both being lists.
	// Note: this is safe to do because it's an value receiver on elem not pointer on *elem.
	if l.children == nil && r.children != nil {
		l.children = []elem{{leaf: l.leaf}}
	} else if l.children != nil && r.children == nil {
		r.children = []elem{{leaf: r.leaf}}
	}

	if l.children == nil && r.children == nil {
		if l.leaf < r.leaf {
			ret = true
		} else if l.leaf > r.leaf {
			ret = false
		} else {
			return nil
		}
	} else {
		for i := 0; i < len(l.children); i++ {
			if i == len(r.children) {
				// Got to the end of the right list first.
				ret = false
				return &ret
			}
			// Compare the item and return the result, if any.
			intermediate := l.children[i].compare(r.children[i])
			if intermediate != nil {
				return intermediate
			}
		}
		if len(l.children) == len(r.children) {
			// The lists are the same length. We can't make a decision.
			return nil
		}
		// Got to the end of the left list first.
		ret = true
	}
	return &ret
}

func NewElem(input string) *elem {
	stack := []elem{{children: make([]elem, 0)}}
	for i := 1; i < len(input)-1; i++ {
		switch input[i] {
		case '[':
			child := elem{
				children: make([]elem, 0),
			}
			stack = append(stack, child)
		case ']':
			stack[len(stack)-2].children = append(stack[len(stack)-2].children, stack[len(stack)-1])
			stack = stack[:len(stack)-1]
		case ',':
			// Do nothing, the existing code handles this.
		default:
			// Otherwise, this is a number.
			end := len(input) - 1
			for tmp := i; tmp < len(input)-1; tmp++ {
				digit := input[tmp]
				if digit < '0' || digit > '9' {
					end = tmp
					break
				}
			}
			val, _ := strconv.Atoi(input[i:end])
			stack[len(stack)-1].children = append(stack[len(stack)-1].children, elem{
				leaf: val,
			})
			// Prepare to resume reading input from the next value.
			i = end - 1
		}
	}
	return &stack[0]
}

func main() {
	flag.Parse()
	bytes, err := ioutil.ReadFile(*inputFile)
	if err != nil {
		return
	}
	contents := string(bytes)
	split := strings.Split(contents, "\n")

	// part A
	var sum int
	for i := 0; i < len(split); i += 3 {
		left := NewElem(split[i])
		right := NewElem(split[i+1])
		if result := left.compare(*right); result != nil && *result {
			sum += 1 + i/3
		}
	}
	fmt.Println(sum)

	// part B
	first := NewElem("[[2]]")
	second := NewElem("[[6]]")
	signals := []*elem{first, second}
	for _, l := range split {
		if len(l) == 0 {
			continue
		}
		signals = append(signals, NewElem(l))
	}
	sort.Slice(signals, func(i, j int) bool {
		ret := signals[i].compare(*signals[j])
		return ret != nil && *ret
	})

	product := 1
	for i, v := range signals {
		if v == first || v == second {
			product *= i + 1
		}
	}
	fmt.Println(product)
}
