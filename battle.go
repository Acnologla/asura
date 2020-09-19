package main

import (
	_ "asura/src/commands" // Initialize all commands and put them into an array
	"asura/src/utils"
	"fmt"
)

func measure() {
    wins := 0

    for j := 0; j < 100000; j++ {
        first := utils.Galo{
            Name: "Papel",
            Xp: utils.CalcXP(5),
            Type: 1,
            Skills: []int{},
            Equipped: []int{},
        }
    
        sec := utils.Galo{
            Name: "Pedra",
            Xp: utils.CalcXP(5),
            Type: 2,
            Skills: []int{},
            Equipped: []int{},
        }
    
        utils.ChooseSkills(&first)
        utils.ChooseSkills(&sec)
    
        battle := utils.CreateBattle(&first,&sec)
        
        for battle.Fighters[0].Life != 0 && battle.Fighters[1].Life != 0  {
            utils.PlayBattle(&battle)
            //fmt.Println(effects)
            //fmt.Println(battle.Fighters[0].Life, " vs ", battle.Fighters[1].Life)
        }

        if battle.Fighters[1].Life == 0 {
            wins++
        } 
    }

	winsPercent := (wins/ 1000);
    fmt.Printf("%d%% vs %d%%",winsPercent, 100 - winsPercent)
}

func one(){
	first := utils.Galo{
		Name: "Papel",
		Xp: utils.CalcXP(5),
		Type: 1,
		Skills: []int{18},
		Equipped: []int{18},
	}

	sec := utils.Galo{
		Name: "Pedra",
		Xp: utils.CalcXP(5),
		Type: 2,
		Skills: []int{},
		Equipped: []int{},
	}

	utils.ChooseSkills(&first)
	utils.ChooseSkills(&sec)

	battle := utils.CreateBattle(&first,&sec)
	
	for battle.Fighters[0].Life != 0 && battle.Fighters[1].Life != 0  {
		if battle.Turn {
			fmt.Println("Turno do de Pedra2")
		} else {
			fmt.Println("Turno do de Papel1")
		}
		effects := utils.PlayBattle(&battle)
		fmt.Println(effects)
		fmt.Println(battle.Fighters[0].Life, " vs ", battle.Fighters[1].Life)
	}
}


func main(){
	measure()
}