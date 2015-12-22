package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	input := 0
	memory := 0

	reader := bufio.NewReader(os.Stdin)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		for i := 0; i < len(line) - 1; i++ {
			input++
			if (i == 0 || i == len(line) - 2) {
				if line[i] != '"' {
					fmt.Println("Malformed input line")
					return
				}
				continue
			}
			memory++
			if line[i] == '\\' {
				if line[i+1] == 'x' {
					input += 3
					i += 3
				} else if line[i+1] == '\\' {
					input += 1
					i += 1
				} else if line[i+1] == '"' {
					input += 1
					i += 1
				}
			}
		}
	}

	fmt.Println(input)
	fmt.Println(memory)
	fmt.Println(input-memory)
}
