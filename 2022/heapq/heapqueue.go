package heapq

import (
	"container/heap"
)

// Like the heap construct, but with generic types and an upsert operation.
type HeapQueue[T comparable] struct {
	elems     *[]T
	score     map[T]int
	positions map[T]int
}

func New[T comparable]() *HeapQueue[T] {
	return &HeapQueue[T]{
		elems:     &[]T{},
		score:     make(map[T]int),
		positions: make(map[T]int),
	}
}

func (h HeapQueue[T]) Len() int           { return len(*h.elems) }
func (h HeapQueue[T]) Less(i, j int) bool { return h.score[(*h.elems)[i]] < h.score[(*h.elems)[j]] }
func (h *HeapQueue[T]) Swap(i, j int) {
	h.positions[(*h.elems)[i]], h.positions[(*h.elems)[j]] = h.positions[(*h.elems)[j]], h.positions[(*h.elems)[i]]
	(*h.elems)[i], (*h.elems)[j] = (*h.elems)[j], (*h.elems)[i]
}

func (h *HeapQueue[T]) Push(x interface{}) {
	h.positions[x.(T)] = len(*h.elems)
	*h.elems = append(*h.elems, x.(T))
}

func (h *HeapQueue[T]) PopSafe() T {
	return heap.Pop(h).(T)
}

func (h *HeapQueue[T]) Pop() interface{} {
	old := *h.elems
	n := len(old)
	x := old[n-1]
	*h.elems = old[0 : n-1]
	delete(h.positions, x)
	return x
}

func (h *HeapQueue[T]) Upsert(n T, score int) {
	h.score[n] = score
	if pos, ok := h.positions[n]; !ok {
		heap.Push(h, n)
	} else {
		heap.Fix(h, pos)
	}
}
