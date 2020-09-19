package main

import (
	_ "asura/src/commands" // Initialize all commands and put them into an array
	"asura/src/utils"
	"fmt"
)

func main() {
    statistics := [2]int{0,0}

    for j := 0; j < 100000; j++ {
        first := utils.Galo{
            Name: "João",
            Xp: 37,
            Type: 1,
            Skills: []int{},
            Equipped: []int{},
        }
    
        sec := utils.Galo{
            Name: "João",
            Xp: 37,
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

        if battle.Fighters[0].Life == 0 {
            statistics[1]++
        } else {
            statistics[0]++
        }
    }

    fmt.Printf("%d%% vs %d%%",(statistics[0]/ 1000),(statistics[1]/ 1000))
}
