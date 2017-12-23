package knot

import (
	"container/ring"
	"fmt"
)

func Key(input string) []int {
	ls := make([]int, len(input)+5)
	for i, l := range input {
		ls[i] = int(l)
	}
	addedLengths := [5]int{17, 31, 73, 47, 23}
	for i := 0; i < len(addedLengths); i++ {
		ls[i+len(input)] = addedLengths[i]
	}
	return ls
}

func Hash(length, rounds int, key []int) *ring.Ring {
	loop := ring.New(length)
	for i := 0; i < length; i++ {
		set(loop, i)
		loop = loop.Next()
	}

	pos := loop
	skip := 0
	offset := 0

	for i := 0; i < rounds; i++ {
		pos = round(pos, length, key, &skip, &offset)
	}

	final := pos.Move(-offset)
	return final
}

func Densify(final *ring.Ring) []int {
	length := final.Len() / 16
	ret := make([]int, length)
	for i := 0; i < length; i++ {
		intermediate := 0
		for j := 0; j < 16; j++ {
			digit := Get(final)
			final = final.Next()
			intermediate ^= digit
		}
		ret[i] = intermediate
	}
	return ret
}

// round performs a round of the "hashing" given the specified constraints.
// Note that this modifies the input ring (danger will robinson)
func round(r *ring.Ring, length int, ls []int, skip, offset *int) *ring.Ring {
	pos := r
	for _, n := range ls {
		if n != 0 {
			leftover := pos.Move(n)
			working := pos.Prev().Unlink(n)
			reversed := reverse(working)

			if n == length {
				pos = reversed
			} else {
				pos = leftover.Prev().Link(reversed)
			}
		}
		pos = pos.Move(*skip)
		*offset = (*offset + *skip + n) % length
		*skip++
	}
	return pos
}

func set(r *ring.Ring, v int) {
	r.Value = v
}

func Get(r *ring.Ring) int {
	return r.Value.(int)
}

func Debug(r *ring.Ring) {
	pos := r
	for i := 0; i < r.Len(); i++ {
		fmt.Printf("%d,", Get(pos))
		pos = pos.Next()
	}
	fmt.Println()
}

func reverse(r *ring.Ring) *ring.Ring {
	// [A], B, C, D -> [D], C, B, A
	pos := r.Prev()
	ret := ring.New(r.Len())
	for i := 0; i < r.Len(); i++ {
		set(ret, Get(pos))
		ret = ret.Next()
		pos = pos.Prev()
	}
	return ret
}
