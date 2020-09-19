package main

import (
	"asura/src/utils"
	"fmt"
)

func measure() {
    wins := 0
	empate := 0
    for j := 0; j < 100000; j++ {
        first := utils.Galo{
            Name: "Papel",
            Xp: utils.CalcXP(6),
            Type: 1,
            Equipped: []int{},
        }
    
        sec := utils.Galo{
            Name: "Pedra",
            Xp: utils.CalcXP(6),
            Type: 2,
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
		} 
		
		if j % 10000 == 0 {
	
			fmt.Printf("Executed %d\n",j)
		}
    }

	winsPercent := (wins/ 1000);
    fmt.Printf("%d%% | %d%%",winsPercent, empate / 1000)
}

func one() {
	first := utils.Galo{
		Name: "Papel",
		Xp: utils.CalcXP(6),
		Type: 2,
		Equipped: []int{},
	}

	sec := utils.Galo{
		Name: "Pedra",
		Xp: utils.CalcXP(6),
		Type: 2,
		Equipped: []int{},
	}

	battle := utils.CreateBattle(&first,&sec)
	
	
	for battle.Fighters[0].Life != 0 && battle.Fighters[1].Life != 0  {
		effects := battle.Play()
		fmt.Println(effects);
		fmt.Printf("%d x %d\n",battle.Fighters[0].Life,battle.Fighters[1].Life)
	}
}
func main(){
	first := utils.Galo{
		Name: "Papel",
		Xp: utils.CalcXP(6),
		Type: 2,
		Equipped: []int{},
	}

	sec := utils.Galo{
		Name: "Pedra",
		Xp: utils.CalcXP(6),
		Type: 2,
		Equipped: []int{},
	}
	
	battle := utils.CreateBattle(&first,&sec)
	battle.Play()

	fmt.Println(battle.Fighters[0].Life)
}