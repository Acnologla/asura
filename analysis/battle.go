package main

import (
	"asura/src/utils/rinha"
	"fmt"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
	"math/rand"
	"time"
)

type stats struct {
	Wins      []int
	TotalWins int
	Level     []int
	Type      int
}

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

	battle := rinha.CreateBattle(&first, &sec, false)

	for battle.Fighters[0].Life != 0 && battle.Fighters[1].Life != 0 {
		turno := 0
		if battle.Turn {
			turno = 1
		}
		effects := battle.Play()
		fmt.Printf("Turno de %d\n%d vs %d\n", turno, battle.Fighters[0].Life, battle.Fighters[1].Life)
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
	} else if battle.Fighters[1].Life == 0 {
		fmt.Println("Primeiro ganhou")
	} else if battle.Fighters[0].Life == 0 {
		fmt.Println("Seugndo ganhou")
	}
}

func measure(firstClass int, secClass int, firstLvl int, secLvl int, times int, log bool) int {
	wins := 0
	wins1 := 0
	empate := 0
	for j := 0; j < times; j++ {
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

		battle := rinha.CreateBattle(&first, &sec, false)

		for battle.Fighters[0].Life != 0 && battle.Fighters[1].Life != 0 {
			if log {
				effects := battle.Play()
				fmt.Println(effects)
				fmt.Printf("%d vs %d\n", battle.Fighters[0].Life, battle.Fighters[1].Life)
			} else {
				battle.Play()
			}
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

func genPts(wins []int, level []int) plotter.XYs {
	pts := make(plotter.XYs, len(wins))
	for i := range pts {
		pts[i].Y = float64(wins[i])
		pts[i].X = float64(level[i])
	}
	return pts
}

func graphic(stat []*stats) {
	p, err := plot.New()
	if err != nil {
		panic(err)
	}

	p.Title.Text = "Rinha classes"
	p.X.Label.Text = "Level"
	p.Y.Label.Text = "% De vitoria "
	arr := []interface{}{}
	for _, class := range stat {
		pts := genPts(class.Wins, class.Level)
		className := rinha.Classes[class.Type]
		arr = append(arr, className.Name, pts)
	}
	err = plotutil.AddLinePoints(p, arr...)
	if err != nil {
		panic(err)
	}
	if err := p.Save(9*vg.Inch, 6*vg.Inch, "points.png"); err != nil {
		panic(err)
	}
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	arr := []*stats{}
	for i := 1; i < len(rinha.Classes); i++ {
		if rinha.Classes[i].Rarity == 1 {
			arr = append(arr, &stats{
				Type: i,
			})
		}
	}
	arr = []*stats{&stats{
		Type: 27,
	}}
	for i := 1; i < 26; i += 1 {
		for _, class := range arr {
			result := measure(class.Type, 1, i, i, 10000, false)
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
