package adapter

import (
	"asura/src/entities"
	"asura/src/rinha"
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
	query := adapter.Db.NewSelect().Where("name = ?", name).Model(&clan)
	for _, relation := range relations {
		query = query.Relation(relation)
	}
	query.Scan(context.Background())
	return
}

func (adapter ClanAdapterPsql) GetMember(id disgord.Snowflake) (member entities.ClanMember) {
	query := adapter.Db.NewSelect().Where("id = ?", id).Model(&member)
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

func (adapter ClanAdapterPsql) GetUserClan(id disgord.Snowflake, relations ...string) UserClan {
	var clanMember entities.ClanMember
	adapter.Db.NewSelect().Model(&clanMember).Where("id =?", id).Scan(context.Background())
	clan := adapter.GetClan(clanMember.Clan, relations...)
	return UserClan{
		Clan:   &clan,
		Member: &clanMember,
	}
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

func (adapter ClanAdapterPsql) RemoveMember(clan *entities.Clan, member disgord.Snowflake, leave bool) error {
	var err error
	if leave {
		m := adapter.GetMember(member)
		if m.Role == entities.Owner {
			if len(clan.Members) == 1 {
				_, err = adapter.Db.NewDelete().Model(clan).Where("name = ?", clan.Name).Exec(context.Background())
			} else {
				_, err = adapter.Db.NewDelete().Model(&entities.ClanMember{}).Where("id = ? ", member).Exec(context.Background())
				adapter.Db.NewUpdate().Model((*entities.ClanMember)(nil)).Set("role = ?", entities.Owner).Where("id = ?", clan.Members[1].ID).Exec(context.Background())
			}
		}
	} else {
		_, err = adapter.Db.NewDelete().Model(&entities.ClanMember{}).Where("id = ? ", member).Exec(context.Background())
	}
	return err
}

func (adapter ClanAdapterPsql) UpdateMember(clan *entities.Clan, member *entities.ClanMember) error {
	_, err := adapter.Db.NewUpdate().Model(member).Where("id = ? ", member.ID).Exec(context.Background())
	return err
}

func (adapter ClanAdapterPsql) SortClan(property string, limit int) (clans []*entities.Clan) {
	adapter.Db.NewSelect().Model(&clans).Order(fmt.Sprintf("%s DESC", property)).Limit(limit).Scan(context.Background())
	return
}
func (adapter ClanAdapterPsql) GetClanPos(clan *entities.Clan) int {
	count, _ := adapter.Db.NewSelect().Model((*entities.Clan)(nil)).Where(fmt.Sprintf("xp > %d", clan.Xp)).Order("xp DESC").Count(context.Background())
	return count
}

//TODO move this function to another file
func (adapter ClanAdapterPsql) CompleteClanMission(clan *entities.Clan, id disgord.Snowflake, xp int) {
	adapter.UpdateClan(clan, func(clan entities.Clan) entities.Clan {
		clan = *rinha.PopulateClanMissions(&clan)
		clan.Xp += xp
		for _, member := range clan.Members {
			if member.ID == id {
				member.Xp++
				adapter.UpdateMember(&clan, member)
				break
			}
		}
		money, xp := rinha.CalcMissionPrize(&clan)
		done := clan.MissionProgress >= rinha.MissionN
		if !done {
			clan.MissionProgress++
			if clan.MissionProgress >= rinha.MissionN {
				for _, member := range clan.Members {
					adapter.Db.NewUpdate()
					adapter.Db.RunInTx(context.Background(), &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
						tx.ExecContext(ctx, fmt.Sprintf("SELECT pg_advisory_xact_lock(%d)", id))
						user := &entities.User{}
						adapter.Db.NewSelect().Model(user).Relation("Galos").Where("id = ?", member.ID).Scan(ctx)
						user.Money += money
						galo := rinha.GetEquippedGalo(user)
						galo.Xp += xp
						adapter.Db.NewUpdate().Model(galo).WherePK().Exec(ctx)
						_, err := adapter.Db.NewUpdate().Model(&user).Where("id = ?", user.ID).Exec(context.Background())
						return err
					})
				}
			}
		}
		return clan
	}, "Members")
}
