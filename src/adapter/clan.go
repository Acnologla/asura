package adapter

import (
	"asura/src/entities"

	"github.com/andersfylling/disgord"
)

type UserClan struct {
	Clan   *entities.Clan
	Member *entities.ClanMember
}

type ClanAdapter interface {
	GetClan(name string, relations ...string) entities.Clan
	CreateClan(clan entities.Clan, id disgord.Snowflake) error
	GetMember(id disgord.Snowflake) (member entities.ClanMember)
	InsertMember(clan *entities.Clan, member *entities.ClanMember) error
	UpdateMember(clan *entities.Clan, member *entities.ClanMember) error
	RemoveMember(clan *entities.Clan, member disgord.Snowflake, leave bool) error
	GetUserClan(id disgord.Snowflake, relations ...string) UserClan
	UpdateClan(clan *entities.Clan, callback func(entities.Clan) entities.Clan, relations ...string) error
	SortClan(property string, limit int) []*entities.Clan
	GetClanPos(clan *entities.Clan) int
	CompleteClanMission(clan *entities.Clan, id disgord.Snowflake, xp int)
}
