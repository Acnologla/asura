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

type EffectType string

const (
    Damaged = "Damaged"
	Effected = "Effected"
	NotEffective = "NotEffective"
	SideEffected = "SideEffected"
    Killed = "Killed"
)


type SideEffect struct {
	Effect EffectType
	Skill Skill
	Damage int
	Self bool
}

type Round struct {
	Effects []SideEffect
	Attacker *Fighter
	Target *Fighter
	Skill *Skill
	SkillId int
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
	skills := GetSkills(*fighter.Galo)

	if len(skills) == 0 {
		skills = append(skills, 0)
	}

	for i := 0; i < len(fighter.Galo.Equipped); i++{
		fighter.Equipped = append(fighter.Equipped, fighter.Galo.Equipped[i])
	}
	need := 5 - len(fighter.Equipped)

	for i := len(skills) - 1; i >= 0 && need != 0; i-- {
		if !IsIntInList(skills[i], fighter.Galo.Equipped) {
			fighter.Equipped = append(fighter.Equipped, skills[i])
			need--
		}
	}
}	

func (battle *Battle) GetReverseTurn() int {
	if !battle.Turn {
		return 1
	} else {
		return 0
	}
} 

func (battle *Battle) GetTurn() int {
	if battle.Turn {
		return 1
	} else {
		return 0
	}
}

func (round *Round) applySkillDamage(firstTurn bool) int {
	round.SkillId = round.Attacker.Equipped[rand.Intn(len(round.Attacker.Equipped))]
	round.Skill = &Skills[round.SkillId]

	attack_damage := Between(round.Skill.Damage)
	not_effective_damage := 0

	if firstTurn {
		attack_damage = attack_damage/2
	}
	
	round.Effects = append(round.Effects, SideEffect{Effect: Damaged, Damage: attack_damage, Skill: *round.Skill, Self: false})

	if IsIntInList(round.Target.Galo.Type, Classes[round.Skill.Type].Disadvantages) && rand.Float64() <= 0.72 {
		not_effective_damage = attack_damage/2
		round.Effects = append(round.Effects, SideEffect{Effect: NotEffective, Damage: -not_effective_damage, Skill: *round.Skill, Self: false})
	}

	return attack_damage - not_effective_damage 
}

func (round *Round) applyEffects() {
	if rand.Float64() <= round.Skill.Effect[0] {
		effect := AttackEffects[int(round.Skill.Effect[1])]
		round.Target.Effect = [3]int{effect.Turns, int(round.Skill.Effect[1]), round.SkillId}
		effect_damage := Between(AttackEffects[round.Target.Effect[1]].Damage)
		round.Target.Life -= effect_damage
		round.Effects = append(round.Effects, SideEffect{Effect: Effected, Damage: effect_damage, Skill: *round.Skill, Self: false})
	} else {
		if round.Target.Effect[0] != 0 {
			round.Target.Effect[0]-- 
			effect_damage := Between(AttackEffects[round.Target.Effect[1]].Damage)
			round.Target.Life -= effect_damage
			round.Effects = append(round.Effects, SideEffect{Effect: SideEffected, Damage: effect_damage, Self: false, Skill: Skills[round.Target.Effect[2]]})
		}
	}
	
	if round.Attacker.Effect[0] != 0 {
		round.Attacker.Effect[0]-- 
		effect_damage := Between(AttackEffects[round.Attacker.Effect[1]].Damage)
		round.Attacker.Life -= effect_damage
		round.Effects = append(round.Effects, SideEffect{Effect: SideEffected, Damage: effect_damage, Self: false, Skill: Skills[round.Attacker.Effect[2]]})
	}
}

func (battle *Battle) Play() []SideEffect {

	round := Round {
		Effects: []SideEffect{},
	}

	if battle.Turn {
        round.Target = &battle.Fighters[0]
        round.Attacker = &battle.Fighters[1]
    } else {
        round.Target = &battle.Fighters[1]
        round.Attacker = &battle.Fighters[0]
	}
	
	main_damage := round.applySkillDamage(battle.FirstRound)
	round.applyEffects()


	round.Target.Life -= main_damage

	if round.Attacker.Life <= 0 {
		round.Attacker.Life = 0
	}

	if round.Target.Life <= 0 {
		round.Target.Life = 0
	}

	battle.FirstRound = false
	battle.Turn = !battle.Turn
	return round.Effects
}