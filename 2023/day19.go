package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
)

var inputFile = flag.String("inputFile", "inputs/day19.input", "Relative file path to use as input.")

type Part [4]int

type Workflow struct {
	Steps []Step
	Final string
}

type Key int

const (
	X = iota
	M
	A
	S
)

type Operator int

const (
	Lt = iota
	Gt
)

type Step struct {
	Disposition string
	CondVal     int
	CondKey     Key
	Cond        Operator
}

func (w Workflow) Process(p Part) string {
	for _, v := range w.Steps {
		comparator := p[v.CondKey]
		switch v.Cond {
		case Lt:
			if comparator < v.CondVal {
				return v.Disposition
			}
		case Gt:
			if comparator > v.CondVal {
				return v.Disposition
			}
		}
	}
	return w.Final
}

type Constraint struct {
	// Open set, not closed set
	Min, Max [4]int
}

const start = "in"
const accepted = "A"

var Unconstrained = Constraint{
	[4]int{0, 0, 0, 0},
	[4]int{4001, 4001, 4001, 4001},
}

func (c Constraint) Possibilities() uint64 {
	ret := uint64(1)
	for i := 0; i < 4; i++ {
		// Closed sets are weird.
		ret *= uint64((c.Max[i] - 1) - (c.Min[i] + 1) + 1)
	}
	return ret
}

func (c Constraint) Combine(o Constraint) *Constraint {
	var ret Constraint
	for i := 0; i < 4; i++ {
		ret.Min[i] = max(c.Min[i], o.Min[i])
		ret.Max[i] = min(c.Max[i], o.Max[i])
		if ret.Max[i]-ret.Min[i] < 2 {
			// between 1 and 3 in an integral open set there's 2. but 1 and 2, or 1 and 1, or 1 and 0 is no good.
			return nil
		}
	}
	return &ret
}

func (s Step) ToConstraint() Constraint {
	ret := Unconstrained
	switch s.Cond {
	case Lt:
		ret.Max[s.CondKey] = s.CondVal
	case Gt:
		ret.Min[s.CondKey] = s.CondVal
	}
	return ret
}

func (s Step) InverseConstraint() Constraint {
	ret := Unconstrained
	switch s.Cond {
	case Lt:
		ret.Min[s.CondKey] = s.CondVal - 1
	case Gt:
		ret.Max[s.CondKey] = s.CondVal + 1
	}
	return ret
}

func main() {
	flag.Parse()
	bytes, err := ioutil.ReadFile(*inputFile)
	if err != nil {
		return
	}
	contents := string(bytes)
	split := strings.Split(contents, "\n")

	queues := make(map[string][]Part)
	workflows := make(map[string]Workflow)

	var partsMode bool
	for _, s := range split[:len(split)-1] {
		if s == "" {
			// Switch to part processing mode.
			partsMode = true
			continue
		}
		if partsMode {
			kvs := strings.Split(s[1:len(s)-1], ",")
			var part Part
			for _, kv := range kvs {
				parts := strings.Split(kv, "=")
				n, err := strconv.Atoi(parts[1])
				if err != nil {
					log.Fatalf("Failed to parse kv %s", kv)
				}
				switch parts[0] {
				case "x":
					part[X] = n
				case "m":
					part[M] = n
				case "a":
					part[A] = n
				case "s":
					part[S] = n
				default:
					log.Fatalf("Unknown key in kv %s", kv)
				}
			}
			queues[start] = append(queues[start], part)
			continue
		}
		// Otherwise this is a workflow.
		parts := strings.Split(s[:len(s)-1], "{")
		var w Workflow
		instructions := strings.Split(parts[1], ",")
		for i, raw := range instructions {
			if i == len(instructions)-1 {
				w.Final = raw
				break
			}
			var s Step
			switch raw[0] {
			case 'x':
				s.CondKey = X
			case 'm':
				s.CondKey = M
			case 'a':
				s.CondKey = A
			case 's':
				s.CondKey = S
			default:
				log.Fatalf("Unknown key in step %s", raw)
			}
			switch raw[1] {
			case '<':
				s.Cond = Lt
			case '>':
				s.Cond = Gt
			default:
				log.Fatalf("Unknown op in step %s", raw)
			}
			final := strings.Split(raw[2:], ":")
			val, err := strconv.Atoi(final[0])
			if err != nil {
				log.Fatalf("Failed to parse step %s", raw)
			}
			s.CondVal = val
			s.Disposition = final[1]
			w.Steps = append(w.Steps, s)
		}

		name := parts[0]
		workflows[name] = w
	}

	for {
		var updated int
		for name, workflow := range workflows {
			for len(queues[name]) > 0 {
				updated++
				item := queues[name][0]
				queues[name] = queues[name][1:]
				disposition := workflow.Process(item)
				queues[disposition] = append(queues[disposition], item)
			}
		}
		if updated == 0 {
			break
		}
	}

	var sum int
	for _, k := range queues[accepted] {
		for _, v := range k {
			sum += v
		}
	}
	fmt.Println(sum)

	// Part B
	constraintQueues := map[string][]Constraint{
		start: []Constraint{Unconstrained},
	}
	for {
		var updated int
		for name, w := range workflows {
		outer:
			for len(constraintQueues[name]) > 0 {
				updated++
				item := &constraintQueues[name][0]
				constraintQueues[name] = constraintQueues[name][1:]
				for _, step := range w.Steps {
					if overlap := item.Combine(step.ToConstraint()); overlap != nil {
						constraintQueues[step.Disposition] = append(constraintQueues[step.Disposition], *overlap)
					}
					item = item.Combine(step.InverseConstraint())
					if item == nil {
						continue outer
					}
				}
				constraintQueues[w.Final] = append(constraintQueues[w.Final], *item)
			}
		}
		if updated == 0 {
			break
		}
	}
	var possibilities uint64
	for _, path := range constraintQueues[accepted] {
		possibilities += path.Possibilities()
	}
	fmt.Println(possibilities)
}
