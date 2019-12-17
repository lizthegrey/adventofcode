package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strconv"
)

var inputFile = flag.String("inputFile", "inputs/day16.input", "Relative file path to use as input.")

const phases int = 100
const repeats int = 10000

type fft []int8

var template [4]int = [4]int{0, 1, 0, -1}

func (f fft) SecondHalf() fft {
	ret := make(fft, len(f))

	// 00000000000000001
	// 00000000000000011
	// 00000000000000111
	// ...
	memo := int8(0)
	// O(N)
	for p := len(f) - 1; p >= len(f)/2; p-- {
		// O(1)
		memo += f[p]
		memo %= 10
		ret[p] = memo
	}

	return ret
}

func (f fft) Phase() fft {
	ret := f.SecondHalf()

	for p := range f {
		if p >= len(f)/2 {
			break
		}

		// O(N**2)
		pattern := make(fft, len(f))
		repeats := p
		tIdx := 0
		for i := 0; i <= len(pattern); i++ {
			// O(N)
			if repeats == -1 {
				repeats = p
				tIdx = (tIdx + 1) % len(template)
			}
			repeats--
			if i == 0 {
				continue
			}
			pattern[i-1] = int8(template[tIdx])
		}

		sum := int64(0)
		for i, v := range f {
			switch pattern[i] {
			case 0:
				continue
			case 1:
				sum += int64(v)
			case -1:
				sum -= int64(v)
			}
		}
		if sum < 0 {
			sum *= -1
		}
		sum %= 10
		ret[p] = int8(sum)
	}
	return ret
}

func main() {
	flag.Parse()
	bytes, err := ioutil.ReadFile(*inputFile)
	if err != nil {
		return
	}
	contents := string(bytes[:len(bytes)-1])
	digits := make(fft, len(contents))
	for i, d := range contents {
		n, err := strconv.Atoi(string(d))
		if err != nil {
			fmt.Printf("Failed to parse digit %s\n", string(d))
		}
		digits[i] = int8(n)
	}

	digitsA := make(fft, len(digits))
	copy(digitsA, digits)
	for p := 1; p <= phases; p++ {
		digitsA = digitsA.Phase()
	}
	for _, d := range digitsA[:8] {
		fmt.Printf("%d", d)
	}
	fmt.Println()

	offset := 0
	for _, d := range digits[0:7] {
		offset *= 10
		offset += int(d)
	}
	fmt.Printf("Offset: %d\n", offset)
	if offset < (len(digits)*repeats)/2 {
		fmt.Println("Can't support an offset smaller than N/2.")
		return
	}
	digitsB := make(fft, len(digits)*repeats)
	for n := 0; n < repeats; n++ {
		copy(digitsB[n*len(digits):(n+1)*len(digits)], digits)
	}
	for p := 1; p <= phases; p++ {
		digitsB = digitsB.SecondHalf()
	}
	for _, d := range digitsB[offset : offset+8] {
		fmt.Printf("%d", d)
	}
	fmt.Println()
}
