package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"sort"
	"strings"
)

var inputFile = flag.String("inputFile", "inputs/day21.input", "Relative file path to use as input.")

func main() {
	flag.Parse()
	bytes, err := ioutil.ReadFile(*inputFile)
	if err != nil {
		return
	}
	contents := string(bytes)
	split := strings.Split(contents, "\n")
	split = split[:len(split)-1]
	foods := make(map[int]map[string]bool)
	inverseContains := make(map[string]map[int]bool)
	for i, s := range split {
		foods[i] = make(map[string]bool)
		parts := strings.Split(s, " (contains ")
		ingreds := strings.Split(parts[0], " ")
		for _, ing := range ingreds {
			foods[i][ing] = true
		}
		allergens := strings.Split(parts[1][0:len(parts[1])-1], ", ")
		for _, a := range allergens {
			if inverseContains[a] == nil {
				inverseContains[a] = make(map[int]bool)
			}
			inverseContains[a][i] = true
		}
	}
	candidates := make(map[string]map[string]bool)
	for allergen, recipes := range inverseContains {
		// For each allergen, look for a food that is found in all of the
		// recipes containing that allergen.
		ingredientsSeen := make(map[string]int)
		candidates[allergen] = make(map[string]bool)
		for rNum := range recipes {
			recipe := foods[rNum]
			for i := range recipe {
				ingredientsSeen[i]++
			}
		}
		for k, v := range ingredientsSeen {
			if v == len(recipes) {
				// This is a candidate for being our allergen.
				candidates[allergen][k] = true
			}
		}
	}

	foodToAllergenMapping := make(map[string]string)
	allergenToFoodMapping := make(map[string]string)
	for {
		changed := false

		for k, ingredients := range candidates {
			for i := range foodToAllergenMapping {
				delete(ingredients, i)
			}
			if len(ingredients) == 1 {
				changed = true
				delete(candidates, k)
				for i := range ingredients {
					foodToAllergenMapping[i] = k
					allergenToFoodMapping[k] = i
				}
			}
		}

		if !changed {
			break
		}
	}
	var result int
	for _, recipe := range foods {
		for i := range recipe {
			if _, found := foodToAllergenMapping[i]; !found {
				result++
			}
		}
	}
	fmt.Println(result)
	var allergens sort.StringSlice
	for allergen := range allergenToFoodMapping {
		allergens = append(allergens, allergen)
	}
	allergens.Sort()

	dangerous := make([]string, len(allergens))
	for i, allergen := range allergens {
		dangerous[i] = allergenToFoodMapping[allergen]
	}
	fmt.Println(strings.Join(dangerous, ","))
}
