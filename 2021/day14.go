package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"math"
	"strings"

	"github.com/lizthegrey/adventofcode/2021/trace"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

var inputFile = flag.String("inputFile", "inputs/day14.input", "Relative file path to use as input.")

var tr = otel.Tracer("day14")

type Pair struct {
	Left, Right byte
}
type Polymer struct {
	InnerElems map[Pair]int
	Outer      Pair
}
type Rules map[Pair]byte

func (r Rules) Grow(ctx context.Context, input Polymer) Polymer {
	_, sp := tr.Start(ctx, "grow")
	defer sp.End()

	pairs := make(map[Pair]int)
	for pair, count := range input.InnerElems {
		if insertion, ok := r[pair]; ok {
			pairs[Pair{pair.Left, insertion}] += count
			pairs[Pair{insertion, pair.Right}] += count
		}
	}
	return Polymer{InnerElems: pairs, Outer: input.Outer}
}

func main() {
	flag.Parse()

	ctx := context.Background()
	hny, tp := trace.InitializeTracing(ctx)
	defer hny.Shutdown(ctx)
	defer tp.Shutdown(ctx)

	bytes, err := ioutil.ReadFile(*inputFile)
	if err != nil {
		return
	}
	contents := string(bytes)
	split := strings.Split(contents, "\n")
	split = split[:len(split)-1]

	ctx, sp := tr.Start(context.Background(), "solution")
	defer sp.End()

	_, iSp := tr.Start(ctx, "parsing")
	polymerRaw := split[0]
	polymer := Polymer{
		InnerElems: make(map[Pair]int),
		Outer:      Pair{polymerRaw[0], polymerRaw[len(polymerRaw)-1]},
	}
	for i := 0; i < len(polymerRaw)-1; i++ {
		polymer.InnerElems[Pair{polymerRaw[i], polymerRaw[i+1]}]++
	}
	rules := make(Rules)
	for _, s := range split[2:] {
		parts := strings.Split(s, " -> ")
		rules[Pair{parts[0][0], parts[0][1]}] = parts[1][0]
	}
	iSp.SetAttributes(
		attribute.Int("input_length", len(polymerRaw)),
		attribute.Int("rules_length", len(rules)),
	)
	iSp.End()

	fmt.Println(runMany(ctx, rules, polymer, 10))
	fmt.Println(runMany(ctx, rules, polymer, 40))
}

func runMany(ctx context.Context, rules Rules, polymer Polymer, iterations int) int {
	for i := 0; i < iterations; i++ {
		polymer = rules.Grow(ctx, polymer)
	}
	return getFrequencies(polymer)
}

func getFrequencies(polymer Polymer) int {
	frequencies := make(map[byte]int)
	for p, c := range polymer.InnerElems {
		frequencies[p.Left] += c
		frequencies[p.Right] += c
	}
	frequencies[polymer.Outer.Left]++
	frequencies[polymer.Outer.Right]++

	most := 0
	least := math.MaxInt
	for _, v := range frequencies {
		if v < least {
			least = v
		}
		if v > most {
			most = v
		}
	}
	return (most - least) / 2
}
