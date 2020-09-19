package utils

import (
	"math/rand"
)

type Fighter struct {
	Galo    *Galo
	Equipped []int
	Life    int 
	Effect  [3]int
}

type Battle struct {
	Stopped bool
	Fighters [2]Fighter
	Turn bool	
	FirstRound bool
}

func CreateBattle(first *Galo, sec *Galo) Battle {

	firstFighter := Fighter{
		Galo: first,
		Life: 100,
		Equipped: []int{},
		Effect: [3]int{},
	}

	secFighter := Fighter{
		Galo: sec,
		Life: 100,
		Equipped: []int{},
		Effect: [3]int{},
	}
	
	initEquips(&firstFighter)
	initEquips(&secFighter)
	
	return Battle {
		Stopped: false,
		Turn: false,
		FirstRound: true,
		Fighters: [2]Fighter{
			firstFighter,		
			secFighter,
		},
	}
}

func initEquips(fighter *Fighter) {

	if len(fighter.Galo.Skills) == 0 {
		fighter.Galo.Skills = append(fighter.Galo.Skills, 0)
	}

	for i := 0; i < len(fighter.Galo.Equipped); i++{
		fighter.Equipped = append(fighter.Equipped, fighter.Galo.Equipped[i])
	}

	need := 5 - len(fighter.Equipped)

	for i := len(fighter.Galo.Skills) - 1; i >= 0 && need != 0; i-- {
		if !IsIntInList(fighter.Galo.Skills[i], fighter.Galo.Equipped) {
			fighter.Equipped = append(fighter.Equipped, fighter.Galo.Skills[i])
			need--
		}
	}
}	


type EffectType string

const (
    Damaged = "Damaged"
	Effected = "Effected"
	Critic = "Critic"
	SideEffected = "SideEffected"
    Killed = "Killed"
)

type SideEffect struct {
	effect EffectType
	skill Skill
	damage int
}

func randomDamage(damage [2]int) int {
	return rand.Intn(damage[1] - damage[0]) + damage[0]
}

func SaturateSub(one int, two int) int {
	if two >= one {
		return 0
	} else {
		return one - two
	}
}

func PlayBattle(battle *Battle) []SideEffect {

	effects := []SideEffect{}

	var attacker *Fighter
	var target *Fighter

	if battle.Turn {
		target = &battle.Fighters[1]
		attacker = &battle.Fighters[0]
	} else {
		target = &battle.Fighters[0]
		attacker = &battle.Fighters[1]
	}
	
	id := attacker.Equipped[rand.Intn(len(attacker.Equipped))]
	skill := Skills[id]
	
	attack_damage := randomDamage(skill.Damage)
	critic_damage := 0
	effect_damage := 0

	effects = append(effects, SideEffect{effect: Damaged, damage: attack_damage, skill: skill})

	if IsIntInList(target.Galo.Type, Classes[skill.Type].Advantages) && rand.Float64() <= 0.6 {
		critic_damage = attack_damage/2
		effects = append(effects, SideEffect{effect: Critic, damage: critic_damage, skill: skill})
	}

	if rand.Float64() <= skill.Effect[0] {
		effect := AttackEffects[int(skill.Effect[1])]
		target.Effect = [3]int{effect.Turns, int(skill.Effect[1]), id}
		effect_damage = randomDamage(AttackEffects[target.Effect[1]].Damage)
		effects = append(effects, SideEffect{effect: Effected, damage: effect_damage, skill: skill})
	} else {
		if target.Effect[0] != 0 {
			target.Effect[0]-- 
			effect_damage = randomDamage(AttackEffects[target.Effect[1]].Damage)
			effects = append(effects, SideEffect{effect: SideEffected, damage: effect_damage, skill: Skills[target.Effect[2]]})
		}
	}

	damage := attack_damage + critic_damage + effect_damage

	if battle.FirstRound {		
		target.Life = SaturateSub(target.Life, damage/2)
	} else {
		target.Life = SaturateSub(target.Life, damage)
	}

	battle.FirstRound = false
	battle.Turn = !battle.Turn

	return effects
}