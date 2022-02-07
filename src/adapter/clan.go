package adapter

import (
	"asura/src/entities"

	"github.com/andersfylling/disgord"
)

type ClanAdapter interface {
	GetClan(name string) entities.Clan
	CreateClan(clan entities.Clan, id disgord.Snowflake) error
	InsertMember(clan entities.Clan, member *entities.ClanMember) error
}
