package rinha

import (
	"asura/src/entities"
	"fmt"
	"math"
	"strings"
	"time"
)

const (
	MissionN     = 20000
	MissionMoney = 625
	MissionXp    = 2125
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
		text += "10% de xp adicional para o ovo\n"
	}
	if level >= 15 {
		text += "Chance adicional de dropar chaves\n"
	}
	return
}

func CalcMissionPrize(clan *entities.Clan) (int, int) {
	return MissionMoney + (10 * clan.MissionsUpgrade), MissionXp + (50 * clan.MissionsUpgrade)
}

func MissionToString(clan *entities.Clan) string {
	money, xp := CalcMissionPrize(clan)
	done := clan.MissionProgress >= MissionN

	if done {
		need := uint64(time.Now().Unix()) - clan.Mission
		return fmt.Sprintf("Espere mais %d dias e %d horas para seu clan receber uma nova missão", 30-(need/60/60/24), 23-(need/60/60%24))
	} else {
		return fmt.Sprintf("Derrote %d/%d galos na rinha\nMoney: **%d**\nXp: **%d**", clan.MissionProgress, MissionN, money, xp)
	}
}

func PopulateClanMissions(clan *entities.Clan) *entities.Clan {
	if uint64((uint64(time.Now().Unix())-clan.Mission)/60/60/24) >= 30 {
		clan.Mission = uint64(time.Now().Unix())
		clan.MissionProgress = 0
	}
	return clan
}

func CalcClanUpgrade(x, price int) int {
	return int(math.Pow(2, float64(x)) * (float64(price * 500)))
}
