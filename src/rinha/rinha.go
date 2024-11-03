package rinha

import (
	"asura/src/entities"
	"asura/src/utils"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/andersfylling/disgord"
	"github.com/google/uuid"
)

type Rarity int

const BASE_LIMIT = 270

var client = &http.Client{}
var TopToken string

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
	return [...]int{30, 160, 500, 1200, 500, 3000}[rarity]
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
	Name    string     `json:"name"`
	Damage  [2]int     `json:"damage"`
	Level   int        `json:"level"`
	Effect  [2]float64 `json:"effect"`
	Self    bool       `json:"self"`
	Evolved bool       `json:"evolved"`
}

var Dungeon []*Room
var Items []*Item
var Effects []*Effect
var Classes []*Class
var Skills []([]*Skill)
var Sprites [][]string
var Cosmetics []*Cosmetic
var Upgrades []Upgrade
var BattlePass []BattlePassLevel
var Ranks []*Rank

func SetTopToken(token string) {
	TopToken = token
}
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
	var Skins []*Cosmetic
	var NewCosmetics []*Cosmetic
	byteValueSkins, _ := ioutil.ReadFile("./resources/galo/skins.json")
	json.Unmarshal([]byte(byteValueSkins), &Skins)
	Cosmetics = append(Cosmetics, Skins...)

	byteValueNewCosmetics, _ := ioutil.ReadFile("./resources/galo/newCosmetics.json")
	json.Unmarshal([]byte(byteValueNewCosmetics), &NewCosmetics)
	Cosmetics = append(Cosmetics, NewCosmetics...)

	byteValueUpgrades, _ := ioutil.ReadFile("./resources/galo/upgrades.json")
	json.Unmarshal([]byte(byteValueUpgrades), &Upgrades)
	byteValueBattlePass, _ := ioutil.ReadFile("./resources/galo/battlePass.json")
	json.Unmarshal([]byte(byteValueBattlePass), &BattlePass)
	byteValueRanks, _ := ioutil.ReadFile("./resources/galo/ranks.json")
	json.Unmarshal([]byte(byteValueRanks), &Ranks)
}

const PityMultiplier = 1

func IsVip(user *entities.User) bool {
	return uint64(time.Now().Unix()) <= user.Vip
}

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

func GetRandMythic() int {
	classTypeArr := []*Class{}
	for i, class := range Classes {
		if class.Rarity == Mythic && i != 50 {
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
	asuraCoins := reset + 1
	if rarity > Rare {
		asuraCoins++
	}
	if rarity == Legendary {
		asuraCoins += 2
	}
	if rarity == Mythic {
		asuraCoins += 4
	}
	return 0, asuraCoins
}

func GetEquippedGalo(user *entities.User) *entities.Rooster {
	for _, galo := range user.Galos {
		if galo.Equip {
			return galo
		}
	}
	if os.Getenv("PRODUCTION") != "" {
		log.Println("excuse me wtf")
	}
	return user.Galos[0]
}

func GetName(username string, galo entities.Rooster) string {
	prefix := ""

	if galo.Evolved {
		prefix = "[⭐] "
	}

	if galo.Name == "" {
		return prefix + username
	}

	return prefix + galo.Name
}

func GetRoosterByID(galos []*entities.Rooster, id uuid.UUID) *entities.Rooster {
	for _, galo := range galos {
		if galo.ID == id {
			return galo
		}
	}
	return nil
}

func GetRarity(galo *entities.Rooster) Rarity {
	if galo.EvolvedRarity != 0 {
		return Rarity(galo.EvolvedRarity)
	}
	return Classes[galo.Type].Rarity
}

func IsIntInList(a int, arry []int) bool {
	for _, b := range arry {
		if b == a {
			return true
		}
	}
	return false
}

func GetSkills(galo entities.Rooster) []int {
	skills := []int{}
	lvl := CalcLevel(galo.Xp)
	if galo.Type == 0 {
		return skills
	}
	for i := 0; i < len(Skills[galo.Type-1]); i++ {
		skill := Skills[galo.Type-1][i]
		if skill.Level > lvl {
			continue
		}

		if skill.Evolved && !galo.Evolved {
			continue
		}

		skills = append(skills, i)
	}
	return skills
}
func SkillToStringFormated(skill *Skill, galo *entities.Rooster) (text string) {
	min, max := CalcDamage(skill, galo)
	text = fmt.Sprintf("`[Dano: %d - %d]`", min, max-1)
	if skill.Effect[0] != 0 || skill.Effect[1] != 0 {
		effect := Effects[int(skill.Effect[1])]
		text += fmt.Sprintf("\n*Tem %d%% de Chance de causar %s*", int(skill.Effect[0]*100), effect.Name)
	}
	return
}

func calcEvolvedRarityMultiplier(galo *entities.Rooster) (multiplier float64) {
	class := Classes[galo.Type]
	for i := class.Rarity; i < Rarity(galo.EvolvedRarity); i++ {
		if i == Common {
			multiplier += 0.15
		}
		if i == Rare {
			multiplier += 0.2
		}
		if i == Epic {
			multiplier += 0.35
		}
		if i == Legendary {
			multiplier += 0.6
		}
	}

	return
}

func calcDamage(min, max int, galo *entities.Rooster) (int, int) {
	if galo.Resets > 0 {
		min += int(float64(min) * 0.15 * float64(galo.Resets))
		max += int(float64(max) * 0.15 * float64(galo.Resets))
	}

	if galo.EvolvedRarity > 0 {
		multiplier := calcEvolvedRarityMultiplier(galo)
		min += int(float64(min) * multiplier)
		max += int(float64(max) * multiplier)
	}

	return min, max
}

func CalcDamage(skill *Skill, galo *entities.Rooster) (min int, max int) {
	min = skill.Damage[0]
	max = skill.Damage[1]
	return calcDamage(min, max, galo)
}
func CalcEffectRange(effect *Effect, galo *entities.Rooster) (min int, max int) {
	min = effect.Range[0]
	max = effect.Range[1]
	return calcDamage(min, max, galo)
}

func Between(damage [2]int) int {
	if damage[1] == damage[0] {
		return damage[1]
	}
	return utils.RandInt(damage[1]-damage[0]) + damage[0]
}

func CalcXP(level int) int {
	if 0 >= level {
		return 1
	}
	return int(math.Pow(float64(level-1), 2)) * 30
}

func GetNextSkill(galo entities.Rooster) []*Skill {
	skills := []*Skill{}
	lvl := CalcLevel(galo.Xp)
	for i := 0; i < len(Skills[galo.Type-1]); i++ {
		if Skills[galo.Type-1][i].Level == lvl {
			skills = append(skills, Skills[galo.Type-1][i])
		}
	}
	return skills
}

func GetNewRooster() int {
	rand := utils.RandInt(100)
	if 4 >= rand {
		return GetRandByType(Epic)
	} else if 10 >= rand {
		return GetRandByType(Rare)
	}
	return GetRandByType(Common)
}

func GetCommonOrRare() int {
	rand := utils.RandInt(100)
	if 4 >= rand {
		return GetRandByType(Rare)
	}
	return GetRandByType(Common)
}

func HaveGalo(galos []*entities.Rooster, id uuid.UUID) bool {
	for _, galo := range galos {
		if galo.ID == id {
			return true
		}
	}
	return false
}

func GetGaloByID(galos []*entities.Rooster, id uuid.UUID) *entities.Rooster {
	for _, galo := range galos {
		if galo.ID == id {
			return galo
		}
	}
	return nil
}
func GetRand() int {
	return utils.RandInt(len(Classes)-1) + 1
}

func GetTrialsMultiplier(user *entities.User) float64 {
	rooster := GetEquippedGalo(user)
	for _, trial := range user.Trials {
		if trial.Rooster == rooster.Type {
			return 1 + (float64(trial.Win) * 0.04)
		}
	}
	return 1
}

func GetItemByID(items []*entities.Item, id uuid.UUID) *entities.Item {
	for _, item := range items {
		if item.ID == id {
			return item
		}
	}
	return nil
}

func CalcLimit(user *entities.User) int {
	limit := BASE_LIMIT
	if HasUpgrade(user.Upgrades, 0, 1, 0, 0) {
		limit += 40
	}
	return limit
}

func IsInLimit(user *entities.User) bool {
	return user.TrainLimit >= CalcLimit(user)
}

func SkillToString(skill *Skill) (text string, effectText string) {
	text = skill.Name
	if skill.Effect[0] != 0 || skill.Effect[1] != 0 {
		effect := Effects[int(skill.Effect[1])]
		effectText = fmt.Sprintf("\n %d%% de causar %s", int(skill.Effect[0]*100), effect.Name)
	}
	return
}

func HasVoted(id disgord.Snowflake) bool {
	req, err := http.NewRequest("GET", fmt.Sprintf("https://top.gg/api/bots/%d/check?userId=%d", 470684281102925844, id), nil)
	if err != nil {
		return false
	}
	req.Header.Add("Authorization", TopToken)
	resp, err := client.Do(req)
	if err == nil {
		defer resp.Body.Close()
		var vote Vote
		json.NewDecoder(resp.Body).Decode(&vote)
		return vote.Voted == 1
	}
	return false
}
