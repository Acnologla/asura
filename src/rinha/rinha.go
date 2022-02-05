package rinha

import (
	"asura/src/entities"
	"asura/src/utils"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"os"
	"strings"
	"time"
)

type Rarity int

var client = &http.Client{}

const (
	Common Rarity = iota
	Rare
	Epic
	Legendary
	Special
	Mythic
)

func (rarity Rarity) String() string {
	return [...]string{"Comum", "Raro", "Epico", "Lendario", "Especial", "Mitico"}[rarity]
}

func (rarity Rarity) Price() int {
	return [...]int{30, 140, 480, 1200, 500, 3000}[rarity]
}

func (rarity Rarity) Color() int {
	return [...]int{13493247, 255, 9699539, 16748544, 16728128, 16777201}[rarity]
}

type Vote struct {
	Voted int `json:"voted"`
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

var Dungeon []*Room
var Items []*Item
var Effects []*Effect
var Classes []*Class
var Skills []([]*Skill)
var Sprites [][]string
var Cosmetics []*Cosmetic
var Upgrades []Upgrade

func init() {
	str, _ := os.Getwd()
	if strings.HasSuffix(str, "test") {
		os.Chdir("..")
	}
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
	byteValueDungeon, _ := ioutil.ReadFile("./resources/galo/dungeon.json")
	json.Unmarshal([]byte(byteValueDungeon), &Dungeon)
	byteValueItems, _ := ioutil.ReadFile("./resources/galo/items.json")
	json.Unmarshal([]byte(byteValueItems), &Items)
	byteValueCosmetics, _ := ioutil.ReadFile("./resources/galo/cosmetics.json")
	json.Unmarshal([]byte(byteValueCosmetics), &Cosmetics)
	byteValueUpgrades, _ := ioutil.ReadFile("./resources/galo/upgrades.json")
	json.Unmarshal([]byte(byteValueUpgrades), &Upgrades)
}

const PityMultiplier = 1

func VipMessage(user *entities.User) string {
	now := uint64(time.Now().Unix())
	if now >= user.Vip {
		return ""
	}
	return fmt.Sprintf("Vip por **%d** dias", (user.Vip-now)/60/60/24)
}

func findClassIndex(class string) int {
	for i, classObj := range Classes {
		if classObj.Name == class {
			return i - 1
		}
	}
	return -1
}
func GetRandByType(classType Rarity) int {
	classTypeArr := []*Class{}
	for _, class := range Classes {
		if class.Rarity == classType {
			classTypeArr = append(classTypeArr, class)
		}
	}
	selected := classTypeArr[utils.RandInt(len(classTypeArr))]
	for i, class := range Classes {
		if class.Name == selected.Name {
			return i
		}
	}
	return -1
}

func GetCommonOrRare() int {
	rand := utils.RandInt(100)
	if 4 >= rand {
		return GetRandByType(Rare)
	}
	return GetRandByType(Common)
}

func HaveRooster(galos []*entities.Rooster, galoType int) bool {
	for _, galo := range galos {
		if galo.Type == galoType {
			return true
		}
	}
	return false

}
func CalcLevel(xp int) int {
	return int(math.Floor(math.Sqrt(float64(xp)/30))) + 1
}
func Sell(rarity Rarity, xp int, reset int) (int, int) {
	level := float64(CalcLevel(xp)+(reset*30)) - 1
	price := float64(rarity.Price())
	if reset == 0 {
		return int(price * (level/5 + 1)), 0
	}
	asuraCoins := reset
	if rarity > Rare {
		asuraCoins++
	}
	if rarity == Legendary {
		asuraCoins += 2
	}
	return 0, asuraCoins
}
