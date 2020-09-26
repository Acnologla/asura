package rinha

import (
	"math/rand"
)

const (
    Damaged = "Damaged"
	Effected = "Effected"
	NotEffective = "NotEffective"
	SideEffected = "SideEffected"
    Killed = "Killed"
)

type Result struct {
	Effect EffectType
	Skill Skill
	Damage int
	Self bool
}

type Round struct {
	Results []Result
	Attacker *Fighter
	Target *Fighter
	Skill *Skill
	SkillId int
}

func (round *Round) applyEffect(self bool){
	receiver := round.Target

	if self {
		receiver = round.Attacker
	}

	receiver.Effect[0]-- 

	effect_damage := Between(AttackEffects[receiver.Effect[1]].Range)
	if effect_damage >= receiver.Life {
		effect_damage = receiver.Life - 1
	}
	receiver.Life -= effect_damage
	round.Results = append(round.Results, Result{Effect: Effected, Damage: effect_damage, Self: self, Skill: Skills[receiver.Effect[2]]})
}


func (round *Round) applySkillDamage(firstTurn bool) int {
	round.SkillId = round.Attacker.Equipped[rand.Intn(len(round.Attacker.Equipped))]
	round.Skill = &Skills[round.SkillId]

	attack_damage := Between(round.Skill.Damage)
	not_effective_damage := 0

	difference := CalcLevel(round.Attacker.Galo.Xp) - CalcLevel(round.Target.Galo.Xp) 

	if firstTurn {
		attack_damage = attack_damage/2
	}
	
	if IsIntInList(round.Target.Galo.Type, Classes[round.Skill.Type].Disadvantages) && rand.Float64() <= 0.72 && difference < 4 {
		not_effective_damage = int(float64(attack_damage)*0.4)
	}

	round.Results = append(round.Results, Result{Effect: Damaged, Damage: attack_damage, Skill: *round.Skill, Self: false})

	if not_effective_damage != 0 {
		round.Results = append(round.Results, Result{Effect: NotEffective, Damage: -not_effective_damage, Skill: *round.Skill, Self: false})
	}

	return attack_damage - not_effective_damage 
}

func (round *Round) applyEffects() {
	if rand.Float64() <= round.Skill.Effect[0] {
		effect := AttackEffects[int(round.Skill.Effect[1])]
		round.Target.Effect = [3]int{effect.Turns, int(round.Skill.Effect[1]), round.SkillId}
		round.applyEffect(false)
	} else {
		if round.Target.Effect[0] != 0 {
			round.applyEffect(false)
		}
	}
	if round.Attacker.Effect[0] != 0 {
		round.applyEffect(true)
	}
}

func (battle *Battle) Play() []Result {

	round := Round {
		Results: []Result{},
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
	return round.Results
}