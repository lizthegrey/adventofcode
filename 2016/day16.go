package main

import (
	"fmt"
)

type BitString []bool

func (a BitString) Generate() BitString {
	ret := make(BitString, len(a)*2+1)
	copy(ret, a)
	ret[len(a)] = false
	for i := 1; i <= len(a); i++ {
		ret[len(a)+i] = !a[len(a)-i]
	}
	return ret
}

func (a BitString) Check() BitString {
	if len(a)%2 == 1 {
		return a
	}
	out := make([]bool, len(a)/2)
	for i := 0; i < len(a); i += 2 {
		out[i/2] = a[i] == a[i+1]
	}
	return out
}

func (a BitString) Print() {
	for _, b := range a {
		if b {
			fmt.Printf("1")
		} else {
			fmt.Printf("0")
		}
	}
	fmt.Println()
}

func main() {
	length := 35651584
	working := BitString{true, false, false, false, true, true, true, false, false, true, true, true, true, false, false, false, false}
	for ; len(working) < length; working = working.Generate() {
		//working.Print()
	}
	working = working[0:length]
	for c := working; ; c = c.Check() {
		//c.Print()
		if len(c)%2 == 1 {
			c.Print()
			return
		}
	}
}
