package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strings"
	"unsafe"
)

var inputFile = flag.String("inputFile", "inputs/day03.input", "Relative file path to use as input.")

type Pair struct {
	arr       *int
	remaining int
}
type Memo map[Pair]uint64

func main() {
	flag.Parse()
	bytes, err := ioutil.ReadFile(*inputFile)
	if err != nil {
		return
	}
	contents := string(bytes)
	split := strings.Split(contents, "\n")

	var sumA, sumB uint64
	for _, s := range split[:len(split)-1] {
		var arr []int
		for _, r := range s {
			arr = append(arr, int(r-'0'))
		}
		memo := make(Memo)
		sumA += memo.max(arr, 2)
		sumB += memo.max(arr, 12)
	}
	fmt.Println(sumA)
	fmt.Println(sumB)
}

func (m Memo) max(arr []int, remaining int) uint64 {
	if remaining == 0 {
		return 0
	}
	key := Pair{unsafe.SliceData(arr), remaining}
	if v, ok := m[key]; ok {
		return v
	}

	var highest uint64
	place := uint64(1)
	for _ = range remaining - 1 {
		place *= 10
	}
	for pos, v := range arr {
		if len(arr)-pos < remaining {
			continue
		}
		if candidate := place*uint64(v) + m.max(arr[pos+1:], remaining-1); candidate > highest {
			highest = candidate
		}
	}
	m[key] = highest
	return highest
}
