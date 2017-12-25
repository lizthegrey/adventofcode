package main

import (
	"fmt"
)

type Tape map[int]bool
type Machine struct {
	S byte
	T Tape
	P int
}

func main() {
	m := Machine{'A', make(Tape), 0}
	for i := 0; i < 12523873; i++ {
		m.Iterate()
	}
	fmt.Println(len(m.T))
}

func (m *Machine) Iterate() {
	switch m.S {
	case 'A':
		if !m.T[m.P] {
			m.T[m.P] = true
			m.P++
			m.S = 'B'
		} else {
			m.T[m.P] = true
			m.P--
			m.S = 'E'
		}
	case 'B':
		if !m.T[m.P] {
			m.T[m.P] = true
			m.P++
			m.S = 'C'
		} else {
			m.T[m.P] = true
			m.P++
			m.S = 'F'
		}
	case 'C':
		if !m.T[m.P] {
			m.T[m.P] = true
			m.P--
			m.S = 'D'
		} else {
			delete(m.T, m.P)
			m.P++
			m.S = 'B'
		}
	case 'D':
		if !m.T[m.P] {
			m.T[m.P] = true
			m.P++
			m.S = 'E'
		} else {
			delete(m.T, m.P)
			m.P--
			m.S = 'C'
		}
	case 'E':
		if !m.T[m.P] {
			m.T[m.P] = true
			m.P--
			m.S = 'A'
		} else {
			delete(m.T, m.P)
			m.P++
			m.S = 'D'
		}
	case 'F':
		if !m.T[m.P] {
			m.T[m.P] = true
			m.P++
			m.S = 'A'
		} else {
			m.T[m.P] = true
			m.P++
			m.S = 'C'
		}
	}
}
