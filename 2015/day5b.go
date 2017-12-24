package main

import (
	"bufio"
	"fmt"
	"os"
)

type Pair struct {
	first, second rune
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	nice := 0
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		penult := rune('+')
		last := rune('-')
		repeat := make(map[Pair]int)
		doubledouble := false
		alternate := false
		for i := 0; i < len(line)-1; i++ {
			cur := rune(line[i])
			pair := Pair{last, cur}
			if repeat[pair] != i-1 && repeat[pair] != 0 {
				doubledouble = true
			}
			if repeat[pair] == 0 {
				repeat[pair] = i
			}
			if cur == penult && last != cur {
				alternate = true
			}
			penult = last
			last = cur
		}
		fmt.Printf("%s: %t, %t\n", line[:len(line)-1], doubledouble, alternate)
		if doubledouble && alternate {
			nice += 1
		}
	}
	fmt.Println(nice)
}
