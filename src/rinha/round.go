package rinha

import (
	"asura/src/entities"
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
	Reset     bool
	ExtraMsg  string
	Skill     *Skill
	Damage    int
	Self      bool
	Reflected bool
	Dodge     bool
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

func (round *Round) applySkillDamage(firstTurn bool, skill int) int {
	reflected := false
	var galo *entities.Rooster
	var user *entities.User
	if round.SkillId == 0 {
		if skill != -1 {
			round.SkillId = round.Attacker.Equipped[skill].Skill
		} else {
			round.SkillId = round.Attacker.Equipped[utils.RandInt(len(round.Attacker.Equipped))].Skill
		}
		round.Skill = Skills[round.Attacker.Galo.Type-1][round.SkillId]
		galo = round.Attacker.Galo
		user = round.Attacker.User
		round.Attacker.Life += int(0.3 * float64(user.Attributes[4]))
		if HasUpgrade(user.Upgrades, 1, 1, 0) {
			if HasUpgrade(user.Upgrades, 1, 1, 0, 1) {
				round.Target.Life -= 40
			}
			if HasUpgrade(user.Upgrades, 1, 1, 0, 0) {
				round.Attacker.Life += 25
			}
			round.Attacker.Life += 20
		}
	} else {
		reflected = true
		round.Skill = Skills[round.Target.Galo.Type-1][round.SkillId]
		galo = round.Target.Galo
		user = round.Target.User
	}

	attack_damage := 0

	if round.Skill.Damage[1] != 0 {
		min, max := CalcDamage(round.Skill, galo)
		attack_damage = int(float32(Between([2]int{min, max})))
	}

	if round.Attacker.IsBlack {
		attack_damage = int(float64(attack_damage) * 1.35)
	}

	attack_damage = int(float64(attack_damage) * GetTrialsMultiplier(user))
	attack_damage = int(float32(attack_damage) * (1 + (0.002 * (float32(user.Attributes[1])))))

	if reflected {
		attack_damage = int(float64(attack_damage) * 0.5)
	}
	if HasUpgrade(user.Upgrades, 1, 0) {
		attack_damage += 10
		if HasUpgrade(user.Upgrades, 1, 0, 0) {
			attack_damage += 20
			if HasUpgrade(user.Upgrades, 1, 0, 0, 0) {
				attack_damage += 30
			}
		}
	}

	if HasUpgrade(user.Upgrades, 2, 0, 1) {
		critChance := 9
		if HasUpgrade(user.Upgrades, 2, 0, 1, 0) {
			critChance = 19
		}
		if critChance >= utils.RandInt(101) {
			crit := attack_damage / 5
			if HasUpgrade(user.Upgrades, 2, 0, 1, 1) {
				crit *= 2
			}
			attack_damage += crit
		}

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
			{Effect: NotEffective, Damage: -real_not_effective, Skill: round.Skill, Self: false, EffectID: 0},
		}, round.Results...)
	}

	round.Results = append([]*Result{
		{Reflected: reflected, Effect: Damaged, Damage: real_damage, Skill: round.Skill, Self: false, EffectID: 0},
	}, round.Results...)
	if 0 > real_damage {
		real_damage = 0
	}
	return real_damage
}

func (round *Round) applyEffectDamage(receiver *Fighter, effect *Effect, ataccker *Fighter) int {
	min := effect.Range[0]
	max := effect.Range[1]
	if effect.Type == 1 || effect.Type == 6 {
		min, max = CalcEffectRange(effect, ataccker.Galo)
	} else if effect.Type == 2 {
		min, max = CalcEffectRange(effect, receiver.Galo)
	}
	isNegative := false
	if 0 > min {
		min = min * -1
		max = max * -1
		isNegative = true
	}
	effect_damage := Between([2]int{min, max})
	if isNegative {
		effect_damage = effect_damage * -1
	}

	if receiver.Galo.Type == 51 && effect.Type > 2 {
		return 0
	}

	switch effect.Type {
	case 1:
		{
			if ataccker.IsBlack {
				effect_damage = int(float64(effect_damage) * 1.35)
			}
			if ataccker.ItemEffect == 6 {
				effect_damage = int(ataccker.ItemPayload * float64(effect_damage))
			}

			if HasUpgrade(receiver.User.Upgrades, 2, 0, 0, 1) {
				if 8 >= utils.RandInt(101) {
					effect_damage = 0
				}
			}
			if HasUpgrade(ataccker.User.Upgrades, 1, 0, 1) {
				effect_damage += 30
				if HasUpgrade(ataccker.User.Upgrades, 1, 0, 1, 0) {
					effect_damage += 30
				}
			}

			effect_damage = int(float32(effect_damage) * (1 + (0.002 * (float32(ataccker.User.Attributes[3])))))

			if effect_damage >= receiver.Life {
				effect_damage = receiver.Life - 1
			}

			if ataccker.Galo.Type == 33 {
				lvl := CalcLevel(ataccker.Galo.Xp)
				n := float64(effect_damage) * (1 + (0.13 * float64(lvl/10)))

				ataccker.Life += int(n * (1 + (float64(ataccker.Galo.Resets) / 5)))
			}

			if receiver.Galo.Type == 51 {
				return 0
			}

			// shields
			if receiver.ItemEffect == 12 {
				effect_damage = int(math.Round(float64(effect_damage) * round.Target.ItemPayload))
			}

			receiver.Life -= effect_damage
		}
	case 2:
		{
			if receiver.ItemEffect == 11 {
				effect_damage = int(receiver.ItemPayload * float64(effect_damage))
			}

			if ataccker.ItemEffect == 11 {
				effect_damage = int(float64(effect_damage) * (2 - ataccker.ItemPayload))
			}

			if round.FragilityUser == receiver {
				effect_damage = int(float32(effect_damage) * (1 - round.Fragility))
			}
			if 0 >= receiver.Life+effect_damage {
				effect_damage = receiver.Life - 1
			}
			receiver.Life += effect_damage

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
			if ataccker.Galo.Type == 51 {
				return 0
			}
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
	ataccker := round.Attacker
	if self {
		receiver = round.Attacker
		ataccker = round.Target
	}

	receiver.Effect[0]--

	effect_damage := round.applyEffectDamage(receiver, effect, ataccker)

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
		if round.Attacker.ItemEffect == 2 && rand.Float64() >= 0.6 {
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

func (battle *Battle) Play(skill int) []*Result {
	if battle.Stun {
		battle.Turn = !battle.Turn
		battle.Stun = false
	} else if battle.ReflexType == 3 {
		battle.ReflexType = 0
		battle.Turn = !battle.Turn
	}
	if !battle.Turn && !battle.FirstRound {
		if len(battle.Waiting) > 0 {
			battle.WaitingN++
			if battle.WaitingN >= len(battle.Waiting) {
				battle.WaitingN = 0
			}
			battle.Fighters[0] = battle.Waiting[battle.WaitingN]
		}
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
		main_damage := round.applySkillDamage(battle.FirstRound, skill)
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
		if HasUpgrade(round.Target.User.Upgrades, 2) {
			num := 2
			if HasUpgrade(round.Target.User.Upgrades, 2, 0) {
				num += 2
			}
			if HasUpgrade(round.Target.User.Upgrades, 2, 0, 0) {
				num += 3
			}
			if HasUpgrade(round.Target.User.Upgrades, 2, 0, 0, 0) {
				num += 2
			}
			if num >= utils.RandInt(101) {
				round.Results = append(round.Results, &Result{Dodge: true})
				main_damage = 0
			}
		}
		// 7 reflect item
		if round.Target.ItemEffect == 7 {
			dm := int(float64(main_damage) * round.Target.ItemPayload)
			if 0 >= round.Attacker.Life-dm {
				dm = round.Attacker.Life - 1
			}
			round.Attacker.Life -= dm
		}
		round.Target.Life -= main_damage
	}
	reset := func(f *Fighter) {
		if !battle.Reseted {
			round.Attacker.Life = round.Attacker.MaxLife
			round.Target.Life = round.Target.MaxLife
			f.Life += 100 + (30 * f.Galo.Resets)
			round.Results = append(round.Results, &Result{Reset: true})
			battle.FirstRound = true
			battle.Reseted = true
		}
	}

	battle.FirstRound = false
	battle.Turn = !battle.Turn
	if round.Attacker.Life <= 0 {
		round.Attacker.Life = 0
		if round.Attacker.Galo.Type == 32 {
			reset(round.Attacker)
		}
	}

	if round.Target.Life <= 0 {
		round.Target.Life = 0
		if round.Target.Galo.Type == 32 {
			reset(round.Target)
		}
	}
	return round.Results

}
