package rinha

import (
	"asura/src/entities"
	"math"
	"strings"
)

const allowedChars = "abcdefghijklmnopqrstuvwxyz123456789 -_"

const MaxMoney = 10000000

func Format(text string) string {
	str := strings.TrimSpace(strings.ToLower(text))
	for _, char := range str {
		if !includesString(char, allowedChars) {
			str = strings.Replace(str, string(char), "", 1)
		}
	}
	return strings.TrimSpace(str)
}

func ClanXpToLevel(xp int) int {
	return int(math.Floor(math.Sqrt(float64(xp)/4500))) + 1
}
func ClanLevelToXp(level int) int {
	return int(math.Pow(float64(level), 2)) * 4500
}

func includesString(strOne rune, strTwo string) bool {
	for _, char := range strTwo {
		if char == strOne {
			return true
		}
	}
	return false
}

func GetMaxMembers(clan *entities.Clan) int {
	return 25 + (clan.MembersUpgrade * 2)
}

func GetBenefits(xp int) (text string) {
	level := ClanXpToLevel(xp)
	text = "10% de xp adicional por rinha ganha\n"
	if level >= 2 {
		text += "10% de xp adicional por rinha ganha\n"
	}
	if level >= 3 {
		text += "Maior chance de galos raros nas caixas\n"
	}
	if level >= 4 {
		text += "1 de ouro adicional por rinha ganha\n"
	}
	if level >= 6 {
		text += "1 de ouro adicional por rinha ganha\n"
	}
	if level >= 8 {
		text += "1 de xp de upgrade a mais por rinha ganha\n"
	}
	if level >= 10 {
		text += "Aumenta o limite de trains"
	}
	return
}
