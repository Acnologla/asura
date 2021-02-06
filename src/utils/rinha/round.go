package rinha

import (
	"asura/src/utils"
	"math"
	"math/rand"
)

const (
	Damaged      = "Damaged"
	Effected     = "Effected"
	NotEffective = "NotEffective"
	Killed       = "Killed"
)

type Result struct {
	Effect    EffectType
	EffectID  int
	Skill     *Skill
	Damage    int
	Self      bool
	Reflected bool
}

type Round struct {
	Results       []*Result
	Attacker      *Fighter
	Target        *Fighter
	Skill         *Skill
	SkillId       int
	Integrity     float32
	ShieldUser    *Fighter
	FragilityUser *Fighter
	Fragility     float32
	Stun          bool
	Reflex        bool
}

func (round *Round) applySkillDamage(firstTurn bool) int {
	reflected := false
	if round.SkillId == 0 {
		round.SkillId = round.Attacker.Equipped[utils.RandInt(len(round.Attacker.Equipped))]
		round.Skill = Skills[round.Attacker.Galo.Type-1][round.SkillId]
	} else {
		reflected = true
		round.Skill = Skills[round.Target.Galo.Type-1][round.SkillId]
	}

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
	if round.FragilityUser != nil && round.Fragility != 0 {
		if round.FragilityUser == round.Target {
			attack_damage = int(float32(attack_damage) * (1 + round.Fragility))
		}
	}
	if round.Integrity != 0 && round.Integrity != 1 && round.ShieldUser == round.Target {
		not_effective_damage += int(float32(attack_damage) * round.Integrity)
	}

	real_damage := attack_damage - not_effective_damage

	// 1 is the ID of ITEM EFFECT that increases damage
	if round.Attacker.ItemEffect == 1 {
		real_damage = int(math.Round(float64(real_damage) * round.Attacker.ItemPayload))
	}

	// 5 is the ID of ITEM EFFECT that increases defense of the player ( reduces the damage )
	if round.Target.ItemEffect == 5 {
		real_damage = int(math.Round(float64(real_damage) * round.Target.ItemPayload))
	}

	if not_effective_damage != 0 {
		real_not_effective := int(math.Round(float64(not_effective_damage) * round.Target.ItemPayload))
		round.Results = append([]*Result{
			&Result{Effect: NotEffective, Damage: -real_not_effective, Skill: round.Skill, Self: false, EffectID: 0},
		}, round.Results...)
	}

	round.Results = append([]*Result{
		&Result{Reflected: reflected, Effect: Damaged, Damage: real_damage, Skill: round.Skill, Self: false, EffectID: 0},
	}, round.Results...)

	return real_damage
}

func (round *Round) applyEffectDamage(receiver *Fighter, effect *Effect) int {

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
				if round.FragilityUser == receiver {
					effect_damage = int(float32(effect_damage) * round.Fragility)
				}
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
	case 5:
		{
			round.Reflex = true
		}
	case 6:
		{
			round.FragilityUser = receiver
			round.Fragility = float32(effect_damage) / 100
		}
	}
	return effect_damage
}

func (round *Round) applyEffect(id int, self bool, to_append bool) {
	effect := Effects[id]

	receiver := round.Target

	if self {
		receiver = round.Attacker
	}

	receiver.Effect[0]--

	effect_damage := round.applyEffectDamage(receiver, effect)

	if to_append {
		round.Results = append(round.Results,
			&Result{
				EffectID: id,
				Effect:   Effected,
				Damage:   effect_damage,
				Self:     self,
				Skill:    nil,
			})
	}
}

func (round *Round) applyEffects() {
	increase := 0.0
	// 3 id is the EFFECT ID for a increase o probability of effects
	if round.Attacker.ItemEffect == 3 {
		increase = round.Attacker.ItemPayload
	}
	if round.Skill.Effect[0] != 0 && rand.Float64() <= round.Skill.Effect[0]+increase {
		effect := Effects[int(round.Skill.Effect[1])]
		effect_phy := [4]int{effect.Turns, int(round.Skill.Effect[1]), round.Attacker.Galo.Type, round.SkillId}
		if effect.Self {
			round.Attacker.Effect = effect_phy
		} else {
			round.Target.Effect = effect_phy
		}
		round.applyEffect(int(round.Skill.Effect[1]), effect.Self, true)
	} else {
		if round.Attacker.ItemEffect == 2 && rand.Float64() >= 0.7 {
			id := int(math.Round(round.Attacker.ItemPayload))
			effect := Effects[id]
			effect_phy := [4]int{effect.Turns, id, round.Attacker.Galo.Type, round.SkillId}
			if effect.Self {
				round.Attacker.Effect = effect_phy
			} else {
				round.Target.Effect = effect_phy
			}
			round.applyEffect(id, effect.Self, true)
		}
	}
}

func (battle *Battle) Play() []*Result {
	if battle.Stun {
		battle.Turn = !battle.Turn
		battle.Stun = false
	} else if battle.ReflexType == 3 {
		battle.ReflexType = 0
		battle.Turn = !battle.Turn
	}

	round := Round{
		Results:   []*Result{},
		Integrity: 1,
	}

	if battle.ReflexType == 2 {
		round.SkillId = battle.ReflexSkill
		battle.ReflexType = 3
	}

	if battle.Turn {
		round.Target = battle.Fighters[0]
		round.Attacker = battle.Fighters[1]
	} else {
		round.Target = battle.Fighters[1]
		round.Attacker = battle.Fighters[0]
	}

	if round.Target.Effect[0] > 0 {
		round.applyEffect(round.Target.Effect[1], false, true)
	}

	if round.Attacker.Effect[0] > 0 {
		round.applyEffect(round.Attacker.Effect[1], true, true)
	}

	if round.Integrity != 0 {
		main_damage := round.applySkillDamage(battle.FirstRound)
		if battle.ReflexType == 1 {
			battle.ReflexSkill = round.SkillId
			battle.ReflexType = 2
			battle.Turn = !battle.Turn
			return []*Result{}
		}
		round.applyEffects()
		battle.Stun = round.Stun
		if round.Reflex {
			battle.ReflexType = 1
		}
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
