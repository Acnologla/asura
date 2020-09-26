package main

import (
	"asura/src/utils/rinha"
	"fmt"
)

func measure(firstClass int, secClass int, firstLvl int, secLvl int, times int, log bool) {
	wins := 0
	wins1 := 0
	empate := 0
    for j := 0; j < times; j++ {
        first := rinha.Galo{
            Name: "Papel",
            Xp: rinha.CalcXP(firstLvl),
            Type: firstClass,
            Equipped: []int{},
        }
    
        sec := rinha.Galo{
            Name: "Pedra",
            Xp: rinha.CalcXP(secLvl),
            Type: secClass,
            Equipped: []int{},
        }
    
		battle := rinha.CreateBattle(&first,&sec)
        
        for battle.Fighters[0].Life != 0 && battle.Fighters[1].Life != 0  {
            if (log){
                effects := battle.Play()
                fmt.Println(effects)
                fmt.Printf("%d vs %d\n",battle.Fighters[0].Life,battle.Fighters[1].Life)
            } else {
                battle.Play()
            }
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
	fmt.Printf("1: Classe: %s, Level: %d Venceu: %d%%\n2: Classe: %s, Level: %d Venceu: %d%%\n\n", rinha.Classes[firstClass].Name, firstLvl, (wins/times)*100, rinha.Classes[secClass].Name, secLvl, wins1/times*100)
}

func main(){
    
	measure(2,2,5,5,1,true)
}