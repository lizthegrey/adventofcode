package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	r := regexp.MustCompile("(-?[0-9]+)")
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		sum := 0
		matches := r.FindAllString(line, -1)
		for i := range matches {
			val, _ := strconv.Atoi(matches[i])
			sum += val
		}
		fmt.Println(sum)
	}
}
