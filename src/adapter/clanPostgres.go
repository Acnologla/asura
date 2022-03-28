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

func (adapter ClanAdapterPsql) GetClan(ctx context.Context, name string, relations ...string) (clan entities.Clan) {
	query := adapter.Db.NewSelect().Where("name = ?", name).Model(&clan)
	for _, relation := range relations {
		query = query.Relation(relation)
	}
	query.Scan(ctx)
	return
}

func (adapter ClanAdapterPsql) GetMember(ctx context.Context, id disgord.Snowflake) (member entities.ClanMember) {
	query := adapter.Db.NewSelect().Where("id = ?", id).Model(&member)
	query.Scan(ctx)
	return
}
func (adapter ClanAdapterPsql) InsertMember(ctx context.Context, clan *entities.Clan, member *entities.ClanMember) error {
	member.Clan = clan.Name
	_, err := adapter.Db.NewInsert().Model(member).Exec(ctx)
	return err
}

func (adapter ClanAdapterPsql) CreateClan(ctx context.Context, clan entities.Clan, id disgord.Snowflake) error {
	clan.CreatedAt = uint64(time.Now().Unix())
	_, err := adapter.Db.NewInsert().Model(&clan).Exec(ctx)
	if err != nil {
		return err
	}
	err = adapter.InsertMember(ctx, &clan, &entities.ClanMember{
		ID:   id,
		Role: entities.Owner,
	})
	return err
}

func (adapter ClanAdapterPsql) GetUserClan(ctx context.Context, id disgord.Snowflake, relations ...string) UserClan {
	var clanMember entities.ClanMember
	adapter.Db.NewSelect().Model(&clanMember).Where("id =?", id).Scan(ctx)
	clan := adapter.GetClan(ctx, clanMember.Clan, relations...)
	return UserClan{
		Clan:   &clan,
		Member: &clanMember,
	}
}

func (adapter ClanAdapterPsql) UpdateClan(ctx context.Context, clanQuery *entities.Clan, callback func(entities.Clan) entities.Clan, relations ...string) error {

	return adapter.Db.RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		tx.ExecContext(ctx, fmt.Sprintf("SELECT pg_advisory_xact_lock(%d)", clanQuery.ID))
		clan := adapter.GetClan(ctx, clanQuery.Name, relations...)
		clan = callback(clan)
		_, err := adapter.Db.NewUpdate().Model(&clan).Where("name = ?", clan.Name).Exec(ctx)
		return err
	})
}

func (adapter ClanAdapterPsql) RemoveMember(ctx context.Context, clan *entities.Clan, member disgord.Snowflake, leave bool) error {
	var err error
	if leave {
		m := adapter.GetMember(ctx, member)
		if m.Role == entities.Owner {
			if len(clan.Members) == 1 {
				_, err = adapter.Db.NewDelete().Model(clan).Where("name = ?", clan.Name).Exec(ctx)
			} else {
				_, err = adapter.Db.NewDelete().Model(&entities.ClanMember{}).Where("id = ? ", member).Exec(ctx)
				adapter.Db.NewUpdate().Model((*entities.ClanMember)(nil)).Set("role = ?", entities.Owner).Where("id = ?", clan.Members[1].ID).Exec(ctx)
			}
			return err
		}
	}
	_, err = adapter.Db.NewDelete().Model(&entities.ClanMember{}).Where("id = ? ", member).Exec(ctx)
	return err
}

func (adapter ClanAdapterPsql) UpdateMember(ctx context.Context, clan *entities.Clan, member *entities.ClanMember) error {
	_, err := adapter.Db.NewUpdate().Model(member).Where("id = ? ", member.ID).Exec(ctx)
	return err
}

func (adapter ClanAdapterPsql) SortClan(ctx context.Context, property string, limit int) (clans []*entities.Clan) {
	adapter.Db.NewSelect().Model(&clans).Order(fmt.Sprintf("%s DESC", property)).Limit(limit).Scan(ctx)
	return
}
func (adapter ClanAdapterPsql) GetClanPos(ctx context.Context, clan *entities.Clan) int {
	count, _ := adapter.Db.NewSelect().Model((*entities.Clan)(nil)).Where(fmt.Sprintf("xp > %d", clan.Xp)).Order("xp DESC").Count(ctx)
	return count
}

//TODO move this function to another file
func (adapter ClanAdapterPsql) CompleteClanMission(ctx context.Context, clan *entities.Clan, id disgord.Snowflake, xp int) {
	adapter.UpdateClan(ctx, clan, func(clan entities.Clan) entities.Clan {
		clan = *rinha.PopulateClanMissions(&clan)
		clan.Xp += xp
		for _, member := range clan.Members {
			if member.ID == id {
				member.Xp++
				adapter.UpdateMember(ctx, &clan, member)
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
					adapter.Db.RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
						tx.ExecContext(ctx, fmt.Sprintf("SELECT pg_advisory_xact_lock(%d)", id))
						user := &entities.User{}
						adapter.Db.NewSelect().Model(user).Relation("Galos").Where("id = ?", member.ID).Scan(ctx)
						user.Money += money
						galo := rinha.GetEquippedGalo(user)
						galo.Xp += xp
						adapter.Db.NewUpdate().Model(galo).WherePK().Exec(ctx)
						_, err := adapter.Db.NewUpdate().Model(&user).Where("id = ?", user.ID).Exec(ctx)
						return err
					})
				}
			}
		}
		return clan
	}, "Members")
}
