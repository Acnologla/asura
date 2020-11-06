package rinha

import (
	"math/rand"
)

const (
	Damaged      = "Damaged"
	Effected     = "Effected"
	NotEffective = "NotEffective"
	SideEffected = "SideEffected"
	Killed       = "Killed"
)

type Result struct {
	Effect EffectType
	Skill  *Skill
	Damage int
	Self   bool
}

type Round struct {
	Results    []*Result
	Attacker   *Fighter
	Target     *Fighter
	Skill      *Skill
	SkillId    int
	Integrity  float32
	ShieldUser *Fighter
	Stun       bool
}

func (round *Round) applySkillDamage(firstTurn bool) int {
	round.SkillId = round.Attacker.Equipped[rand.Intn(len(round.Attacker.Equipped))]
	round.Skill = Skills[round.Attacker.Galo.Type-1][round.SkillId]

	attack_damage := 0

	if round.Skill.Damage[1] != 0 {
		attack_damage = int(float32(Between(round.Skill.Damage)))
	}

	not_effective_damage := 0

	difference := CalcLevel(round.Attacker.Galo.Xp) - CalcLevel(round.Target.Galo.Xp)

	if firstTurn {
		attack_damage = attack_damage/2 - attack_damage/15
	}

	if IsIntInList(round.Target.Galo.Type, Classes[round.Attacker.Galo.Type].Disadvantages) && rand.Float64() <= 0.72 && difference < 4 {
		not_effective_damage = int(float64(attack_damage) * 0.4)
	}

	if round.Integrity != 0 && round.Integrity != 1 && round.ShieldUser == round.Target {
		not_effective_damage += int(float32(attack_damage) * round.Integrity)
	}

	if not_effective_damage != 0 {
		round.Results = append([]*Result{
			&Result{Effect: NotEffective, Damage: -not_effective_damage, Skill: round.Skill, Self: false},
		}, round.Results...)
	}

	round.Results = append([]*Result{
		&Result{Effect: Damaged, Damage: attack_damage - not_effective_damage, Skill: round.Skill, Self: false},
	}, round.Results...)

	return attack_damage - not_effective_damage
}

func (round *Round) applyEffect(self bool, to_append bool) {
	receiver := round.Target

	if self {
		receiver = round.Attacker
	}

	receiver.Effect[0]--

	effect := Effects[receiver.Effect[1]]

	effect_damage := Between(effect.Range)

	switch effect.Type {
	case 1:
		{
			if effect_damage >= receiver.Life {
				effect_damage = receiver.Life - 1
			}
			receiver.Life -= effect_damage
		}
	case 2:
		{
			if effect_damage >= receiver.MaxLife {
				receiver.Life = receiver.MaxLife
			} else {
				receiver.Life += effect_damage
			}
		}
	case 3:
		{ // Stun
			round.Stun = true
		}
	case 4:
		{ // Defesa
			round.ShieldUser = receiver
			round.Integrity = float32(effect_damage)/100 + 0.01
		}
	}

	if to_append {
		round.Results = append(round.Results, &Result{Effect: Effected, Damage: effect_damage, Self: self, Skill: Skills[receiver.Effect[2]-1][receiver.Effect[3]]})
	}
}

func (round *Round) applyEffects() {
	if rand.Float64() <= round.Skill.Effect[0] {
		effect := Effects[int(round.Skill.Effect[1])]
		if effect.Self {
			round.Attacker.Effect = [4]int{effect.Turns, int(round.Skill.Effect[1]), round.Attacker.Galo.Type, round.SkillId}
			round.applyEffect(true, true)
		} else {
			round.Target.Effect = [4]int{effect.Turns, int(round.Skill.Effect[1]), round.Attacker.Galo.Type, round.SkillId}
			round.applyEffect(false, true)

		}
	}
}

func (battle *Battle) Play() []*Result {

	if battle.Stun {
		battle.Turn = !battle.Turn
		battle.Stun = false
	}

	round := Round{
		Results:   []*Result{},
		Integrity: 1,
	}

	if battle.Turn {
		round.Target = battle.Fighters[0]
		round.Attacker = battle.Fighters[1]
	} else {
		round.Target = battle.Fighters[1]
		round.Attacker = battle.Fighters[0]
	}

	if round.Target.Effect[0] != 0 {
		round.applyEffect(false, true)
	}

	if round.Attacker.Effect[0] != 0 {
		round.applyEffect(true, true)
	}

	if round.Integrity != 0 {
		main_damage := round.applySkillDamage(battle.FirstRound)
		round.applyEffects()
		battle.Stun = round.Stun
		round.Target.Life -= main_damage
	}

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
