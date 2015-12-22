package main

import (
	"fmt"
)

type Entity struct {
	Hitpoints, Damage, Defense int
}

type Weapon struct {
	Cost, Damage int
}

type Armor struct {
	Cost, Defense int
}

type Ring struct {
	Cost, Damage, Defense int
}

type EquipmentSet struct {
	Wielded Weapon
	Worn Armor
	LeftRing, RightRing Ring
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

func (es EquipmentSet) TotalCost() int {
	return es.Wielded.Cost + es.Worn.Cost + es.LeftRing.Cost + es.RightRing.Cost
}

func (es EquipmentSet) GenerateEntity(startHealth int) Entity {
	return Entity{
		Hitpoints: startHealth,
		Damage: es.Wielded.Damage + es.LeftRing.Damage + es.RightRing.Damage,
		Defense: es.Worn.Defense + es.LeftRing.Defense + es.RightRing.Defense,
	}
}

func Fight(player, boss Entity) bool {
	for {
		if player.AttackOnce(&boss) {
			return true
		}
		if boss.AttackOnce(&player) {
			return false
		}
	}
}

func main() {
	allArmor := []Armor {
		{0, 0},
		{13, 1},
		{31, 2},
		{53, 3},
		{75, 4},
		{102, 5},
	}
	allWeapons := []Weapon {
		{8, 4},
		{10, 5},
		{25, 6},
		{40, 7},
		{74, 8},
	}
	allRings := []Ring {
		{0, 0, 0},
		{0, 0, 0},
		{25, 1, 0},
		{50, 2, 0},
		{100, 3, 0},
		{20, 0, 1},
		{40, 0, 2},
		{80, 0, 3},
	}

	boss := Entity {
		Hitpoints: 109,
		Damage: 8,
		Defense: 2,
	}

	cheapest := 9999999
	expensivest := 0

	for w := range allWeapons {
		for a := range allArmor {
			for lr := range allRings {
				for rr := range allRings {
					if rr >= lr {
						// Don't allow buying the same ring twice; treat rings as symmetric.
						continue
					}
					es := EquipmentSet {
						Wielded: allWeapons[w],
						Worn: allArmor[a],
						LeftRing: allRings[lr],
						RightRing: allRings[rr],
					}
					cost := es.TotalCost()
					player := es.GenerateEntity(100)
					if Fight(player, boss) && cost < cheapest {
						cheapest = cost
					}
					if !Fight(player, boss) && cost > expensivest {
						expensivest = cost
					}
				}
			}
		}
	}
	fmt.Println(cheapest, expensivest)
}
