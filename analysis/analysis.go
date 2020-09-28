package main

import (
	"asura/src/utils/rinha"
    "fmt"
    "encoding/json"
    "io/ioutil"
)

type Result struct{
    Class [2]int
    Level [2]int
    Ties int
    Wins [2]int
    Skills [2][5][2]int
    Frequency [2][5]int
}

func battle(firstClass int, secClass int, firstLvl int, secLvl int, times int) Result {

    result := Result{}

	for time := 0; time < times; time++ {
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
    
        result.Skills = [2][5][2]int{}
    
        for i := 0; i < 2; i++ {
            for j, v := range battle.Fighters[i].Equipped {
                result.Skills[i][j] = [2]int{battle.Fighters[i].Galo.Type,v}
            }
        }
        for battle.Fighters[0].Life > 0 && battle.Fighters[1].Life > 0  {
            battle.Play()
            
        }
    
        if battle.Fighters[1].Life <= 0 && battle.Fighters[0].Life <= 0 {
            result.Ties++
        } else if battle.Fighters[1].Life <= 0  {
            result.Wins[0]++
        } else if battle.Fighters[0].Life <= 0  {
            result.Wins[1]++
        } 
    }

    return result
}

	
func check(e error) {
    if e != nil {
        panic(e)
    }
}

func main(){
    data := []Result{}
    
    for i := 1; i < 50; i++ {
        for j := 1; j <= 4; j++ {    
            res :=  battle(1,j,i,i,10000)
            res.Class = [2]int{1,j}
            res.Level = [2]int{i,i}
            data = append(data, res)
        }
    }

    b, err := json.Marshal(data)
    if err != nil {
        fmt.Println(err)
        return 
    } else {
        err := ioutil.WriteFile("./data.json", b, 0644)
        check(err)
    }
}