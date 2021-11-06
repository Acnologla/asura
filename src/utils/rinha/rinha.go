package rinha

import (
	"asura/src/database"
	"asura/src/utils"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"strings"
	"time"

	"firebase.google.com/go/db"
	"github.com/andersfylling/disgord"
)

type Rarity int

const (
	Common Rarity = iota
	Rare
	Epic
	Legendary
	Special
)

func (rarity Rarity) String() string {
	return [...]string{"Comum", "Raro", "Epico", "Lendario", "Especial"}[rarity]
}

func (rarity Rarity) Price() int {
	return [...]int{30, 140, 500, 1600, 500}[rarity]
}

func (rarity Rarity) Color() int {
	return [...]int{13493247, 255, 9699539, 16748544, 16728128}[rarity]
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

type SubGalo struct {
	Type      int    `json:"type"`
	Xp        int    `json:"xp"`
	Name      string `json:"name"`
	GaloReset int    `json:"galoReset"`
}

type TrainLimit struct {
	Times     int    `json:"times"`
	LastReset uint64 `json:"lastReset"`
}

type Arena struct {
	Active    bool              `json:"active"`
	Win       int               `json:"wins"`
	Lose      int               `json:"lose"`
	LastFight disgord.Snowflake `json:"lastFight"`
}

type Daily struct {
	Last   uint64 `json:"last"`
	Strike int    `json:"strike"`
}

type Galo struct {
	UserXp           int        `json:"userXp"`
	Upgrades         []int      `json:"upgrades"`
	GaloReset        int        `json:"galoReset"`
	Name             string     `json:"name"`
	Xp               int        `json:"xp"`
	Type             int        `json:"type"`
	Equipped         []int      `json:"equipped"`
	Win              int        `json:"win"`
	Lose             int        `json:"lose"`
	Lootbox          int        `json:"lootbox"`
	CommonLootbox    int        `json:"commonLootbox"`
	RareLootbox      int        `json:"rareLootbox"`
	EpicLootbox      int        `json:"epicLootbox"`
	Pity             int        `json:"pity"`
	LegendaryLootbox int        `json:"legendaryLootbox"`
	ItemsLootbox     int        `json:"itemsLootbox"`
	Galos            []SubGalo  `json:"galos"`
	Daily            Daily      `json:"realDaily"`
	Money            int        `json:"money"`
	Clan             string     `json:"clan"`
	Dungeon          int        `json:"dungeon"`
	DungeonReset     int        `json:"dungeonreset"`
	Items            []int      `json:"items"`
	Missions         []Mission  `json:"missions"`
	MissionTrade     uint64     `json:"missionTrade"`
	LastMission      uint64     `json:"lastMission"`
	Vip              uint64     `json:"vip"`
	Cosmetics        []int      `json:"cosmetics"`
	Background       int        `json:"bg"`
	CosmeticLootbox  int        `json:"cosmeticLootbox"`
	TrainLimit       TrainLimit `json:"trainLimit"`
	Arena            Arena      `json:"arena"`
	AsuraCoin        int        `json:"asuraCoin"`
}

var Dungeon []*Room
var Items []*Item
var Effects []*Effect
var Classes []*Class
var Skills []([]*Skill)
var Sprites [][]string
var Cosmetics []*Cosmetic

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

func IsVip(galo Galo) bool {
	return uint64(time.Now().Unix()) <= galo.Vip
}

func findClassIndex(class string) int {
	for i, classObj := range Classes {
		if classObj.Name == class {
			return i - 1
		}
	}
	return -1
}

func CalcDamage(skill *Skill, galo Galo) (min int, max int) {
	min = skill.Damage[0]
	max = skill.Damage[1]
	if galo.GaloReset > 0 {
		min += min / 10 * galo.GaloReset
		max += max / 10 * galo.GaloReset
	}
	return
}
func CalcEffectRange(effect *Effect, galo Galo) (min int, max int) {
	min = effect.Range[0]
	max = effect.Range[1]
	if galo.GaloReset > 0 {
		min += min / 10 * galo.GaloReset
		max += max / 10 * galo.GaloReset
	}
	return
}

func SkillToStringFormated(skill *Skill, galo Galo) (text string) {
	min, max := CalcDamage(skill, galo)
	text = fmt.Sprintf("`[Dano: %d - %d]`", min, max-1)
	if skill.Effect[0] != 0 || skill.Effect[1] != 0 {
		effect := Effects[int(skill.Effect[1])]
		text += fmt.Sprintf("\n*Tem %d%% de Chance de causar %s*", int(skill.Effect[0]*100), effect.Name)
	}
	return
}

func SkillToString(skill *Skill) (text string, effectText string) {
	text = skill.Name
	if skill.Effect[0] != 0 || skill.Effect[1] != 0 {
		effect := Effects[int(skill.Effect[1])]
		effectText = fmt.Sprintf("\n %d%% de causar %s", int(skill.Effect[0]*100), effect.Name)
	}
	return
}

// Database manipulation

func GetGaloDB(id disgord.Snowflake) (Galo, error) {
	var acc Galo
	err := database.Database.NewRef(fmt.Sprintf("galo/%d", id)).Get(context.Background(), &acc)
	if err != nil {
		return acc, errors.New("THERES SOME ERROR")
	}
	return acc, nil
}

func ChangeMoney(id disgord.Snowflake, money int, onlyIf int) error {
	fn := func(tn db.TransactionNode) (interface{}, error) {
		var galo Galo
		err := tn.Unmarshal(&galo)
		if err == nil {
			if galo.Money >= onlyIf {
				galo.Money += money
				return galo, nil
			}
			return nil, errors.New("DONT HAVE MONEY")
		}
		fmt.Println("ChangeMoney transaction error\n", err)
		return nil, err
	}
	return database.Database.NewRef(fmt.Sprintf("galo/%d", id)).Transaction(context.Background(), fn)
}

func UpdateGaloDB(id disgord.Snowflake, callback func(galo Galo) (Galo, error)) {
	fn := func(tn db.TransactionNode) (interface{}, error) {
		var galo Galo
		err := tn.Unmarshal(&galo)
		if err == nil {
			return callback(galo)
		}
		fmt.Println("UpdateGaloDB transaction error\n", err)
		return nil, err
	}
	database.Database.NewRef(fmt.Sprintf("galo/%d", id)).Transaction(context.Background(), fn)
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

func HaveGalo(galo int, galos []SubGalo) bool {
	for _, gal := range galos {
		if gal.Type == galo {
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
	return utils.RandInt(damage[1]-damage[0]) + damage[0]
}

func GetName(username string, galo Galo) string {
	if galo.Name == "" {
		return username
	}
	return galo.Name
}

func GetRand() int {
	return utils.RandInt(len(Classes)-1) + 1
}

func GetRarityPlusOne(rarity Rarity) int {
	classTypeArr := []*Class{}
	for _, class := range Classes {
		if class.Rarity == rarity || class.Rarity == rarity+1 {
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

// Effect functions
func GetEffectFromIndex(idx int) *Effect {
	return Effects[idx]
}

func IsInLimit(galo Galo, id disgord.Snowflake) bool {
	max := 200
	if HasUpgrade(galo.Upgrades, 0, 1, 0, 0) {
		max += 30
	}
	if galo.Clan != "" {
		clan := GetClan(galo.Clan)
		level := ClanXpToLevel(clan.Xp)
		if level >= 10 {
			max += 30
		}
	}
	if galo.TrainLimit.LastReset == 0 || 1 <= (uint64(time.Now().Unix())-galo.TrainLimit.LastReset)/60/60/24 {
		UpdateGaloDB(id, func(galo Galo) (Galo, error) {
			galo.TrainLimit.LastReset = uint64(time.Now().Unix())
			galo.TrainLimit.Times = 0
			return galo, nil
		})
	} else if galo.TrainLimit.Times >= max {
		return true
	}
	return false
}
