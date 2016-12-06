package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
)

func main() {
	found := 0
	result := []byte("________")
	scrambled := true
	for j := 0; found != 8; j++ {
		data := []byte(fmt.Sprintf("%s%d", "cxdnnyjw", j))
		h := md5.Sum(data)
		hash := hex.EncodeToString(h[:])
		if hash[0] == '0' && hash[1] == '0' && hash[2] == '0' && hash[3] == '0' && hash[4] == '0' {
			if !scrambled {
				result[found] = hash[5]
			} else if pos := hash[5] - '0'; pos >= 0 && pos < 8 && result[pos] == '_' {
				result[pos] = hash[6]
			} else {
				continue
			}
			fmt.Println(string(result))
			found++
		}
	}
}
