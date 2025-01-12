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

func (role Role) ToString() string {
	return [...]string{"Membro", "Administrador", "Dono"}[role]
}

type ClanMember struct {
	bun.BaseModel `bun:"table:clanmember,alias:clanmember"`

	ID   disgord.Snowflake `bun:"id,pk"`
	Clan string            `bun:"clan"`
	Xp   int               `bun:"xp"`
	Role Role              `bun:"role"`
}

type Clan struct {
	bun.BaseModel `bun:"table:clan,alias:clan"`

	Members         []*ClanMember       `bun:"rel:has-many,join:name=clan"`
	Name            string              `bun:"name,pk"`
	ID              uint                `bun:"id"`
	Xp              int                 `bun:"xp"`
	CreatedAt       uint64              `bun:"createdat"`
	Background      string              `bun:"background"`
	Money           int                 `bun:"money"`
	MembersUpgrade  int                 `bun:"membersupgrade"`
	MissionsUpgrade int                 `bun:"missionsupgrade"`
	Mission         uint64              `bun:"mission"`
	MissionProgress int                 `bun:"missionprogress"`
	BossLife        int                 `bun:"bosslife"`
	BossDate        uint64              `bun:"bossdate"`
	BossBattles     []disgord.Snowflake `bun:"bossbattles,array"`
}
