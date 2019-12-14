package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"math"
	"strconv"
	"strings"
)

var inputFile = flag.String("inputFile", "inputs/day14.input", "Relative file path to use as input.")
var quantity = flag.Int64("quantity", 1, "Quantity of fuel to produce.")

type Ingredient string

type Recipe struct {
	Inputs         map[Ingredient]int64
	OutputQuantity int64
}

const ore Ingredient = "ORE"
const fuel Ingredient = "FUEL"

func main() {
	flag.Parse()
	bytes, err := ioutil.ReadFile(*inputFile)
	if err != nil {
		return
	}
	contents := string(bytes)
	split := strings.Split(contents, "\n")

	// Assumes only one way to make each ingredient
	recipes := make(map[Ingredient]Recipe)
	for _, s := range split {
		if s == "" {
			continue
		}

		parts := strings.Split(s, " => ")

		inputs := make(map[Ingredient]int64)
		for _, v := range strings.Split(parts[0], ", ") {
			input := strings.Split(v, " ")
			inputQuantity, err := strconv.Atoi(input[0])
			if err != nil {
				fmt.Printf("Failed to parse %s\n", s)
				break
			}
			inputType := Ingredient(input[1])
			inputs[inputType] = int64(inputQuantity)
		}

		output := strings.Split(parts[1], " ")
		outputQuantity, err := strconv.Atoi(output[0])
		if err != nil {
			fmt.Printf("Failed to parse %s\n", s)
			break
		}
		outputType := Ingredient(output[1])
		if _, ok := recipes[outputType]; ok {
			fmt.Printf("Invariant invalid: more than one way to make ingredient %d.\n", outputType)
		}
		recipes[outputType] = Recipe{inputs, int64(outputQuantity)}
	}

	fmt.Println(process(*quantity, recipes))
	fmt.Println(binSearch(func(i int64) bool {
		return process(i, recipes) < int64(1000000000000)
	}, 1, math.MaxInt32))
}

func binSearch(f func(i int64) bool, min, max int64) int64 {
	midpoint := int64(int64(min)/2 + int64(max)/2)
	if min == max || min+1 == max {
		return min
	}
	fmt.Printf("f(%d) = %v\n", midpoint, f(midpoint))
	if f(midpoint) {
		return binSearch(f, midpoint, max)
	} else {
		return binSearch(f, min, midpoint)
	}
}

func process(desired int64, recipes map[Ingredient]Recipe) int64 {
	spare := make(map[Ingredient]int64)
	needed := map[Ingredient]int64{
		fuel: desired,
	}
	for {
		if len(needed) == 1 && needed[ore] != 0 {
			return needed[ore]
		}
		for k, q := range needed {
			if k == ore {
				continue
			}
			r := recipes[k]
			copies := int64(math.Ceil(float64(q-spare[k]) / float64(r.OutputQuantity)))
			for i, q := range r.Inputs {
				needed[i] += copies * q
			}
			spare[k] += (r.OutputQuantity * copies) - q
			delete(needed, k)
		}
	}
}
