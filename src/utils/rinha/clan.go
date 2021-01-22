package rinha

import (
	"asura/src/database"
	"context"
	"firebase.google.com/go/db"
	"fmt"
	"github.com/andersfylling/disgord"
	"math"
	"strings"
	"time"
)

const allowedChars = "abcdefghijklmnopqrstuvwxyz123456789 -_"

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

type Clan struct {
	Xp        int          `json:"xp"`
	CreatedAt uint64       `json:"createdAt"`
	Members   []ClanMember `json:"members"`
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
	return str
}

func ClanXpToLevel(xp int) int {
	return int(math.Floor(math.Sqrt(float64(xp)/2000))) + 1

}
func ClanLevelToXp(level int) int {
	return int(math.Pow(float64(level), 2)) * 2000
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

func UpdateClan(clan string, update map[string]interface{}) {
	database.Database.NewRef("clan/"+clan).Update(context.Background(), update)
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
		text += "5 membros adicionais\n"
	}
	if level >= 4 {
		text += "2 de ouro adicional por rinha ganha\n"
	}
	if level >= 5 {
		text += "3 de ouro adicional por rinha ganha"
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

func AddClanXp(clan string, id disgord.Snowflake, xp int) {
	fn := func(tn db.TransactionNode) (interface{}, error) {
		var clan Clan
		err := tn.Unmarshal(&clan)
		if err == nil {
			clan.Xp += xp
			for i, member := range clan.Members {
				if member.ID == uint64(id) {
					member.Xp += uint(xp)
					clan.Members[i] = member
					break
				}
			}
			return clan, nil
		}
		fmt.Println(err)
		return nil, err
	}
	database.Database.NewRef("clan/"+clan).Transaction(context.Background(), fn)
}
