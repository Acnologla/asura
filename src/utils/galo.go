package utils

import (
	"math"
	"asura/src/database"
	"fmt"
	"context"
	"errors"
	"github.com/andersfylling/disgord"
	"io/ioutil"
	"encoding/json"
	"math/rand"
)
type AttackEffect struct {
	Name string `json:"name"`
	Type int `json:"type"`
	Turns int `json:"turns"`
	Damage [2]int `json:"damage"`
}

var AttackEffects []AttackEffect

type Class struct {
	Name string `json:"name"`
	Desc string `json:"desc"`
	Disadvantages []int `json:"disadvantages"`
}

var Classes []Class

type Skill struct {
	Name   string `json:"name"`
	Damage [2]int `json:"damage"`
	Type   int    `json:"type"`
	MinLevel   int    `json:"minlevel"`
	Effect [2] float64 `json:"effect"`
}

var Skills []Skill

type Galo struct {
	Name string `json:"name"`
	Xp   int `json:"xp"`
	Type int `json:"type"`
	Skills []int `json:""`
	Equipped []int `json:""`
}

func IsIntInList(a int, arry []int) bool {
    for _, b := range arry {
        if b == a {
            return true
        }
    }
    return false
}

func init() {
	byteValue, _ := ioutil.ReadFile("./resources/galo/attacks.json")
	json.Unmarshal([]byte(byteValue), &Skills)
	byteValueClass, _ := ioutil.ReadFile("./resources/galo/class.json")
	json.Unmarshal([]byte(byteValueClass), &Classes)
	byteValueEffect, _ := ioutil.ReadFile("./resources/galo/effects.json")
	json.Unmarshal([]byte(byteValueEffect), &AttackEffects)
	
}

func GetGaloDB(id disgord.Snowflake) (Galo, error) {
	var acc Galo
	err := database.Database.NewRef(fmt.Sprintf("galo/%d",id)).Get(context.Background(), &acc);
	if err != nil {
		return acc, errors.New("Not bro")
	}
	return acc, nil
}

func SaveGaloDB(id disgord.Snowflake, galo Galo) {
	database.Database.NewRef(fmt.Sprintf("galo/%d", id)).Set(context.Background(), &galo)
}

func RaffleSkill(galo Galo, exclude []int) (*Skill, int) {
	dontHave := []int{}
	
	if len(galo.Skills) == len(Skills) {
		return nil, 0
	}

	level := CalcLevel(galo.Xp)

	for i := 0; i < len(Skills); i++ {
		if !IsIntInList(i, exclude) && !IsIntInList(i, dontHave) && math.Abs(float64(level - Skills[i].MinLevel)) < 2 && (Skills[i].Type == 0 || Skills[i].Type == galo.Type) {
			dontHave = append(dontHave, i)
		}
	}

	if len(dontHave) == 0 {
		for i := 0; i < len(Skills); i++ {
			if !IsIntInList(i, exclude) && !IsIntInList(i, dontHave) && Skills[i].MinLevel - level < 2 && (Skills[i].Type == 0 || Skills[i].Type == galo.Type) {
				dontHave = append(dontHave, i)
			}
		}
	}

	if len(dontHave) == 0 {
		return nil, 0
	}

	random := rand.Intn(len(dontHave))
	return &Skills[dontHave[random]], dontHave[random] 
}

func CalcLevel(xp int) int {
	return int(math.Floor(math.Sqrt(float64(xp) / 30)))+1
}

func CalcXP(level int) int{
	return int(math.Pow(float64(level-1), 2)) * 30
}

func ChooseSkills(galo *Galo){
	for i := 0; i < 10; i++ {
		skill, id := RaffleSkill(*galo, galo.Skills)
		if skill != nil {
			galo.Skills = append(galo.Skills,id)
		} else {
			break;
		}
	}
}