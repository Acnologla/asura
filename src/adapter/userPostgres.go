package adapter

import (
	"asura/src/entities"
	"asura/src/rinha"
	"asura/src/utils"
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/andersfylling/disgord"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type UserAdapterPsql struct {
	Db *bun.DB
}

func (adapter UserAdapterPsql) GetUser(ctx context.Context, id disgord.Snowflake, relations ...string) (user entities.User) {
	query := adapter.Db.NewSelect().Model(&user)
	for _, relation := range relations {
		query = query.Relation(relation)
	}
	query.Where("id = ?", id).Scan(ctx)
	if user.ID == 0 {
		user.ID = id
		adapter.SetUser(ctx, user)
		adapter.InsertRooster(ctx, &entities.Rooster{
			Type:   rinha.GetNewRooster(),
			UserID: id,
			Equip:  true,
		})
		user = adapter.GetUser(ctx, id, relations...)
	}
	return
}

func (adapter UserAdapterPsql) UpdateUser(ctx context.Context, id disgord.Snowflake, callback func(entities.User) entities.User, relations ...string) error {

	return adapter.Db.RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		tx.ExecContext(ctx, fmt.Sprintf("SELECT pg_advisory_xact_lock(%d)", id))
		user := adapter.GetUser(ctx, id, relations...)
		user = callback(user)
		_, err := adapter.Db.NewUpdate().Model(&user).Where("id = ?", user.ID).Exec(ctx)
		return err
	})
}

func (adapter UserAdapterPsql) SetUser(ctx context.Context, user entities.User) error {
	_, err := adapter.Db.NewInsert().Model(&user).Exec(ctx)
	return err
}

func (adapter UserAdapterPsql) GetRoosters(ctx context.Context, id disgord.Snowflake) []entities.Rooster {
	var rooster []entities.Rooster
	adapter.Db.NewInsert().Model(&rooster).Where("userid = ?", id).Exec(ctx)
	return rooster
}

func (adapter UserAdapterPsql) GetItems(ctx context.Context, id disgord.Snowflake) []entities.Item {
	var items []entities.Item
	adapter.Db.NewInsert().Model(&items).Where("userid = ?", id).Exec(ctx)
	return items
}

func (adapter UserAdapterPsql) InsertItemQuantity(ctx context.Context, id disgord.Snowflake, items []*entities.Item, itemID int, itemType entities.ItemType, quantity int) error {
	var itemUpdate *entities.Item
	for _, item := range items {
		if item.Type == itemType && item.ItemID == itemID {
			itemUpdate = item
		}
	}
	if itemUpdate != nil {
		itemUpdate.Quantity += quantity
		_, err := adapter.Db.NewUpdate().Model(itemUpdate).Where("id = ?", itemUpdate.ID).Exec(ctx)

		return err
	}
	newItem := entities.Item{
		Type:     itemType,
		Quantity: quantity,
		ItemID:   itemID,
		UserID:   id,
	}
	_, err := adapter.Db.NewInsert().Model(&newItem).Exec(ctx)
	return err
}

func (adapter UserAdapterPsql) InsertItem(ctx context.Context, id disgord.Snowflake, items []*entities.Item, itemID int, itemType entities.ItemType) error {
	return adapter.InsertItemQuantity(ctx, id, items, itemID, itemType, 1)
}

func (adapter UserAdapterPsql) RemoveItemQuantity(ctx context.Context, items []*entities.Item, itemUUID uuid.UUID, quantity int) error {
	for _, item := range items {
		if item.ID == itemUUID {
			if item.Quantity > quantity {
				item.Quantity -= quantity
				_, err := adapter.Db.NewUpdate().Model(item).Where("id = ?", item.ID).Exec(ctx)
				return err
			} else {
				_, err := adapter.Db.NewDelete().Model(item).Where("id = ?", item.ID).Exec(ctx)
				return err
			}
		}
	}
	return errors.New("item not found")
}

func (adapter UserAdapterPsql) RemoveItem(ctx context.Context, items []*entities.Item, itemUUID uuid.UUID) error {
	for _, item := range items {
		if item.ID == itemUUID {
			if item.Quantity > 1 {
				item.Quantity--
				_, err := adapter.Db.NewUpdate().Model(item).Where("id = ?", item.ID).Exec(ctx)
				return err
			} else {
				_, err := adapter.Db.NewDelete().Model(item).Where("id = ?", item.ID).Exec(ctx)
				return err
			}
		}
	}
	return errors.New("item not found")
}

func (adapter UserAdapterPsql) InsertRooster(ctx context.Context, rooster *entities.Rooster) error {
	if rooster.Type == 0 {
		panic("rooster type is 0")
	}
	_, err := adapter.Db.NewInsert().Model(rooster).Exec(ctx)
	return err
}

func (adapter UserAdapterPsql) RemoveRooster(ctx context.Context, id uuid.UUID) error {
	_, err := adapter.Db.NewDelete().Model(&entities.Rooster{}).Where("id=?", id).Exec(ctx)
	return err
}

func (adapter UserAdapterPsql) UpdateEquippedRooster(ctx context.Context, user entities.User, callback func(entities.Rooster) entities.Rooster) error {
	galo := rinha.GetEquippedGalo(&user)
	cb := callback(*galo)
	galo = &cb
	_, err := adapter.Db.NewUpdate().Model(galo).Where("id = ?", galo.ID).Exec(ctx)
	return err
}

func (adapter UserAdapterPsql) SortUsers(ctx context.Context, limit int, propertys ...string) (users []*entities.User) {
	query := adapter.Db.NewSelect().Model(&users)
	for _, property := range propertys {
		query = query.Order(fmt.Sprintf("%s DESC", property))
	}
	query.Limit(limit).Scan(ctx)
	return
}

func (adapter UserAdapterPsql) GetRoosterPosMultiple(ctx context.Context, user entities.User, p1, p2 string, v1, v2 int) int {
	count, _ := adapter.Db.NewSelect().Model((*entities.Rooster)(nil)).
		Where(fmt.Sprintf("%s > ? OR (%s = ? AND %s > ?)", p1, p1, p2), v1, v1, v2).
		Column("userid").
		Group("userid").
		Count(ctx)

	return count
}

func (adapter UserAdapterPsql) GetPosMultiple(ctx context.Context, user entities.User, p1, p2 string, v1, v2 int) int {
	count, _ := adapter.Db.NewSelect().Model((*entities.User)(nil)).
		Where(fmt.Sprintf("%s > ? OR (%s = ? AND %s > ?)", p1, p1, p2), v1, v1, v2).
		Count(ctx)

	return count
}

func (adapter UserAdapterPsql) GetPos(ctx context.Context, user entities.User, property string, data func(u *entities.User) int) int {
	count, _ := adapter.Db.NewSelect().Model((*entities.User)(nil)).Where(fmt.Sprintf("%s > %d", property, data(&user))).Count(ctx)
	return count
}

func (adapter UserAdapterPsql) SortUsersByRooster(ctx context.Context, limit int, propertys ...string) (users []*entities.User) {
	var roosters []*entities.Rooster
	query := adapter.Db.NewSelect().Model(&roosters)
	for _, property := range propertys {
		query = query.Order(fmt.Sprintf("%s DESC", property))
	}
	query.Limit(limit).Scan(ctx)
	for _, rooster := range roosters {
		users = append(users, &entities.User{
			ID:    rooster.UserID,
			Galos: []*entities.Rooster{rooster},
		})
	}
	return
}
func (adapter UserAdapterPsql) InsertMission(ctx context.Context, id disgord.Snowflake, mission *entities.Mission) {
	mission.UserID = id
	adapter.Db.NewInsert().Model(mission).Exec(ctx)
}

func (adapter UserAdapterPsql) UpdateMissions(ctx context.Context, id disgord.Snowflake, mission *entities.Mission, done bool) {
	if done {
		adapter.Db.NewDelete().Model(mission).WherePK().Exec(ctx)
	} else {
		adapter.Db.NewUpdate().Model(mission).WherePK().Exec(ctx)
	}
}

func (adapter UserAdapterPsql) AddTrialWin(ctx context.Context, trial *entities.Trial) error {
	trial.Win++
	_, err := adapter.Db.NewUpdate().Model(trial).WherePK().Exec(ctx)
	return err
}

func (adapter UserAdapterPsql) UpdateRooster(ctx context.Context, user *entities.User, id uuid.UUID, callback func(entities.Rooster) entities.Rooster) error {
	var galo *entities.Rooster
	for _, gal := range user.Galos {
		if gal.ID == id {
			galo = gal
			continue
		}
	}
	if galo == nil {
		return errors.New("galo not found")
	}
	cb := callback(*galo)
	galo = &cb
	_, err := adapter.Db.NewUpdate().Model(galo).Where("id = ?", galo.ID).Exec(ctx)
	return err
}

func (adapter UserAdapterPsql) InsertTransaction(ctx context.Context, id disgord.Snowflake, transaction *entities.Transaction) error {
	transaction.UserID = id
	_, err := adapter.Db.NewInsert().Model(transaction).Exec(ctx)
	return err
}

func (adapter UserAdapterPsql) InsertTrial(ctx context.Context, id disgord.Snowflake, trial *entities.Trial) error {
	trial.UserID = id
	_, err := adapter.Db.NewInsert().Model(trial).Exec(ctx)
	return err
}

func (adapter UserAdapterPsql) UpdateItem(ctx context.Context, user *entities.User, id uuid.UUID, callback func(entities.Item) entities.Item) error {
	var item *entities.Item
	for _, it := range user.Items {
		if it.ID == id {
			item = it
			continue
		}
	}
	if item == nil {
		return errors.New("item not found")
	}
	cb := callback(*item)
	item = &cb
	_, err := adapter.Db.NewUpdate().Model(item).Where("id = ?", item.ID).Exec(ctx)
	return err
}

func (adapter UserAdapterPsql) UpdateBp(ctx context.Context, user *entities.User, rooster *entities.Rooster) int {
	level := rinha.CalcBPLevel(user.BattlePass)
	isVip := rinha.IsVip(user)
	if !isVip && level >= len(rinha.BattlePass)/2 {
		return 0
	}

	xpOb := utils.RandInt(6) + 2

	if isVip {
		xpOb += 2
	}

	user.BattlePass += xpOb
	if level != rinha.CalcBPLevel(user.BattlePass) {
		level++
		if len(rinha.BattlePass) > level {
			currentLevel := rinha.BattlePass[level]
			switch currentLevel.Type {
			case rinha.BattlePassItem:
				adapter.InsertItem(ctx, user.ID, user.Items, currentLevel.Value, currentLevel.ItemType)
			case rinha.BattlePassMoney:
				user.Money += currentLevel.Value
			case rinha.BattlePassXp:
				divider := 1
				if rooster.Resets > 1 {
					divider = 2
				}
				rooster.Xp += currentLevel.Value / divider
			case rinha.BattlePassCoins:
				user.AsuraCoin += currentLevel.Value
			}

		}
	}
	return xpOb
}
