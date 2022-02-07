package adapter

import (
	"asura/src/entities"
	"context"
	"time"

	"github.com/andersfylling/disgord"
	"github.com/uptrace/bun"
)

type ClanAdapterPsql struct {
	Db *bun.DB
}

func (adapter ClanAdapterPsql) GetClan(name string) (clan entities.Clan) {
	clan.Name = name
	query := adapter.Db.NewSelect().Model(&clan)
	query.Scan(context.Background())
	return
}

func (adapter ClanAdapterPsql) InsertMember(clan entities.Clan, member *entities.ClanMember) error {
	member.Clan = clan.Name
	_, err := adapter.Db.NewInsert().Model(member).Exec(context.Background())
	return err
}

func (adapter ClanAdapterPsql) CreateClan(clan entities.Clan, id disgord.Snowflake) error {
	clan.CreatedAt = uint64(time.Now().Unix())
	_, err := adapter.Db.NewInsert().Model(&clan).Exec(context.Background())
	if err != nil {
		return err
	}
	err = adapter.InsertMember(clan, &entities.ClanMember{
		ID:   id,
		Role: entities.Owner,
	})
	return err
}
