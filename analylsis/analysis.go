package main

import (
	"asura/src/entities"
	"asura/src/rinha"
	"fmt"
	"math/rand"
	"time"
)

type stats struct {
	Wins      []int
	TotalWins int
	Level     []int
	Type      int
}

/*
	func one(firstClass int, secClass int, firstLvl int, secLvl int) {
		fmt.Printf("Only one --------\n1: Classe: %s, Level: %d\n2: Classe: %s, Level: %d\n\n", rinha.Classes[firstClass].Name, firstLvl, rinha.Classes[secClass].Name, secLvl)

		first := rinha.Galo{
			Name:     "Papel",
			Xp:       rinha.CalcXP(firstLvl),
			Type:     firstClass,
			Equipped: []int{},
		}

		sec := rinha.Galo{
			Name:     "Pedra",
			Xp:       rinha.CalcXP(secLvl),
			Type:     secClass,
			Equipped: []int{},
		}

		battle := rinha.CreateBattle(first, sec, false, 0, 0)

		for battle.Fighters[0].Life != 0 && battle.Fighters[1].Life != 0 {
			turno := 0
			if battle.Turn {
				turno = 1
			}
			battle.Play(-1)
			fmt.Printf("Turno de %d\n%d vs %d\n", turno, battle.Fighters[0].Life, battle.Fighters[1].Life)
		}

		if battle.Fighters[1].Life == 0 && battle.Fighters[0].Life == 0 {
			fmt.Println("Empate")
		} else if battle.Fighters[1].Life == 0 {
			fmt.Println("Primeiro ganhou")
		} else if battle.Fighters[0].Life == 0 {
			fmt.Println("Seugndo ganhou")
		}
	}
*/
func measure(firstClass int, secClass int, firstLvl int, secLvl int, times int, log bool) int {
	wins := 0
	wins1 := 0
	empate := 0
	for j := 0; j < times; j++ {
		first := &entities.Rooster{
			Name:     "Papel",
			Xp:       rinha.CalcXP(firstLvl),
			Type:     firstClass,
			Equipped: []int{},
			Evolved:  true,
		}

		sec := &entities.Rooster{
			Name:     "Pedra",
			Xp:       rinha.CalcXP(secLvl),
			Type:     secClass,
			Equipped: []int{},
			Evolved:  true,
		}

		user1 := &entities.User{
			Galos: []*entities.Rooster{first},
		}
		user2 := &entities.User{
			Galos: []*entities.Rooster{sec},
		}
		battle := rinha.CreateBattle(user1, user2, false, 0, 0, []*entities.User{}, []string{})
		rounds := 0
		for battle.Fighters[0].Life != 0 && battle.Fighters[1].Life != 0 {
			if log {
				effects := battle.Play(-1, rounds)
				fmt.Println(effects)
				fmt.Printf("%d vs %d\n", battle.Fighters[0].Life, battle.Fighters[1].Life)
			} else {
				battle.Play(-1, rounds)
			}
			rounds++
		}

		if battle.Fighters[1].Life == 0 && battle.Fighters[0].Life == 0 {
			empate++
		} else if battle.Fighters[1].Life == 0 {
			wins++
		} else if battle.Fighters[0].Life == 0 {
			wins1++
		}

	}

	//fmt.Println("\n------------- Analitical tests for ----------- ")
	//fmt.Printf("1: Classe: %s, Level: %d Venceu: %d%%\n2: Classe: %s, Level: %d Venceu: %d%%\n\n", rinha.Classes[firstClass].Name, firstLvl, wins/(times/100), rinha.Classes[secClass].Name, secLvl, wins1/(times/100))
	return wins / (times / 100)
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	arr := []*stats{}
	for i := 1; i < len(rinha.Classes); i++ {
		if rinha.Classes[i].Rarity == 5 {
			arr = append(arr, &stats{
				Type: i,
			})
		}
	}
	/*

		arr = []*stats{&stats{
			Type: 55,
		}}

	*/
	for i := 30; i < 35; i += 1 {
		for _, class := range arr {
			result := measure(class.Type, 55, i, i, 3500, false)
			class.TotalWins += result
			class.Wins = append(class.Wins, result)
			class.Level = append(class.Level, i)
		}
	}
	for _, class := range arr {
		fmt.Printf("Classe: %s\nTotal %d%%\n", rinha.Classes[class.Type].Name, class.TotalWins/len(class.Wins))
	}

	//    graphic(arr)
}
