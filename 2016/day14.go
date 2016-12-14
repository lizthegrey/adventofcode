package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
)

type Repeateds struct {
	FirstTrip, Quints uint16
}

type Cache map[int]Repeateds

var Salt string = "qzyelonm"
var StretchFactor int = 2016

func (cache Cache) GetRepeateds(i int) Repeateds {
	if r, ok := cache[i]; ok {
		return r
	}
	data := []byte(fmt.Sprintf("%s%d", Salt, i))
	h := md5.Sum(data)

	for i := 0; i < StretchFactor; i++ {
		h = md5.Sum([]byte(hex.EncodeToString(h[:])))
	}

	var r Repeateds

	var consecChar byte
	consecRepeats := 0
	for _, b := range h {
		eval := func(c byte) {
			if consecChar != c {
				consecRepeats = 1
				consecChar = c
			} else {
				consecRepeats++
				if consecRepeats == 3 && r.FirstTrip == 0 {
					r.FirstTrip = 1 << consecChar
				} else if consecRepeats == 5 {
					r.Quints |= 1 << consecChar
				}
			}
		}

		lowHalf := b & 0xf
		highHalf := b >> 4

		eval(highHalf)
		eval(lowHalf)
	}
	cache[i] = r
	return r
}

func main() {
	cache := make(Cache)
	searchPos := 0
	for found := 0; found < 64; searchPos++ {
		r := cache.GetRepeateds(searchPos)
		if r.FirstTrip == 0 {
			continue
		}
		for i := searchPos + 1; i < searchPos+1001; i++ {
			if cache.GetRepeateds(i).Quints&r.FirstTrip != 0 {
				fmt.Println(searchPos)
				found++
				break
			}
		}
	}
}
