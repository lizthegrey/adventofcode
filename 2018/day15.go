package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"sort"
)

var inputFile = flag.String("inputFile", "inputs/day15.input", "Relative file path to use as input.")
var maxRounds = flag.Int("maxRounds", -1, "Maximum number of rounds to run.")
var verbose = flag.Bool("verbose", false, "Whether to print verbose debug output.")
var partB = flag.Bool("partB", false, "Whether to use the Part B logic.")

type Coord struct {
	X, Y int
}

type CandidateDest struct {
	Dest      Coord
	FirstStep *Coord
	Steps     int
}
type CDByRO []CandidateDest

func (a CDByRO) Len() int           { return len(a) }
func (a CDByRO) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a CDByRO) Less(i, j int) bool { return LessRO(a[i].Dest, a[j].Dest) }

type Unit struct {
	Force bool // true = elf, false is goblin
	HP    int
	Pos   Coord
}

type AttackOrder []*Unit

func (a AttackOrder) Len() int      { return len(a) }
func (a AttackOrder) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a AttackOrder) Less(i, j int) bool {
	// Always push nulls to the right.
	if a[i] == nil {
		return false
	}
	if a[j] == nil {
		return true
	}

	// Sort by HP from smallest to largest.
	if a[i].HP != a[j].HP {
		return a[i].HP < a[j].HP
	}
	// Then sort by Reading Order from smallest to largest.
	return LessRO(a[i].Pos, a[j].Pos)
}

func (u *Unit) PickAdjacentTarget(us Units) *Unit {
	p := u.Pos

	up := p
	up.Y--
	right := p
	right.X++
	down := p
	down.Y++
	left := p
	left.X--

	ao := []*Unit{us[up], us[right], us[down], us[left]}
	for i, t := range ao {
		if t != nil {
			if t.Force == u.Force {
				// Don't attack our friends.
				ao[i] = nil
			}
		}
	}
	sort.Sort(AttackOrder(ao))
	return ao[0] // Could return null if nothing to attack.
}

func (p Coord) PassableNeighbors(t Terrain, us Units) []Coord {
	up := p
	up.Y--
	right := p
	right.X++
	down := p
	down.Y++
	left := p
	left.X--

	ret := make([]Coord, 0)

	// always prioritize up, then left, then right, then down given
	// multiple paths to same destination.

	for _, c := range []Coord{up, left, right, down} {
		if !t[c] {
			// Not a passable tile.
			continue
		}
		if us[c] != nil {
			// Occupied by another unit.
			continue
		}
		ret = append(ret, c)
	}
	return ret
}

type Terrain map[Coord]bool // Passable has value true; not passable not present.
type Units map[Coord]*Unit
type CByRO []Coord
type UByRO []*Unit

func LessRO(a, b Coord) bool {
	if a.Y != b.Y {
		return a.Y < b.Y
	}
	return a.X < b.X
}

func (a CByRO) Len() int           { return len(a) }
func (a CByRO) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a CByRO) Less(i, j int) bool { return LessRO(a[i], a[j]) }

func (a UByRO) Len() int           { return len(a) }
func (a UByRO) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a UByRO) Less(i, j int) bool { return LessRO(a[i].Pos, a[j].Pos) }

func main() {
	flag.Parse()
	f, err := os.Open(*inputFile)
	if err != nil {
		return
	}
	defer f.Close()

	t := make(Terrain)
	uByC := make(Units)
	uByRO := make(UByRO, 0)

	r := bufio.NewReader(f)
	maxY := 0
	maxX := 0
	for y := 0; ; y++ {
		l, err := r.ReadString('\n')
		if err != nil || len(l) == 0 {
			break
		}
		maxY = y

		l = l[:len(l)-1]
		for x, c := range l {
			if c == ' ' {
				break
			}

			maxX = x
			loc := Coord{x, y}
			if c == '#' {
				continue
			}
			// Passable terrain.
			t[loc] = true

			if c == '.' {
				continue
			}

			// We have a member of a force.
			u := Unit{c == 'E', 200, loc}
			uByC[loc] = &u
			uByRO = append(uByRO, &u)
		}
	}

	if !*partB {
		RunRounds(3, uByRO, t, uByC, maxX, maxY)

		// Print final state.
		PrintBoard(t, uByC, maxX, maxY)
		return
	}

	// We are in part B.
	for p := 4; p <= 200; p++ {
		pUByC := make(Units)
		pUByRO := make(UByRO, len(uByRO))

		// Clone our starting position.
		for i, v := range uByRO {
			u := *v
			pUByRO[i] = &u
			pUByC[u.Pos] = &u
		}

		if RunRounds(p, pUByRO, t, pUByC, maxX, maxY) {
			PrintBoard(t, uByC, maxX, maxY)
			fmt.Println(p)
			return
		}
	}
	fmt.Println("Couldn't find a viable elf power.")
}

func PrintBoard(t Terrain, uByC Units, maxX, maxY int) {
	for y := 0; y <= maxY; y++ {
		hps := make([]int, 0)
		for x := 0; x <= maxX; x++ {
			c := Coord{x, y}
			if !t[c] {
				fmt.Printf("#")
				continue
			}
			if u := uByC[c]; u != nil {
				if u.Force {
					fmt.Printf("E")
				} else {
					fmt.Printf("G")
				}
				hps = append(hps, u.HP)
			} else {
				fmt.Printf(".")
			}
		}
		for _, v := range hps {
			fmt.Printf(" %d", v)
		}
		fmt.Println()
	}
}

func RunRounds(elfPower int, uByRO UByRO, t Terrain, uByC Units, maxX, maxY int) bool {
	for rounds := 0; *maxRounds < 0 || rounds < *maxRounds; rounds++ {
		if *verbose {
			fmt.Printf("Starting round %d...\n", rounds+1)
		}
		sort.Sort(uByRO)

		survivors := make(UByRO, 0)
		for _, u := range uByRO {
			// Tick the unit if it wasn't killed earlier in turn.
			if u.HP <= 0 {
				continue
			}

			foundAliveTarget := false
			for _, v := range uByC {
				if v.Force != u.Force {
					foundAliveTarget = true
				}
			}

			if !foundAliveTarget {
				fmt.Printf("Combat ends after %d full rounds\n", rounds)
				winner := ""
				if u.Force {
					winner = "Elves"
				} else {
					winner = "Goblins"
				}
				totalHP := 0
				for _, c := range uByC {
					totalHP += c.HP
				}
				fmt.Printf("%s win with %d total hit points left\n", winner, totalHP)
				fmt.Printf("Outcome: %d * %d = %d\n", rounds, totalHP, rounds*totalHP)
				return true
			}

			// Move towards a target if we're not next to something already.
			target := u.PickAdjacentTarget(uByC)
			if target == nil {
				// Then we need to move towards a square. Nothing's next to us.
				destinations := make(map[Coord]bool)

				for k, v := range uByC {
					if v.Force == u.Force {
						continue
					}
					for _, v := range k.PassableNeighbors(t, uByC) {
						destinations[v] = true
					}
				}

				// Rank the destinations: BFS outwards, until we touch a destination.
				// Map locations to the first direction to travel to get there.
				// However, this doesn't account for ties so we need to resolve those.

				// First step towards that coord.
				seen := make(map[Coord]bool)
				dests := make([]CandidateDest, 0)

				// always prioritize up, then left, then right, then down given
				// multiple paths to same destination.

				// Perform a breadth-first search, stopping after we're latched.
				for worklist := []CandidateDest{{u.Pos, nil, 0}}; len(worklist) != 0; {
					w := worklist[0]
					worklist = worklist[1:]

					if len(dests) > 0 && w.Steps > dests[0].Steps {
						// We've gone past everything with the same distance and can
						// resolve ties.
						break
					}

					if destinations[w.Dest] {
						dests = append(dests, w)
						continue
					}
					for _, v := range w.Dest.PassableNeighbors(t, uByC) {
						if seen[v] {
							continue
						}

						var firstStep Coord
						// Propagate our starting direction.
						if w.FirstStep == nil {
							firstStep = v
						} else {
							firstStep = *w.FirstStep
						}

						entry := CandidateDest{v, &firstStep, w.Steps + 1}

						// Explore it next.
						seen[v] = true
						worklist = append(worklist, entry)
					}
				}

				// Move towards the chosen destination.
				if len(dests) != 0 {
					sort.Sort(CDByRO(dests))
					dest := dests[0]

					delete(uByC, u.Pos)
					u.Pos = *dest.FirstStep
					uByC[u.Pos] = u
				}

				// Re-compute targets after moving.
				target = u.PickAdjacentTarget(uByC)
			}

			// Attack an adjacent target if available.
			if target == nil {
				continue
			}
			if u.Force {
				target.HP -= elfPower
			} else {
				target.HP -= 3
			}
			if target.HP <= 0 {
				// Remove target from uByC to make the square passable.
				// No need to clean up from uByRO, we skip it if it's invalid.
				delete(uByC, target.Pos)
				if target.Force && *partB {
					// an elf died, no joy.
					return false
				}
			}
		}

		// Propagate and tabulate only survivors.
		for _, u := range uByC {
			survivors = append(survivors, u)
		}
		uByRO = survivors
	}
	return true
}
