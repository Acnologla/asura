package rinha

import (
	"asura/src/database"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/andersfylling/disgord"
	"io/ioutil"
	"math"
	"math/rand"
	"strings"
)

type Rarity int

const (
	Common Rarity = iota
	Rare
	Epic
	Legendary
)

func (rarity Rarity) String() string {
	return [...]string{"Comum", "Raro", "Epico", "Lendario"}[rarity]
}

func (rarity Rarity) Color() int {
	return [...]int{13493247, 255, 9699539, 16748544}[rarity]
}

type Effect struct {
	Name   string `json:"name"`
	Class  int    `json:"class"`
	Type   int    `json:"type"`
	Self   bool   `json:"self"`
	Phrase string `json:"phrase"`
	Turns  int    `json:"turns"`
	Range  [2]int `json:"range"`
}

type Class struct {
	Name          string `json:"name"`
	Desc          string `json:"desc"`
	Disadvantages []int  `json:"disadvantages"`
	Rarity        Rarity `json:"rarity"`
}

type Skill struct {
	Name   string     `json:"name"`
	Damage [2]int     `json:"damage"`
	Level  int        `json:"level"`
	Effect [2]float64 `json:"effect"`
	Self   bool       `json:"self"`
}

type SubGalo struct{
	Type     int    `json:"type"`
	Xp       int    `json:"xp"`

}

type Galo struct {
	Name     string `json:"name"`
	Xp       int    `json:"xp"`
	Type     int    `json:"type"`
	Equipped []int  `json:"equipped"`
	Win      int    `json:"win"`
	Lose     int    `json:"lose"`
	Lootbox  int    `json:"lootbox"`
	Galos    []SubGalo  `json:"galos"`
	Money    int    `json:"money"`
}

var Effects []*Effect
var Classes []*Class
var Skills []([]*Skill)
var Sprites [][]string

func init() {

	byteValueClass, _ := ioutil.ReadFile("./resources/galo/class.json")
	json.Unmarshal([]byte(byteValueClass), &Classes)
	for i := 0; i < len(Classes)-1; i++ {
		Skills = append(Skills, []*Skill{})
	}
	atacks, _ := ioutil.ReadDir("./resources/galo/attacks")
	for _, class := range atacks {
		byteValueAtack, _ := ioutil.ReadFile(fmt.Sprintf("./resources/galo/attacks/%s", class.Name()))
		name := strings.Split(class.Name(), ".")[0]
		index := findClassIndex(name)
		if index != -1 {
			skils := []*Skill{}
			json.Unmarshal([]byte(byteValueAtack), &skils)
			Skills[index] = skils
		}
	}
	byteValueEffect, _ := ioutil.ReadFile("./resources/galo/effects.json")
	json.Unmarshal([]byte(byteValueEffect), &Effects)
	byteValueSprites, _ := ioutil.ReadFile("./resources/galo/sprites.json")
	json.Unmarshal([]byte(byteValueSprites), &Sprites)
}

func findClassIndex(class string) int {
	for i, classObj := range Classes {
		if classObj.Name == class {
			return i - 1
		}
	}
	return -1
}

func SkillToString(skill *Skill) (text string) {
	text = fmt.Sprintf("Dano: %d - %d", skill.Damage[0], skill.Damage[1]-1)
	if skill.Effect[0] != 0 || skill.Effect[1] != 0 {
		effect := Effects[int(skill.Effect[1])]
		text += fmt.Sprintf("\nTem %d%% de Chance de causar %s", int(skill.Effect[0]*100), effect.Name)
	}
	return
}

// Database manipulation

func GetGaloDB(id disgord.Snowflake) (Galo, error) {
	var acc Galo
	err := database.Database.NewRef(fmt.Sprintf("galo/%d", id)).Get(context.Background(), &acc)
	if err != nil {
		return acc, errors.New("There's some error")
	}
	return acc, nil
}

func SaveGaloDB(id disgord.Snowflake, galo Galo) {
	database.Database.NewRef(fmt.Sprintf("galo/%d", id)).Set(context.Background(), &galo)
}

func GetSkills(galo Galo) []int {
	skills := []int{}
	lvl := CalcLevel(galo.Xp)
	if galo.Type == 0 {
		return skills
	}
	for i := 0; i < len(Skills[galo.Type-1]); i++ {
		if Skills[galo.Type-1][i].Level > lvl {
			continue
		}
		skills = append(skills, i)
	}
	return skills
}

func GetNextSkill(galo Galo) []*Skill {
	skills := []*Skill{}
	lvl := CalcLevel(galo.Xp)
	for i := 0; i < len(Skills[galo.Type-1]); i++ {
		if Skills[galo.Type-1][i].Level == lvl {
			skills = append(skills, Skills[galo.Type-1][i])
		}
	}
	return skills
}

// Helper functions

func IsIntInList(a int, arry []int) bool {
	for _, b := range arry {
		if b == a {
			return true
		}
	}
	return false
}

func CalcLevel(xp int) int {
	return int(math.Floor(math.Sqrt(float64(xp)/30))) + 1
}

func HaveGalo(galo int,galos []SubGalo) bool{
	for _,gal := range galos{
		if gal.Type == galo{
			return true
		}
	}
	return false
}

func CalcXP(level int) int {
	return int(math.Pow(float64(level-1), 2)) * 30
}

func Between(damage [2]int) int {
	if damage[1] == damage[0] {
		return damage[1]
	}
	return rand.Intn(damage[1]-damage[0]) + damage[0]
}

func GetRandByType(classType Rarity) int {
	classTypeArr := []*Class{}
	for _, class := range Classes {
		if class.Rarity == classType {
			classTypeArr = append(classTypeArr, class)
		}
	}
	selected := classTypeArr[rand.Intn(len(classTypeArr))]
	for i, class := range Classes {
		if class.Name == selected.Name {
			return i
		}
	}
	return -1
}

func SaturateSub(one int, two int) int {
	if two >= one {
		return 0
	} else {
		return one - two
	}
}

// Effect functions

func GetEffectFromSkill(skill *Skill) *Effect {
	return Effects[int(skill.Effect[1])]
}
