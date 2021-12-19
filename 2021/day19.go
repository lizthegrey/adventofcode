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
	// "go.opentelemetry.io/otel/attribute"
)

var inputFile = flag.String("inputFile", "inputs/day19.input", "Relative file path to use as input.")

var tr = otel.Tracer("day19")

type Coord3 struct {
	X, Y, Z int
}

// X,Y,Z -> X,Y,Z: [[1,0,0],[0,1,0],[0,0,1]]
// X,Y,Z -> Z,Y,X: [[0,0,1],[0,1,0],[1,0,0]]
type MappingMatrix [3][3]int

func generateAllTransforms() [24]MappingMatrix {
	return [24]MappingMatrix{
		{{1, 0, 0}, {0, 1, 0}, {0, 0, 1}},
		{{1, 0, 0}, {0, 0, -1}, {0, 1, 0}},
		{{1, 0, 0}, {0, -1, 0}, {0, 0, -1}},
		{{1, 0, 0}, {0, 0, 1}, {0, -1, 0}},

		{{0, -1, 0}, {1, 0, 0}, {0, 0, 1}},
		{{0, 0, 1}, {1, 0, 0}, {0, 1, 0}},
		{{0, 1, 0}, {1, 0, 0}, {0, 0, -1}},
		{{0, 0, -1}, {1, 0, 0}, {0, -1, 0}},

		{{-1, 0, 0}, {0, -1, 0}, {0, 0, 1}},
		{{-1, 0, 0}, {0, 0, -1}, {0, -1, 0}},
		{{-1, 0, 0}, {0, 1, 0}, {0, 0, -1}},
		{{-1, 0, 0}, {0, 0, 1}, {0, 1, 0}},

		{{0, 1, 0}, {-1, 0, 0}, {0, 0, 1}},
		{{0, 0, 1}, {-1, 0, 0}, {0, -1, 0}},
		{{0, -1, 0}, {-1, 0, 0}, {0, 0, -1}},
		{{0, 0, -1}, {-1, 0, 0}, {0, 1, 0}},

		{{0, 0, -1}, {0, 1, 0}, {1, 0, 0}},
		{{0, 1, 0}, {0, 0, 1}, {1, 0, 0}},
		{{0, 0, 1}, {0, -1, 0}, {1, 0, 0}},
		{{0, -1, 0}, {0, 0, -1}, {1, 0, 0}},

		{{0, 0, -1}, {0, -1, 0}, {-1, 0, 0}},
		{{0, -1, 0}, {0, 0, 1}, {-1, 0, 0}},
		{{0, 0, 1}, {0, 1, 0}, {-1, 0, 0}},
		{{0, 1, 0}, {0, 0, -1}, {-1, 0, 0}},
	}
}

func (c Coord3) Sub(o Coord3) Coord3 {
	return Coord3{
		X: c.X - o.X,
		Y: c.Y - o.Y,
		Z: c.Z - o.Z,
	}
}

func (c Coord3) Add(o Coord3) Coord3 {
	return Coord3{
		X: c.X + o.X,
		Y: c.Y + o.Y,
		Z: c.Z + o.Z,
	}
}

func (c Coord3) Map(m MappingMatrix) Coord3 {
	return Coord3{
		X: m[0][0]*c.X + m[0][1]*c.Y + m[0][2]*c.Z,
		Y: m[1][0]*c.X + m[1][1]*c.Y + m[1][2]*c.Z,
		Z: m[2][0]*c.X + m[2][1]*c.Y + m[2][2]*c.Z,
	}
}

func (m MappingMatrix) Map(n MappingMatrix) MappingMatrix {
	return MappingMatrix{
		{
			m[0][0]*n[0][0] + m[0][1]*n[1][0] + m[0][2]*n[2][0],
			m[0][0]*n[0][1] + m[0][1]*n[1][1] + m[0][2]*n[2][1],
			m[0][0]*n[0][2] + m[0][1]*n[1][2] + m[0][2]*n[2][2],
		},
		{
			m[1][0]*n[0][0] + m[1][1]*n[1][0] + m[1][2]*n[2][0],
			m[1][0]*n[0][1] + m[1][1]*n[1][1] + m[1][2]*n[2][1],
			m[1][0]*n[0][2] + m[1][1]*n[1][2] + m[1][2]*n[2][2],
		},
		{
			m[2][0]*n[0][0] + m[2][1]*n[1][0] + m[2][2]*n[2][0],
			m[2][0]*n[0][1] + m[2][1]*n[1][1] + m[2][2]*n[2][1],
			m[2][0]*n[0][2] + m[2][1]*n[1][2] + m[2][2]*n[2][2],
		},
	}
}

type Pair [2]int

type Scanner struct {
	AxisMap *MappingMatrix
	Seen    []Coord3
	// Position is relative to scanner 0
	Position *Coord3
	Deltas   map[Coord3]Pair
}

func (s *Scanner) ComputeDeltas() {
	s.Deltas = make(map[Coord3]Pair)
	for i, n := range s.Seen {
		for j, m := range s.Seen {
			if j >= i {
				break
			}
			s.Deltas[n.Sub(m)] = Pair{i, j}
		}
	}
}

func (s Scanner) TransformedSeen() []Coord3 {
	if s.AxisMap == nil {
		return nil
	}
	transformed := make([]Coord3, len(s.Seen))
	for i, v := range s.Seen {
		transformed[i] = v.Map(*s.AxisMap).Add(*s.Position)
	}
	return transformed
}

func main() {
	flag.Parse()

	ctx := context.Background()
	hny, tp := trace.InitializeTracing(ctx)
	defer hny.Shutdown(ctx)
	defer tp.Shutdown(ctx)

	allTransforms := generateAllTransforms()

	bytes, err := ioutil.ReadFile(*inputFile)
	if err != nil {
		return
	}
	contents := string(bytes)
	split := strings.Split(contents, "\n")

	var scanners []*Scanner
	var scanner *Scanner
	for _, s := range split {
		if len(s) == 0 {
			scanner.ComputeDeltas()
			continue
		}
		if s[0:3] == "---" {
			scanner = &Scanner{}
			scanners = append(scanners, scanner)
			continue
		}
		parts := strings.Split(s, ",")
		x, _ := strconv.Atoi(parts[0])
		y, _ := strconv.Atoi(parts[1])
		z, _ := strconv.Atoi(parts[2])
		scanner.Seen = append(scanner.Seen, Coord3{x, y, z})
	}

	scanners[0].AxisMap = &allTransforms[0]
	scanners[0].Position = &Coord3{0, 0, 0}
	beacons := make(map[Coord3]bool)
	for _, v := range scanners[0].TransformedSeen() {
		beacons[v] = true
	}

outer:
	for {
		// Iterate through each sensor that does not yet have a transform and position.
		// We are done when there are no more sensors without positions and axes.
		for _, s := range scanners {
			if s.Position != nil && s.AxisMap != nil {
				continue
			}
			// This is an unaligned beacon.
			for _, aligned := range scanners {
				if aligned.Position == nil || aligned.AxisMap == nil {
					continue
				}
				// Search all aligned beacons to see if we have at least 12 matches.
				deltasMatched := make(map[MappingMatrix][][2]Pair)
				for d, p := range aligned.Deltas {
					for _, mapping := range allTransforms {
						if pair, ok := s.Deltas[d.Map(mapping)]; ok {
							deltasMatched[mapping] = append(deltasMatched[mapping], [2]Pair{p, pair})
							break
						}
					}
				}
				for mapping, pairs := range deltasMatched {
					if len(pairs) < 12 {
						// Keep looking for other orientations that might match.
						continue
					}
					// Success! we've found a orientation that results in >= 12 overlaps.
					mapped := aligned.AxisMap.Map(mapping)
					s.AxisMap = &mapped
					// Take one such overlap pair, and compute a consistent diff.
					pair := pairs[0]
					referenceOne := aligned.Seen[pair[0][0]].Map(*aligned.AxisMap).Add(*aligned.Position)
					referenceTwo := aligned.Seen[pair[0][1]].Map(*aligned.AxisMap).Add(*aligned.Position)
					mineOne := s.Seen[pair[1][0]].Map(*s.AxisMap)
					mineTwo := s.Seen[pair[1][1]].Map(*s.AxisMap)

					position := referenceOne.Sub(mineOne)
					if position.Add(mineTwo) != referenceTwo {
						position = referenceTwo.Sub(mineOne)
						if position.Add(mineTwo) != referenceOne {
							fmt.Printf("inconsistent position: %s, %s\n", position, referenceOne.Sub(mineTwo))
						}
					}
					s.Position = &position

					for _, v := range s.TransformedSeen() {
						beacons[v] = true
					}
					// Now restart the outer loop, since there's new data to work with.
					continue outer
				}
				// If we get here, we've failed to find at least 12 overlaps.
				// Go on and check the next aligned beacon.
			}
			// We were not able to match against any aligned beacons.
		}
		// We have not found any new beacon matches, we are probably done (or stuck).
		break
	}
	for _, s := range scanners {
		if s.AxisMap == nil || s.Position == nil {
			fmt.Println("Incomplete mapping, result will be an undercount.")
		}
		fmt.Println(*s.Position)
	}
	fmt.Println(len(beacons))
}
