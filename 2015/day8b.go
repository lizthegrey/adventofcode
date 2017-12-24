package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	input := 0
	output := 0

	reader := bufio.NewReader(os.Stdin)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		input += len(line) - 1
		output += len(line) - 1 + 2
		for i := 0; i < len(line)-1; i++ {
			if line[i] == '\\' {
				output++
			} else if line[i] == '"' {
				output++
			}
		}
	}

	fmt.Println(input)
	fmt.Println(output)
	fmt.Println(output - input)
}
