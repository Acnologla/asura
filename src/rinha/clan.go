package rinha

import (
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
