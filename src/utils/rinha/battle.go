package rinha

import (
	"math"

	"github.com/andersfylling/disgord"
)

type EffectType string

type NewSkill struct {
	Skill    int
	Cooldown int
}

type Fighter struct {
	Galo        *Galo
	Equipped    []*NewSkill
	Revived     bool
	Life        int
	ID          disgord.Snowflake
	ItemEffect  int
	ItemPayload float64
	Username    string
	MaxLife     int
	Effect      [4]int
}

type Battle struct {
	Stopped     bool
	Fighters    [2]*Fighter
	Waiting     []*Fighter
	WaitingN    int
	Turn        bool
	FirstRound  bool
	Stun        bool
	ReflexType  int
	Reseted     bool
	ReflexSkill int
}

func CheckItem(galo *Galo) (int, float64) {
	if len(galo.Items) > 0 {
		id := galo.Items[0]
		item := Items[id]
		return item.Effect, item.Payload
	}
	return 0, 0
}

func InitFighter(galo *Galo, noItems bool) *Fighter {
	life := 100 + (CalcLevel(galo.Xp) * 3)
	if HasUpgrade(galo.Upgrades, 1) {
		life += 5
		if HasUpgrade(galo.Upgrades, 1, 1) {
			life += 5
			if HasUpgrade(galo.Upgrades, 1, 1, 1) {
				life += 10
			}
		}

	}
	var itemEffect int
	var payload float64
	if !noItems {
		itemEffect, payload = CheckItem(galo)
	}
	// 4 is the ID of Item EFFECT that increase life
	if itemEffect == 4 {
		life = int(math.Round(float64(life) * payload))
	}

	return &Fighter{
		Galo:        galo,
		Life:        life,
		MaxLife:     life,
		ItemEffect:  itemEffect,
		ItemPayload: payload,
		Equipped:    []*NewSkill{},
		Effect:      [4]int{},
	}
}

func CreateBattle(first Galo, sec Galo, noItems bool, firstID, secondID disgord.Snowflake, waiting []Galo, usernames []string) Battle {
	firstFighter := InitFighter(&first, noItems)
	secFighter := InitFighter(&sec, noItems)
	if HasUpgrade(firstFighter.Galo.Upgrades, 2, 1) {
		if HasUpgrade(firstFighter.Galo.Upgrades, 2, 1, 1) {
			secFighter.Life -= 10
			secFighter.MaxLife -= 10
			if HasUpgrade(firstFighter.Galo.Upgrades, 2, 1, 1, 0) {
				secFighter.Life -= 15
				secFighter.MaxLife -= 15
			}
		}
		secFighter.Life -= 5
		secFighter.MaxLife -= 5
	}
	if HasUpgrade(secFighter.Galo.Upgrades, 2, 1) {
		if HasUpgrade(secFighter.Galo.Upgrades, 2, 1, 1) {
			firstFighter.Life -= 10
			firstFighter.MaxLife -= 10
			if HasUpgrade(secFighter.Galo.Upgrades, 2, 1, 1, 0) {
				firstFighter.Life -= 15
				firstFighter.MaxLife -= 15
			}
		}
		firstFighter.Life -= 5
		firstFighter.MaxLife -= 5
	}
	initEquips(firstFighter)
	initEquips(secFighter)
	firstFighter.ID = firstID
	secFighter.ID = secondID
	waitingBattle := []*Fighter{}
	if len(waiting) > 0 {
		for i, galo := range waiting {
			galoFighter := InitFighter(&galo, noItems)
			initEquips(galoFighter)
			if i == 0 {
				galoFighter = firstFighter
			}
			galoFighter.Username = GetName(usernames[i], galo)

			waitingBattle = append(waitingBattle, galoFighter)
		}
	}
	return Battle{
		Stopped:    false,
		Turn:       false,
		FirstRound: true,
		Fighters: [2]*Fighter{
			firstFighter,
			secFighter,
		},
		Waiting: waitingBattle,
	}
}

func GetEquipedSkills(galo *Galo) []*NewSkill {
	skills := GetSkills(*galo)
	newSkill := []*NewSkill{}
	if len(skills) == 0 {
		skills = append(skills, 0)
	}
	equipedSkills := []int{}
	for i := 0; i < len(galo.Equipped); i++ {
		equipedSkills = append(equipedSkills, galo.Equipped[i])
	}
	need := 5 - len(equipedSkills)
	for i := len(skills) - 1; i >= 0 && need != 0; i-- {
		if !IsIntInList(skills[i], galo.Equipped) {
			equipedSkills = append(equipedSkills, skills[i])
			need--
		}
	}
	for _, skill := range equipedSkills {
		newSkill = append(newSkill, &NewSkill{
			Skill: skill,
		})
	}
	return newSkill
}

func initEquips(fighter *Fighter) {
	fighter.Equipped = GetEquipedSkills(fighter.Galo)
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
