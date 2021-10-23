package rinha

import (
	"asura/src/database"
	"context"
	"fmt"
	"math"
	"strings"
	"time"

	"firebase.google.com/go/db"
	"github.com/andersfylling/disgord"
)

const allowedChars = "abcdefghijklmnopqrstuvwxyz123456789 -_"

const MaxMoney = 10000000

type Role uint

const (
	Member Role = iota
	Admin
	Owner
)

func (role Role) ToString() string {
	return [...]string{"Membro", "Administrador", "Dono"}[role]
}

type ClanMember struct {
	ID   uint64 `json:"id"`
	Role Role   `json:"role"`
	Xp   uint   `json:"xp"`
}

type ClanUpgrades struct {
	Members int `json:"members"`
	Banks   int `json:"banks"`
	Mission int `json:"mission"`
}

type Clan struct {
	Xp              int          `json:"xp"`
	CreatedAt       uint64       `json:"createdAt"`
	Members         []ClanMember `json:"members"`
	Money           int          `json:"money"`
	Upgrades        ClanUpgrades `json:"upgrades"`
	LastIncome      uint64       `json:"lastIncome"`
	Mission         uint64       `json:"mission"`
	MissionProgress int          `json:"missionProgress"`
}

func GetClan(name string) Clan {
	text := Format(name)
	var clan Clan
	database.Database.NewRef("clan/"+text).Get(context.Background(), &clan)
	return clan
}

func CreateClan(name string, owner disgord.Snowflake) {
	ownerMember := ClanMember{
		ID:   uint64(owner),
		Role: Owner,
	}
	clan := Clan{
		Members:   []ClanMember{ownerMember},
		CreatedAt: uint64(time.Now().Unix()),
	}
	database.Database.NewRef("clan/"+name).Set(context.Background(), &clan)
}

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
	return int(math.Floor(math.Sqrt(float64(xp)/4000))) + 1
}
func ClanLevelToXp(level int) int {
	return int(math.Pow(float64(level), 2)) * 4000
}

func GetMember(clan Clan, id disgord.Snowflake) ClanMember {
	for _, member := range clan.Members {
		if member.ID == uint64(id) {
			return member
		}
	}
	return ClanMember{}
}

func IsInClan(clan Clan, id disgord.Snowflake) bool {
	for _, member := range clan.Members {
		if member.ID == uint64(id) {
			return true
		}
	}
	return false
}

func DeleteClan(clan string) {
	database.Database.NewRef("clan/" + clan).Delete(context.Background())
}

func UpdateClan(clan string, callback func(clan Clan) (Clan, error)) {
	fn := func(tn db.TransactionNode) (interface{}, error) {
		var clan Clan
		err := tn.Unmarshal(&clan)
		if err == nil {
			return callback(clan)
		}
		fmt.Println("UpdateClan transaction error\n", err)
		return nil, err
	}
	database.Database.NewRef("clan/"+clan).Transaction(context.Background(), fn)
}
func FindMemberIndex(clan Clan, memberID disgord.Snowflake) int {
	for i, member := range clan.Members {
		if member.ID == uint64(memberID) {
			return i
		}
	}
	return -1
}

func PromoteMember(clan Clan, memberID disgord.Snowflake) []ClanMember {
	index := FindMemberIndex(clan, memberID)
	if index == -1 {
		return clan.Members
	}
	clan.Members[index].Role = Admin
	return clan.Members
}

func RemoveMember(clan Clan, memberID disgord.Snowflake) []ClanMember {
	index := FindMemberIndex(clan, memberID)
	if index == -1 {
		return clan.Members
	}
	for i := index; i < len(clan.Members)-1; i++ {
		clan.Members[i] = clan.Members[i+1]
	}
	clan.Members = clan.Members[0 : len(clan.Members)-1]
	return clan.Members
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
		text += "2 de ouro adicional por rinha ganha\n"
	}
	if level >= 8 {
		text += "1 de xp de upgrade a mais por rinha ganha\n"
	}
	if level >= 10 {
		text += "Aumenta o limite de trains"
	}
	return
}
func includesString(strOne rune, strTwo string) bool {
	for _, char := range strTwo {
		if char == strOne {
			return true
		}
	}
	return false
}

func PopulateClanMissions(clan Clan, name string, save bool) Clan {
	if uint64((uint64(time.Now().Unix())-clan.Mission)/60/60/24) >= 30 {
		clan.Mission = uint64(time.Now().Unix())
		clan.MissionProgress = 0
		if save {
			UpdateClan(name, func(clanUpdate Clan) (Clan, error) {
				clanUpdate.Mission = clan.Mission
				clanUpdate.MissionProgress = clan.MissionProgress
				return clanUpdate, nil
			})
		}
	}
	return clan
}

const (
	missionN     = 20000
	missionMoney = 625
	missionXp    = 2125
)

func calcMissionPrize(clan Clan) (int, int) {
	return missionMoney + (10 * clan.Upgrades.Mission), missionXp + (50 * clan.Upgrades.Mission)
}

func MissionToString(clan Clan) string {
	money, xp := calcMissionPrize(clan)
	done := clan.MissionProgress >= missionN

	if done {
		need := uint64(time.Now().Unix()) - clan.Mission
		return fmt.Sprintf("Espere mais %d dias e %d horas para seu clan receber uma nova missão", 30-(need/60/60/24), 23-(need/60/60%24))
	} else {
		return fmt.Sprintf("Derrote %d/%d galos na rinha\nMoney: **%d**\nXp: **%d**", clan.MissionProgress, missionN, money, xp)
	}
}

func CompleteClanMission(clanName string, id disgord.Snowflake) {
	UpdateClan(clanName, func(clan Clan) (Clan, error) {
		clan = PopulateClanMissions(clan, clanName, false)
		clan.Xp++
		for i, member := range clan.Members {
			if member.ID == uint64(id) {
				member.Xp++
				clan.Members[i] = member
				break
			}
		}
		money, xp := calcMissionPrize(clan)
		done := clan.MissionProgress >= missionN
		if !done {
			clan.MissionProgress++
			if clan.MissionProgress >= missionN {
				for _, member := range clan.Members {
					UpdateGaloDB(disgord.Snowflake(member.ID), func(galo Galo) (Galo, error) {
						galo.Xp += xp
						galo.Money += money
						return galo, nil
					})
				}
			}
		}
		return clan, nil
	})

}

func GetMaxMembers(clan Clan) int {
	return 15 + clan.Upgrades.Members
}

func CalcClanUpgrade(x, price int) int {
	return int(math.Pow(2, float64(x)) * (float64(price * 1000)))
}

func CalcClanIncomeTime(clan Clan) int {
	return 50 - int((uint64(time.Now().Unix())-clan.LastIncome)/60/60)
}
func UpdateClanBank(clan Clan, clanName string) Clan {
	income := int((uint64(time.Now().Unix()) - clan.LastIncome) / 60 / 60 / 50)
	if income == 0 {
		return clan
	}
	if clan.LastIncome == 0 {
		income = 1
	}
	money := 0
	for i := 0; i < income; i++ {
		val := int(float64(clan.Money) * (float64(1+clan.Upgrades.Banks) / 100))
		money += val
		clan.Money += val
	}
	UpdateClan(clanName, func(clan Clan) (Clan, error) {
		clan.Money += money
		if clan.Money > MaxMoney {
			clan.Money = MaxMoney
		}
		clan.LastIncome = uint64(time.Now().Unix())
		return clan, nil
	})
	clan.LastIncome = uint64(time.Now().Unix())
	return clan
}
