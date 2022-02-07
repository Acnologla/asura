package adapter

import (
	"asura/src/entities"
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/andersfylling/disgord"
	"github.com/uptrace/bun"
)

type ClanAdapterPsql struct {
	Db *bun.DB
}

func (adapter ClanAdapterPsql) GetClan(name string, relations ...string) (clan entities.Clan) {
	clan.Name = name
	query := adapter.Db.NewSelect().Model(&clan)
	for _, relation := range relations {
		query = query.Relation(relation)
	}
	query.Scan(context.Background())
	return
}

func (adapter ClanAdapterPsql) InsertMember(clan *entities.Clan, member *entities.ClanMember) error {
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
	err = adapter.InsertMember(&clan, &entities.ClanMember{
		ID:   id,
		Role: entities.Owner,
	})
	return err
}

func (adapter ClanAdapterPsql) GetUserClan(id disgord.Snowflake, relations ...string) entities.Clan {
	var clanMember entities.ClanMember
	clanMember.ID = id
	adapter.Db.NewSelect().Model(&clanMember).Scan(context.Background())
	return adapter.GetClan(clanMember.Clan, relations...)

}

func (adapter ClanAdapterPsql) UpdateClan(clanQuery *entities.Clan, callback func(entities.Clan) entities.Clan, relations ...string) error {

	return adapter.Db.RunInTx(context.Background(), &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		tx.ExecContext(ctx, fmt.Sprintf("SELECT pg_advisory_xact_lock(%d)", clanQuery.ID))
		clan := adapter.GetClan(clanQuery.Name, relations...)
		clan = callback(clan)
		_, err := adapter.Db.NewUpdate().Model(&clan).Where("name = ?", clan.Name).Exec(context.Background())
		return err
	})
}
