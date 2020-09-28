package main

import (
	"asura/src/utils/rinha"
    "fmt"
    "math/rand"
    "time"
)

func one(firstClass int, secClass int, firstLvl int, secLvl int) {
    fmt.Printf("Only one --------\n1: Classe: %s, Level: %d\n2: Classe: %s, Level: %d\n\n", rinha.Classes[firstClass].Name, firstLvl, rinha.Classes[secClass].Name, secLvl)
    
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
        turno := 0
        if battle.Turn {
            turno = 1
        }
        effects := battle.Play()
        fmt.Printf("Turno de %d\n%d vs %d\n",turno, battle.Fighters[0].Life,battle.Fighters[1].Life)
        for i := 0; i < len(effects); i++ {
            if effects[i].Effect == rinha.Effected {
                fmt.Println(effects[i].Effect, rinha.GetEffectFromSkill(effects[i].Skill))
            } else {
                fmt.Println(effects[i].Effect, *effects[i].Skill)
            }
        }
        fmt.Println("---")
    }

    if battle.Fighters[1].Life == 0 && battle.Fighters[0].Life == 0 {
        fmt.Println("Empate")
    } else if battle.Fighters[1].Life == 0  {
        fmt.Println("Primeiro ganhou")
    } else if battle.Fighters[0].Life == 0  {
        fmt.Println("Seugndo ganhou")
    } 
}


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
	fmt.Printf("1: Classe: %s, Level: %d Venceu: %d%%\n2: Classe: %s, Level: %d Venceu: %d%%\n\n", rinha.Classes[firstClass].Name, firstLvl, wins/(times/100), rinha.Classes[secClass].Name, secLvl, wins1/(times/100))
}

func main(){
    rand.Seed(time.Now().UTC().UnixNano())
	one(2,4,3,3)
	//measure(2,4,3,3,10000,false)
}