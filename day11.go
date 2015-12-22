package main

import (
	"fmt"
)

func main() {
	input := []byte("hxbxwxba")
	for ; !valid(input); inc(input) {
	}
	fmt.Println(string(input))
	for inc(input); !valid(input); inc(input) {
	}
	fmt.Println(string(input))
}

func valid(input []byte) bool {
	straight := false
	for i := 0; i < len(input)-2; i++ {
		if input[i]+1 == input[i+1] && input[i+1]+1 == input[i+2] {
			straight = true
		}
	}
	if !straight {
		return false
	}
	for i := range input {
		if input[i] == 'i' || input[i] == 'o' || input[i] == 'l' {
			return false
		}
	}

	pairs := make(map[byte]bool)
	for i := 0; i < len(input)-1; i++ {
		if input[i] == input[i+1] {
			pairs[input[i]] = true
			i++
		}
	}
	return len(pairs) >= 2
}

func inc(input []byte) {
	carry := true
	for i := len(input) - 1; i >= 0 && carry; i-- {
		if input[i] == 'z' {
			input[i] = 'a'
		} else {
			input[i]++
			carry = false
		}
	}
}
