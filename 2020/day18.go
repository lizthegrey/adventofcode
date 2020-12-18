package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

var inputFile = flag.String("inputFile", "inputs/day18.input", "Relative file path to use as input.")
var partB = flag.Bool("partB", false, "Whether to use part B logic.")

type Expression interface {
	Evaluate() int
}

type Number struct {
	Value int
}

func (n Number) Evaluate() int {
	return n.Value
}

type Sequence struct {
	Operands  []Expression
	Operators []Operator
}

func (s Sequence) Evaluate() int {
	var result int
	if !*partB {
		result = s.Operands[0].Evaluate()
		for i, v := range s.Operands {
			if i == 0 {
				continue
			}
			result = s.Operators[i-1](result, v.Evaluate())
		}
	} else {
		i := 1
		expr := s
		for i < len(expr.Operands) {
			// If it's an add...
			if s.Operators[i-1](1, 1) == 2 {
				preVal := expr.Operands[i-1].Evaluate()
				curVal := expr.Operands[i].Evaluate()
				expr.Operands = append(expr.Operands[0:i], expr.Operands[i+1:len(expr.Operands)]...)
				expr.Operands[i-1] = &Number{preVal + curVal}
				expr.Operators = append(expr.Operators[0:i-1], expr.Operators[i:len(expr.Operators)]...)
				continue
			}
			i++
		}

		result = 1
		for _, v := range expr.Operands {
			result *= v.Evaluate()
		}
	}
	return result
}

type Operator func(int, int) int

func Add(x, y int) int {
	return x + y
}
func Mul(x, y int) int {
	return x * y
}

func main() {
	flag.Parse()
	bytes, err := ioutil.ReadFile(*inputFile)
	if err != nil {
		return
	}
	contents := string(bytes)
	contents = strings.ReplaceAll(contents, "(", "( ")
	contents = strings.ReplaceAll(contents, ")", " )")
	split := strings.Split(contents, "\n")
	split = split[:len(split)-1]

	var homework []Expression

	for _, s := range split {
		var line Sequence
		tokens := strings.Split(s, " ")
		curExpr := []*Sequence{&line}
		for _, tok := range tokens {
			peek := curExpr[len(curExpr)-1]
			switch tok {
			case "(":
				child := Sequence{}
				peek.Operands = append(peek.Operands, &child)
				curExpr = append(curExpr, &child)
			case ")":
				curExpr = curExpr[:len(curExpr)-1]
			case "*":
				peek.Operators = append(peek.Operators, Mul)
			case "+":
				peek.Operators = append(peek.Operators, Add)
			default:
				n, err := strconv.Atoi(tok)
				if err != nil {
					fmt.Printf("Failed to parse token %s in line %s\n", tok, s)
					return
				}
				peek.Operands = append(peek.Operands, Number{n})
			}
		}
		homework = append(homework, line)
	}

	sum := 0
	for _, line := range homework {
		sum += line.Evaluate()
	}
	fmt.Println(sum)
}
