package main

import (
	"container/ring"
	"flag"
	"fmt"
)

var recipes = flag.Int("recipes", 16, "The number of recipes to try or pattern to look for.")
var partB = flag.Bool("partB", false, "Whether to use the Part B logic.")

// Function check determines if the trailing digits match the pattern.
func Check(last *ring.Ring, pattern []int) bool {
	for i := 0; i < len(pattern); i++ {
		if last.Value.(int) != pattern[i] {
			return false
		}
		last = last.Prev()
	}
	return true
}

func main() {
	flag.Parse()

	// Save our starting position as we need start.Prev() to jump to end.
	start := ring.New(2)

	// Initialize our two elven workers to start at the first and second recipes.
	e1 := start
	e1.Value = 3
	e2 := e1.Next()
	e2.Value = 7

	// Expected contains the digits from least to greatest in the part B pattern.
	expected := make([]int, 0)
	for i := *recipes; i != 0; i /= 10 {
		expected = append(expected, i%10)
	}

	// Record how many recipes we've seen so that we can stop on time
	// or so we can know how many in we are for reporting pattern matches.
	n := 2
	// If using part B logic loop forever and use a break; otherwise loop
	// for the specified number of recipes plus be able to look at the 10 after.
	for *partB || n <= *recipes+10 {
		// Get the elves' current scores.
		v1 := e1.Value.(int)
		v2 := e2.Value.(int)

		// Compute the digits of the sum. Note this can only be 2 digits (9+9=18)
		sum := v1 + v2
		sum1 := sum / 10
		sum2 := sum % 10

		// Ignore leading 0s in base 10.
		if sum1 != 0 {
			added := ring.New(1)
			added.Value = sum1
			// Add to the end of our circularly linked list and increment count.
			start.Prev().Link(added)
			n++
			if *partB && Check(added, expected) {
				fmt.Printf("Found match at %d recipes.\n", n-len(expected))
				break
			}
		}

		// Add the ones digit next.
		added := ring.New(1)
		added.Value = sum2
		// Add to the end of our circularly linked list and increment count.
		start.Prev().Link(added)
		n++
		if *partB && Check(added, expected) {
			fmt.Printf("Found match at %d recipes.\n", n-len(expected))
			break
		}

		// Move the elves along by their value plus one, looping as needed.
		for i := 0; i < 1+v1; i++ {
			e1 = e1.Next()
		}
		for i := 0; i < 1+v2; i++ {
			e2 = e2.Next()
		}
	}

	if !*partB {
		// Rewind to the *recipe recipe, then read forward by N
		cur := start
		for i := 0; i <= n-(*recipes)-1; i++ {
			cur = cur.Prev()
		}
		fmt.Printf("The next 10 digits after recipe %d are: ", *recipes)
		for i := 0; i < 10; i++ {
			fmt.Printf("%d", cur.Value.(int))
			cur = cur.Next()
		}
		fmt.Println()
	}
}
