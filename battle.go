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
	class1 := 3
	class2 := 4
	for i := 5; i < 14; i++ {	
		measure(class1,class2,5,i)
	}
}