package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

var inputFile = flag.String("inputFile", "inputs/day07.input", "Relative file path to use as input.")

func main() {
	flag.Parse()
	bytes, err := ioutil.ReadFile(*inputFile)
	if err != nil {
		return
	}
	contents := string(bytes)
	split := strings.Split(contents, "\n")

	var validA, validB int64
	for _, s := range split[:len(split)-1] {
		parts := strings.Split(s, ": ")
		total, _ := strconv.Atoi(parts[0])
		var nums []int64
		for _, part := range strings.Split(parts[1], " ") {
			n, _ := strconv.Atoi(part)
			nums = append(nums, int64(n))
		}
		if check(int64(total), nums[0], nums[1:], false) {
			validA += int64(total)
		}
		if check(int64(total), nums[0], nums[1:], true) {
			validB += int64(total)
		}
	}
	fmt.Println(validA)
	fmt.Println(validB)
}

func check(total, cumulative int64, nums []int64, partB bool) bool {
	// Each op makes total bigger, so bail early if we overshot.
	if cumulative > total {
		return false
	}
	sum := cumulative + nums[0]
	product := cumulative * nums[0]
	cat := concat(cumulative, nums[0])
	// No further numbers to add.
	if len(nums) == 1 {
		if sum == total || product == total || (partB && cat == total) {
			return true
		}
		return false
	}
	return check(total, sum, nums[1:], partB) || check(total, product, nums[1:], partB) || (partB && check(total, cat, nums[1:], partB))
}

func concat(left, right int64) int64 {
	shift := 0
	digits := right
	for digits > 0 {
		digits = digits / 10
		shift++
	}
	ret := left
	for range shift {
		ret *= 10
	}
	ret += right
	return ret
}
