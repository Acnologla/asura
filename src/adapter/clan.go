package adapter

import (
	"asura/src/entities"

	"github.com/andersfylling/disgord"
)

type ClanAdapter interface {
	GetClan(name string, relations ...string) entities.Clan
	CreateClan(clan entities.Clan, id disgord.Snowflake) error
	InsertMember(clan entities.Clan, member *entities.ClanMember) error
	GetUserClan(id disgord.Snowflake, relations ...string) entities.Clan
}
