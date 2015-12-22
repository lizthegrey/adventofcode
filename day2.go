package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	paper := 0
	ribbon := 0
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		dim := strings.Split(line[:len(line)-1], "x")
		w, _ := strconv.Atoi(dim[0])
		l, _ := strconv.Atoi(dim[1])
		h, _ := strconv.Atoi(dim[2])
		paper += 2 * w * l + 2 * w * h + 2 * l * h
		if w * l < w * h && w * l < l * h {
			paper += w * l
		} else if  w * h < l * h {
			paper += w * h
		} else {
			paper += l * h
		}
		ribbon += w * l * h
		if w+l < w+h && w+l < l+h {
			ribbon += 2*(w+l)
		} else if w+h < l+h {
			ribbon += 2*(w+h)
		} else {
			ribbon += 2*(l+h)
		}
	}
	fmt.Println(paper)
	fmt.Println(ribbon)
}
