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

type Class struct {
	Name string `json:"name"`
	Desc string `json:"desc"`
	Disadvantages []int `json:"disadvantages"`
}

type Skill struct {
	Name   string `json:"name"`
	Damage [2]int `json:"damage"`
	Type   int    `json:"type"`
	Level   int    `json:"level"`
	Effect [2] float64 `json:"effect"`
}

type Galo struct {
	Name string `json:"name"`
	Xp   int `json:"xp"`
	Type int `json:"type"`
	Equipped []int `json:"equipped"`
	Ignore bool `json:"ignore"`
}

var AttackEffects []AttackEffect
var Classes []Class
var Skills []Skill
var Sprites [][]string

func GetEffectFromSkill(skill Skill) AttackEffect {
	return AttackEffects[int(skill.Effect[1])]
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
	byteValueSprites, _ := ioutil.ReadFile("./resources/galo/sprites.json")
	json.Unmarshal([]byte(byteValueSprites), &Sprites)
	
} 

func GetGaloDB(id disgord.Snowflake) (Galo, error) {
	var acc Galo
	err := database.Database.NewRef(fmt.Sprintf("galo/%d",id)).Get(context.Background(), &acc);
	if err != nil {
		return acc, errors.New("There's some error")
	}
	return acc, nil
}

func GetSkills(galo Galo) []int {
	skills := []int{}
	lvl := CalcLevel(galo.Xp)
	for i := 0; i < len(Skills); i++ {
		if Skills[i].Level > lvl || (Skills[i].Type != galo.Type && Skills[i].Type != 1) {
			continue;
		}
		skills = append(skills, i)
	}
	return skills
}

func GetNextSkill(galo Galo) []Skill {
	skills := []Skill{}
	lvl := CalcLevel(galo.Xp)
	for i := 0; i < len(Skills); i++ {
		if Skills[i].Level == lvl+1 && (Skills[i].Type != galo.Type && Skills[i].Type != 1) {
			skills = append(skills, Skills[i])
		}
	}
	return skills
}


func SaveGaloDB(id disgord.Snowflake, galo Galo) {
	database.Database.NewRef(fmt.Sprintf("galo/%d", id)).Set(context.Background(), &galo)
}

// Math functions

func CalcLevel(xp int) int {
	return int(math.Floor(math.Sqrt(float64(xp) / 30)))+1
}

func CalcXP(level int) int{
	return int(math.Pow(float64(level-1), 2)) * 30
}

func Between(damage [2]int) int {
	return rand.Intn(damage[1] - damage[0]) + damage[0]
}

func SaturateSub(one int, two int) int {
	if two >= one {
		return 0
	} else {
		return one - two
	}
}