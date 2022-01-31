package entities

import (
	"github.com/andersfylling/disgord"
	"github.com/uptrace/bun"
)

type Role int

const (
	Member Role = iota
	Administrator
	Owner
)

type ClanMember struct {
	bun.BaseModel `bun:"table:clanmember,alias:clanmember"`

	ID   disgord.Snowflake `bun:"id,pk"`
	Clan string            `bun:"clan"`
	Xp   int               `bun:"xp"`
	Role Role              `bun:"role"`
}

type Clan struct {
	bun.BaseModel `bun:"table:clan,alias:clan"`

	Members         []*ClanMember `bun:"rel:has-many,join:name=clan"`
	Name            string        `bun:"name, pk"`
	Xp              int           `bun:"xp"`
	CreatedAt       uint64        `bun:"createdAt"`
	Background      string        `bun:"background"`
	Money           int           `bun:"money"`
	MembersUpgrade  int           `bun:"membersUpgrade"`
	BanksUpgrade    int           `bun:"banksUpgrade"`
	MissionsUpgrade int           `bun:"missionsUpgrade"`
	Mission         uint64        `bun:"mission"`
	MissionProgress int           `bun:"missionProgress"`
}
