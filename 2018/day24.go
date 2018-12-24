package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

var inputFile = flag.String("inputFile", "inputs/day24.input", "Relative file path to use as input.")
var verbose = flag.Bool("verbose", false, "Whether to print verbose debug output.")

type Group struct {
	Side       string
	Count      int
	PerUnitHP  int
	Weaknesses map[string]bool
	Immunities map[string]bool
	Damage     map[string]int
	Init       int
}

func (u *Group) ComputeDamage(t *Group) int {
	totalDmg := 0
	for kind, dmg := range u.Damage {
		if t.Immunities[kind] {
			continue
		}
		totalDmg += dmg * u.Count
		if t.Weaknesses[kind] {
			totalDmg += dmg * u.Count
		}
	}
	return totalDmg
}

type ByTurnOrder []*Group

func (b ByTurnOrder) Len() int      { return len(b) }
func (b ByTurnOrder) Swap(i, j int) { b[i], b[j] = b[j], b[i] }
func (b ByTurnOrder) Less(i, j int) bool {
	iDmg := 0
	for _, v := range b[i].Damage {
		iDmg += v
	}
	jDmg := 0
	for _, v := range b[j].Damage {
		jDmg += v
	}
	if iDmg*b[i].Count != jDmg*b[j].Count {
		return iDmg*b[i].Count > jDmg*b[j].Count
	}
	return b[i].Init > b[j].Init
}

func main() {
	flag.Parse()
	for b := 0; ; b++ {
		winner, units := RunCombat(b)
		if b == 0 {
			fmt.Printf("Number of units remaining with no boost: %d\n", units)
		}
		if winner {
			fmt.Printf("Immune system wins with %d boost: %d units alive\n", b, units)
			break
		}
	}
}

func RunCombat(b int) (bool, int) {
	f, err := os.Open(*inputFile)
	if err != nil {
		return true, -1
	}
	defer f.Close()

	combatants := make([]*Group, 0)
	combatantsByInit := make(map[int]*Group)

	r := bufio.NewReader(f)
	var side string
	var maxInit int
	for {
		l, err := r.ReadString('\n')
		if err != nil {
			break
		}
		if len(l) <= 1 {
			continue
		}
		l = l[:len(l)-1]
		if l[len(l)-1] == ':' {
			// This is a new side definition.
			side = l[:len(l)-1]
			continue
		}
		parsed := strings.SplitN(l, " ", 8)
		count, _ := strconv.Atoi(parsed[0])
		hp, _ := strconv.Atoi(parsed[4])
		parsed = strings.Split(parsed[7], "with an attack that does ")
		immunitiesString := parsed[0]
		immune := make(map[string]bool)
		weak := make(map[string]bool)
		if len(immunitiesString) > 1 {
			parts := strings.Split(immunitiesString[1:len(immunitiesString)-2], "; ")
			for _, v := range parts {
				spec := strings.SplitN(v, " ", 3)
				var kind map[string]bool
				switch spec[0] {
				case "immune":
					kind = immune
				case "weak":
					kind = weak
				}
				for _, t := range strings.Split(spec[2], ", ") {
					kind[t] = true
				}
			}
		}

		damageAndInitString := strings.Split(parsed[1], " damage at initiative ")
		damage := strings.Split(damageAndInitString[0], " ")
		dmgType := damage[1]
		dmg, _ := strconv.Atoi(damage[0])
		if side == "Immune System" {
			dmg += b
		}
		damageMap := map[string]int{dmgType: dmg}
		init, _ := strconv.Atoi(damageAndInitString[1])
		if init > maxInit {
			maxInit = init
		}
		group := Group{side, count, hp, weak, immune, damageMap, init}
		combatants = append(combatants, &group)
		combatantsByInit[init] = &group
	}

	for round := 1; ; round++ {
		if *verbose {
			fmt.Printf("\nRound %d\n", round)
		}
		// Pick targets.
		turnOrder := make([]*Group, len(combatants))

		// Indexed by targeting group.
		targets := make(map[*Group]*Group)

		// Indexed by targeted group.
		targeted := make(map[*Group]bool)

		copy(turnOrder, combatants)
		sort.Sort(ByTurnOrder(turnOrder))

		for _, u := range turnOrder {
			if u.Count <= 0 {
				// We're dead, skip.
				continue
			}
			highestDamage := 0
			potentialTargets := make([]*Group, 0)
			for _, t := range combatants {
				if u.Side == t.Side || t.Count <= 0 || targeted[t] {
					continue
				}
				// We know this is an enemy unit. Compute the damage we would do.
				totalDmg := u.ComputeDamage(t)
				if totalDmg > 0 {
					if totalDmg > highestDamage {
						potentialTargets = []*Group{t}
						highestDamage = totalDmg
					} else if totalDmg < highestDamage {
						continue
					}
					potentialTargets = append(potentialTargets, t)
				}
			}
			if len(potentialTargets) == 0 {
				// Couldn't damage anything
				continue
			}
			sort.Sort(ByTurnOrder(potentialTargets))
			target := potentialTargets[0]
			targeted[target] = true
			targets[u] = target
		}

		deaths := 0
		for i := maxInit; i >= 0; i-- {
			u := combatantsByInit[i]
			if u == nil || u.Count <= 0 {
				continue
			}
			t := targets[u]
			if t == nil {
				continue
			}
			killed := u.ComputeDamage(t) / t.PerUnitHP
			if killed > t.Count {
				killed = t.Count
			}
			t.Count -= killed
			deaths += killed
			if *verbose {
				fmt.Printf("Unit at init %d killed %d units of init %d.\n", u.Init, killed, t.Init)
			}
		}
		if deaths == 0 {
			if *verbose {
				fmt.Println("Stalemated with boost %d.\n", b)
			}
			return false, -1
		}

		// See if only one side is left alive.
		alive := make(map[string]int)
		for _, c := range combatants {
			if c.Count > 0 {
				alive[c.Side] += c.Count
			}
		}

		if len(alive) <= 1 {
			for k, v := range alive {
				if *verbose {
					fmt.Printf("Winning side is %s with %d alive.\n", k, v)
				}
				return k == "Immune System", v
			}
		}
	}
}
