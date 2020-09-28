package rinha

type EffectType string

type Fighter struct {
	Galo     *Galo
	Equipped []int
	Life     int 
	MaxLife  int
	Effect   [4]int
}

type Battle struct {
	Stopped bool
	Fighters [2]*Fighter
	Turn bool	
	FirstRound bool
	Stun bool
}

func CreateBattle(first *Galo, sec *Galo) Battle {
	firstFighter := &Fighter{
		Galo: first,
		Life: 100 + (CalcLevel(first.Xp) * 3),
		MaxLife: 100  + (CalcLevel(first.Xp) * 3),
		Equipped: []int{},
		Effect: [4]int{},
	}

	secFighter := &Fighter{
		Galo: sec,
		Life: 100  + (CalcLevel(sec.Xp) * 3),
		MaxLife: 100  + (CalcLevel(sec.Xp) * 3),
		Equipped: []int{},
		Effect: [4]int{},
	}
	
	initEquips(firstFighter)
	initEquips(secFighter)
	
	return Battle {
		Stopped: false,
		Turn: false,
		FirstRound: true,
		Fighters: [2]*Fighter{
			firstFighter,		
			secFighter,
		},
	}
}

func GetEquipedSkills(galo *Galo)[]int{
	skills := GetSkills(*galo)
	if len(skills) == 0 {
		skills = append(skills, 0)
	}
	equipedSkills := []int{}
	for i := 0; i < len(galo.Equipped); i++{
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
