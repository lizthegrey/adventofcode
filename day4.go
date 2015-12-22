package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"strconv"
)

func main() {
	for i := 0; ; i++ {
		h := md5.New()
		io.WriteString(h, "yzbqklnj")
		io.WriteString(h, strconv.Itoa(i))
		if fmt.Sprintf("%x", h.Sum(nil))[0:6] == "000000" {
			fmt.Println(i)
			break
		}
	}
}
