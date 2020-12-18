package rinha

import "math"

type EffectType string

type Fighter struct {
	Galo        *Galo
	Equipped    []int
	Life        int
	ItemEffect  int
	ItemPayload float64
	MaxLife     int
	Effect      [4]int
}

type Battle struct {
	Stopped    bool
	Fighters   [2]*Fighter
	Turn       bool
	FirstRound bool
	Stun       bool
}

func checkItem(galo *Galo) (int, float64) {
	if len(galo.Items) > 0 {
		id := galo.Items[0]
		item := Items[id]
		return item.Effect, item.Payload
	}
	return 0, 0
}

func initFighter(galo *Galo) *Fighter {
	life := 100 + (CalcLevel(galo.Xp) * 3)
	itemEffect, payload := checkItem(galo)

	// 3 is the ID of Item EFFECT that increase life
	if itemEffect == 3 {
		life = int(math.Round(float64(life) * payload))
	}

	return &Fighter{
		Galo:        galo,
		Life:        life,
		MaxLife:     life,
		ItemEffect:  itemEffect,
		ItemPayload: payload,
		Equipped:    []int{},
		Effect:      [4]int{},
	}
}

func CreateBattle(first *Galo, sec *Galo) Battle {
	firstFighter := initFighter(first)
	secFighter := initFighter(sec)

	initEquips(firstFighter)
	initEquips(secFighter)

	return Battle{
		Stopped:    false,
		Turn:       false,
		FirstRound: true,
		Fighters: [2]*Fighter{
			firstFighter,
			secFighter,
		},
	}
}

func GetEquipedSkills(galo *Galo) []int {
	skills := GetSkills(*galo)
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
	return equipedSkills
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
