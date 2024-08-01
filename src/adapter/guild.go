package adapter

import (
	"asura/src/entities"
	"context"

	"github.com/andersfylling/disgord"
)

type GuildAdapter interface {
	GetGuild(ctx context.Context, id disgord.Snowflake) entities.Guild
	SetGuild(ctx context.Context, guild entities.Guild) error
	UpdateGuild(ctx context.Context, guild entities.Guild) error
}
