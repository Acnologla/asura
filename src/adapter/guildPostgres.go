package adapter

import (
	"asura/src/entities"
	"context"

	"github.com/andersfylling/disgord"
	"github.com/uptrace/bun"
)

type GuildAdapterPsql struct {
	Db *bun.DB
}

func (adapter GuildAdapterPsql) GetGuild(ctx context.Context, id disgord.Snowflake) (guild entities.Guild) {
	adapter.Db.NewSelect().Model(&guild).Where("id = ?", id).Scan(ctx)
	if guild.ID == 0 {
		guild.ID = id
		adapter.SetGuild(ctx, guild)
	}
	return
}

func (adapter GuildAdapterPsql) SetGuild(ctx context.Context, guild entities.Guild) error {
	_, err := adapter.Db.NewInsert().Model(&guild).Exec(ctx)
	return err
}

func (adapter GuildAdapterPsql) UpdateGuild(ctx context.Context, guild entities.Guild) error {
	_, err := adapter.Db.NewUpdate().Model(&guild).Where("id = ?", guild.ID).Exec(ctx)
	return err
}
