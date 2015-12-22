package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	nice := 0
	vowels := map[rune]bool{
		'a': true,
		'e': true,
		'i': true,
		'o': true,
		'u': true,
	}
	blacklist := map[rune]rune{
		'a': 'b',
		'c': 'd',
		'p': 'q',
		'x': 'y',
	}
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		last := rune('-')
		vowelCount := 0
		double := false
		blacklisted := false
		for i := 0; i < len(line)-1; i++ {
			cur := rune(line[i])
			if last == cur {
				double = true
			}

			if vowels[cur] {
				vowelCount += 1
			}
			if blacklist[last] == cur {
				blacklisted = true
			}
			last = cur
		}
		fmt.Printf("%s: %d, %t, %t\n", line[:len(line)-1], vowelCount, double, blacklisted)
		if vowelCount >= 3 && double && !blacklisted {
			nice += 1
		}
	}
	fmt.Println(nice)
}
