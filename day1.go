package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	line, _ := reader.ReadString('\n')
	floor := 0
	hitNegative := false
	for i := range line[:len(line)-1] {
		switch line[i] {
		case '(':
			floor++
		case ')':
			floor--
		}
		if floor < 0 && !hitNegative {
			hitNegative = true
			fmt.Println(i + 1)
		}
	}
	fmt.Println(floor)
}
