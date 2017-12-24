package main

import (
	"fmt"
	"math"
)

var parallelize = true
var hard = true
var recursionDepth = 13
var verbose = false

type Entity struct {
	Hitpoints, Damage, Defense, Mana int
	Effects                          []Effect
}

var boss = Entity{
	Hitpoints: 71,
	Damage:    10,
	Defense:   0,
	Mana:      0,
}

var player = Entity{
	Hitpoints: 50,
	Damage:    0,
	Defense:   0,
	Mana:      500,
}

type Effect struct {
	Name      string
	TurnsLeft int
	Mana      int
	Armor     int
	Damage    int
}

type Spell struct {
	Cost   int
	Self   Effect
	Target Effect
}

type State struct {
	Player, Boss Entity
	ManaSpent    int
}

func (ent *Entity) ApplyEffects() {
	newEffects := make([]Effect, 0, len(ent.Effects))
	for i := range ent.Effects {
		eff := ent.Effects[i]
		ent.Mana += eff.Mana
		ent.Hitpoints += eff.Damage
		eff.TurnsLeft--
		if eff.TurnsLeft != 0 {
			newEffects = append(newEffects, eff)
		} else {
			ent.Defense -= eff.Armor
		}
	}
	ent.Effects = newEffects
}
func (s *Spell) CastSpell(player, target *Entity) bool {
	if s.Cost > player.Mana {
		return false
	}
	for i := range player.Effects {
		if player.Effects[i].Name == s.Self.Name {
			return false
		}
	}
	for i := range target.Effects {
		if target.Effects[i].Name == s.Target.Name {
			return false
		}
	}

	player.Mana -= s.Cost

	// Armor instantaneously applies, and will be removed by the timer.
	player.Defense += s.Self.Armor
	if s.Self.TurnsLeft > 0 {
		player.Effects = append(player.Effects, s.Self)
	}
	if s.Target.TurnsLeft > 0 {
		target.Effects = append(target.Effects, s.Target)
	}
	return true
}

func (e *Entity) AttackOnce(target *Entity) bool {
	hit := e.Damage - target.Defense
	if hit < 1 {
		hit = 1
	}
	target.Hitpoints -= hit
	if target.Hitpoints <= 0 {
		return true
	}
	return false
}

func FightIncremental(state State, spell Spell) (bool, bool, bool, State) {
	if hard {
		state.Player.Hitpoints--
		if state.Player.Hitpoints <= 0 {
			return false, false, true, state
		}
	}

	state.Player.ApplyEffects()
	state.Boss.ApplyEffects()

	if !spell.CastSpell(&state.Player, &state.Boss) {
		return false, false, false, state
	}
	state.ManaSpent += spell.Cost

	state.Player.ApplyEffects()
	state.Boss.ApplyEffects()

	if state.Boss.Hitpoints <= 0 {
		// Win!
		return true, false, true, state
	}
	if state.Boss.AttackOnce(&state.Player) {
		// Loss!
		return false, false, true, state
	}

	return false, true, true, state
}

var allSpells = []Spell{
	{53, Effect{}, Effect{"MM", 1, 0, 0, -4}},
	{73, Effect{"Drain", 1, 0, 0, 2}, Effect{"Drain", 1, 0, 0, -2}},
	{113, Effect{"Shield", 6, 0, 7, 0}, Effect{}},
	{173, Effect{}, Effect{"Poison", 6, 0, 0, -3}},
	{229, Effect{"Recharge", 5, 101, 0, 0}, Effect{}},
}

func main() {
	leastManaSpent := TestSpells([]Spell{}, State{player, boss, 0}, recursionDepth)
	fmt.Println(leastManaSpent)
}

func TestSpells(prev []Spell, os State, depth int) int {
	values := make(chan int, len(allSpells))

	for s := range allSpells {
		if allSpells[s].Cost > 100 {
			lookBack := 2
			i := len(prev) - lookBack
			if i < 0 {
				i = 0
			}
			for ; i < len(prev); i++ {
				if allSpells[s].Cost == prev[i].Cost {
					// Save fan-out.
					continue
				}
			}
		}

		compute := func(previous []Spell, oldState State, maxDepth, spell int) {
			spellList := make([]Spell, len(previous)+1)
			copy(spellList[:len(previous)], previous)
			spellList[len(previous)] = allSpells[spell]

			// win, proceed, state := Fight(player, boss, spellList)
			win, proceed, valid, state := FightIncremental(oldState, allSpells[spell])

			finalState := state
			finalState.Player.ApplyEffects()
			finalState.Boss.ApplyEffects()

			// Part (b)
			if hard {
				finalState.Player.Hitpoints--
				if finalState.Boss.Hitpoints <= 0 && finalState.Player.Hitpoints > 0 && valid {
					win = true
				}
			} else {
				if finalState.Boss.Hitpoints <= 0 && valid {
					win = true
				}
			}

			if win {
				if verbose {
					fmt.Printf("%d mana spent: ", finalState.ManaSpent)
					for i := range spellList {
						fmt.Printf("%d, ", spellList[i].Cost)
					}
					fmt.Printf("\n")
				}
				values <- finalState.ManaSpent
				// Definitionally our children will spend more mana.
				return
			}

			if maxDepth >= 1 && proceed {
				values <- TestSpells(spellList, state, maxDepth-1)
			} else {
				values <- math.MaxInt32
			}
		}
		if depth == recursionDepth && parallelize {
			go compute(prev, os, depth, s)
		} else {
			compute(prev, os, depth, s)
		}
	}

	leastManaSpent := math.MaxInt32
	for i := 0; i < len(allSpells); i++ {
		v := <-values
		if v < leastManaSpent {
			leastManaSpent = v
		}
	}
	return leastManaSpent
}
