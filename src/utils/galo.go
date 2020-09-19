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


type Skill struct {
	Name   string `json:"name"`
	Damage [2]int `json:"damage"`
	Type   int    `json:"type"`
	MinLevel   int    `json:"minlevel"`
}

var Skills []Skill

type Galo struct {
	Name string `json:"name"`
	Xp   int `json:"xp"`
	Type int `json:"type"`
	Skills []int `json:""`
	Equipped []int `json:""`
}

func IdInSkills(a int, arry []int) bool {
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

func RaffleSkill(galo Galo) (*Skill, int) {
	dontHave := []int{}
	
	if len(galo.Skills) == len(Skills) {
		return nil, 0
	}

	level := CalcLevel(galo.Xp)

	for i := 0; i < len(Skills); i++ {
		if !IdInSkills(i, galo.Skills) && !IdInSkills(i, dontHave) && math.Abs(float64(level - Skills[i].MinLevel)) < 2 && (Skills[i].Type == 0 || Skills[i].Type == galo.Type) {
			dontHave = append(dontHave, i)
		}
	}

	if len(dontHave) == 0 {
		for i := 0; i < len(Skills); i++ {
			if !IdInSkills(i, galo.Skills) && !IdInSkills(i, dontHave) && Skills[i].MinLevel - level < 2 && (Skills[i].Type == 0 || Skills[i].Type == galo.Type) {
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