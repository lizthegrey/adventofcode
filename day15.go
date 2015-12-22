package main

import (
	"fmt"
)

type Ingredient struct {
	Capacity, Durability, Flavor, Texture, Calories int
}
type Recipe map[string]int

func main() {
	ingredients := make(map[string]Ingredient)
	ingredients["Sugar"] = Ingredient{3, 0, 0, -3, 2}
	ingredients["Sprinkles"] = Ingredient{-3, 3, 0, 0, 9}
	ingredients["Candy"] = Ingredient{-1, 0, 4, 0, 1}
	ingredients["Chocolate"] = Ingredient{0, 0, -2, 2, 8}

	mostPoints := 0
	maxVolume := 100

	count := 1
	for i := 0; i < len(ingredients)-1; i++ {
		count *= (maxVolume + 1)
	}

	for i := 0; i < count; i++ {
		encoded := i
		r := make(Recipe)
		total := 0
		for name := range ingredients {
			if len(r) == len(ingredients)-1 {
				r[name] = maxVolume - total
			} else {
				quantity := encoded % (maxVolume + 1)
				encoded /= (maxVolume + 1)
				r[name] = quantity
				total += quantity
			}
		}
		points := evaluateRecipe(r, ingredients)
		if points > mostPoints {
			mostPoints = points
		}
	}

	fmt.Println(mostPoints)
}

func evaluateRecipe(r Recipe, ingredients map[string]Ingredient) int {
	capacity := 0
	durability := 0
	flavor := 0
	texture := 0
	calories := 0
	for n := range r {
		i := ingredients[n]
		quantity := r[n]
		if quantity < 0 {
			return -1
		}
		capacity += i.Capacity * quantity
		durability += i.Durability * quantity
		flavor += i.Flavor * quantity
		texture += i.Texture * quantity
		calories += i.Calories * quantity
	}
	if capacity <= 0 || durability <= 0 || flavor <= 0 || texture <= 0 || calories != 500 {
		return 0
	}
	return capacity * durability * flavor * texture
}
