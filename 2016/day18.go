package main

import (
	"fmt"
)

type Row []bool

func (r Row) Next() Row {
	ret := make(Row, len(r))
	for i := range r {
		left := !(i == 0 || !r[i-1])
		center := r[i]
		right := !(i == len(r)-1 || !r[i+1])
		ret[i] = (left && center && !right) || (!left && center && right) || (left && !center && !right) || (!left && !center && right)
	}
	return ret
}

func (r Row) CountSafe() int {
	ret := 0
	for _, v := range r {
		if !v {
			ret++
		}
	}
	return ret
}

func main() {
	input := ".^^..^...^..^^.^^^.^^^.^^^^^^.^.^^^^.^^.^^^^^^.^...^......^...^^^..^^^.....^^^^^^^^^....^^...^^^^..^"
	r := make(Row, len(input))
	for i, c := range input {
		r[i] = c == '^'
	}
	sum := 0
	for i := 0; i < 400000; i++ {
		sum += r.CountSafe()
		r = r.Next()
	}
	fmt.Println(sum)
}
