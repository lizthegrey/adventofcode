package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/lizthegrey/adventofcode/2021/trace"
	"go.opentelemetry.io/otel"
	//"go.opentelemetry.io/otel/attribute"
)

var inputFile = flag.String("inputFile", "inputs/day22.input", "Relative file path to use as input.")

var tr = otel.Tracer("day22")

type Coord3 [3]int

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

	cubes := make(map[Coord3]bool)
	for i, s := range split {
		if i >= 20 {
			fmt.Println(len(cubes))
			return
		}
		parts := strings.Split(s, " ")
		coords := strings.Split(parts[1], ",")
		var lower, upper Coord3

		for p, v := range coords {
			innerParts := strings.Split(v, "..")
			lower[p], _ = strconv.Atoi(innerParts[0][2:])
			upper[p], _ = strconv.Atoi(innerParts[1])
		}

		for x := lower[0]; x <= upper[0]; x++ {
			for y := lower[1]; y <= upper[1]; y++ {
				for z := lower[2]; z <= upper[2]; z++ {
					switch parts[0] {
					case "on":
						cubes[Coord3{x, y, z}] = true
					case "off":
						delete(cubes, Coord3{x, y, z})
					}
				}
			}
		}
	}
}
