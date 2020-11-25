package rinha

import (
	"asura/src/database"
	"context"
	"github.com/andersfylling/disgord"
	"math"
	"strings"
	"firebase.google.com/go/db"
	"time"
	"fmt"
)

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
	return strings.TrimSpace(strings.ToLower(text))
}

func ClanXpToLevel(xp int) int {
	return int(math.Floor(math.Sqrt(float64(xp)/1000))) + 1

}
func ClanLevelToXp(level int) int {
	return int(math.Pow(float64(level), 2)) * 1000
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


func GetBenefits(xp int) (text string){
	level := ClanXpToLevel(xp)
	text = "10% de xp adicional por rinha ganha\n"
	if level >= 2{
		text += "10% de xp adicional por rinha ganha\n"
	}
	if level >= 3{
		text += "2 de ouro adicional por rinha ganha\n"
	}
	if level >= 4{
		text += "3 de ouro adicionadl por rinha ganha"
	}
	return
}


func AddClanXp(clan string, xp int){
	fn := func(tn db.TransactionNode) (interface{}, error) {
		var clan Clan
		err := tn.Unmarshal(&clan)
		if err == nil {
			clan.Xp += xp
			return clan, nil
		}
		fmt.Println(err)
		return nil, err
	}
	database.Database.NewRef("clan/"+clan).Transaction(context.Background(), fn)
}
