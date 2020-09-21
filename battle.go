package main

import (
	"asura/src/utils"
	"fmt"
)

func measure(firstClass int, secClass int, firstLvl int, secLvl int) {
	wins := 0
	wins1 := 0
	empate := 0
    for j := 0; j < 100000; j++ {
        first := utils.Galo{
            Name: "Papel",
            Xp: utils.CalcXP(firstLvl),
            Type: firstClass,
            Equipped: []int{},
        }
    
        sec := utils.Galo{
            Name: "Pedra",
            Xp: utils.CalcXP(secLvl),
            Type: secClass,
            Equipped: []int{},
        }
    
		battle := utils.CreateBattle(&first,&sec)
        
        for battle.Fighters[0].Life != 0 && battle.Fighters[1].Life != 0  {
            battle.Play()
        }

		if battle.Fighters[1].Life == 0 && battle.Fighters[0].Life == 0 {
			empate++;
		} else if battle.Fighters[1].Life == 0  {
            wins++
		} else if battle.Fighters[0].Life == 0  {
            wins1++
		} 
		
    }

	fmt.Println("\n------------- Analitical tests for ----------- ")
	fmt.Printf("1: Classe: %s, Level: %d Venceu: %d%%\n2: Classe: %s, Level: %d Venceu: %d%%\n\n", utils.Classes[firstClass].Name, firstLvl, wins/1000, utils.Classes[secClass].Name, secLvl, wins1/1000)
}

func main(){
	measure(2,2,6,6)
	measure(3,3,6,6)
	measure(4,4,6,6)
	measure(2,3,6,6)
	measure(3,2,6,6)
	measure(3,4,6,6)
	measure(4,3,6,6)
	measure(2,4,6,6)
	measure(2,2,6,7)
	measure(3,3,6,7)
	measure(4,4,6,7)
	measure(2,2,6,14)
	measure(3,3,6,14)
	measure(4,4,6,14)
}