package adapter

import (
	"asura/src/entities"
	"context"

	"github.com/andersfylling/disgord"
)

type UserClan struct {
	Clan   *entities.Clan
	Member *entities.ClanMember
}

type ClanAdapter interface {
	GetClan(ctx context.Context, name string, relations ...string) entities.Clan
	CreateClan(ctx context.Context, clan entities.Clan, id disgord.Snowflake) error
	GetMember(ctx context.Context, id disgord.Snowflake) (member entities.ClanMember)
	InsertMember(ctx context.Context, clan *entities.Clan, member *entities.ClanMember) error
	UpdateMember(ctx context.Context, clan *entities.Clan, member *entities.ClanMember) error
	RemoveMember(ctx context.Context, clan *entities.Clan, member disgord.Snowflake, leave bool) error
	GetUserClan(ctx context.Context, id disgord.Snowflake, relations ...string) UserClan
	UpdateClan(ctx context.Context, clan *entities.Clan, callback func(entities.Clan) entities.Clan, relations ...string) error
	SortClan(ctx context.Context, property string, limit int) []*entities.Clan
	GetClanPos(ctx context.Context, clan *entities.Clan) int
	CompleteClanMission(ctx context.Context, clan *entities.Clan, id disgord.Snowflake, xp int)
}
