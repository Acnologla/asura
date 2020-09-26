package rinha

type EffectType string

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
		Life: 100 + (CalcLevel(first.Xp) * 3),
		Equipped: []int{},
		Effect: [3]int{},
	}

	secFighter := Fighter{
		Galo: sec,
		Life: 100  + (CalcLevel(sec.Xp) * 3),
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
