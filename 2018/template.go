package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	// "strconv"
)

var inputFile = flag.String("inputFile", "inputs/dayXX.input", "Relative file path to use as input.")
var partB = flag.Bool("partB", false, "Whether to use the Part B logic.")

func main() {
	flag.Parse()
	f, err := os.Open(*inputFile)
	if err != nil {
		return
	}
	defer f.Close()

	r := bufio.NewReader(f)
	for {
		l, err := r.ReadString('\n')
		if err != nil || len(l) == 0 {
			break
		}
		l = l[:len(l)-1]
		fmt.Println(l)
	}

	result := ""
	fmt.Printf("Result is %s\n", result)
}
