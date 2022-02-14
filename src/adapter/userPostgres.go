package adapter

import (
	"asura/src/entities"
	"asura/src/rinha"
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

func (adapter UserAdapterPsql) GetUser(id disgord.Snowflake, relations ...string) (user entities.User) {
	query := adapter.Db.NewSelect().Model(&user)
	for _, relation := range relations {
		query = query.Relation(relation)
	}
	query.Where("id = ?", id).Scan(context.Background())
	if user.ID == 0 {
		user.ID = id
		adapter.SetUser(user)
		adapter.InsertRooster(&entities.Rooster{
			Type:   rinha.GetCommonOrRare(),
			UserID: id,
			Equip:  true,
		})
		user = adapter.GetUser(id, relations...)
	}
	return
}

func (adapter UserAdapterPsql) UpdateUser(id disgord.Snowflake, callback func(entities.User) entities.User, relations ...string) error {

	return adapter.Db.RunInTx(context.Background(), &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		tx.ExecContext(ctx, fmt.Sprintf("SELECT pg_advisory_xact_lock(%d)", id))
		user := adapter.GetUser(id, relations...)
		user = callback(user)
		_, err := adapter.Db.NewUpdate().Model(&user).Where("id = ?", user.ID).Exec(context.Background())
		return err
	})
}

func (adapter UserAdapterPsql) SetUser(user entities.User) error {
	_, err := adapter.Db.NewInsert().Model(&user).Exec(context.Background())
	return err
}

func (adapter UserAdapterPsql) GetRoosters(id disgord.Snowflake) []entities.Rooster {
	var rooster []entities.Rooster
	adapter.Db.NewInsert().Model(&rooster).Where("userid = ?", id).Exec(context.Background())
	return rooster
}

func (adapter UserAdapterPsql) GetItems(id disgord.Snowflake) []entities.Item {
	var items []entities.Item
	adapter.Db.NewInsert().Model(&items).Where("userid = ?", id).Exec(context.Background())
	return items
}

func (adapter UserAdapterPsql) InsertItem(id disgord.Snowflake, items []*entities.Item, itemID int, itemType entities.ItemType) error {
	var itemUpdate *entities.Item
	for _, item := range items {
		if item.Type == itemType && item.ItemID == itemID {
			itemUpdate = item
		}
	}
	if itemUpdate != nil {
		itemUpdate.Quantity++
		_, err := adapter.Db.NewUpdate().Model(itemUpdate).Where("id = ?", itemUpdate.ID).Exec(context.Background())

		return err
	}
	newItem := entities.Item{
		Type:     itemType,
		Quantity: 1,
		ItemID:   itemID,
		UserID:   id,
	}
	_, err := adapter.Db.NewInsert().Model(&newItem).Exec(context.Background())
	return err
}

func (adapter UserAdapterPsql) RemoveItem(items []*entities.Item, itemUUID uuid.UUID) error {
	for _, item := range items {
		if item.ID == itemUUID {
			if item.Quantity > 1 {
				item.Quantity--
				_, err := adapter.Db.NewUpdate().Model(item).Where("id = ?", item.ID).Exec(context.Background())
				return err
			} else {
				_, err := adapter.Db.NewDelete().Model(item).Where("id = ?", item.ID).Exec(context.Background())
				return err
			}
		}
	}
	return errors.New("item not found")
}

func (adapter UserAdapterPsql) InsertRooster(rooster *entities.Rooster) error {
	_, err := adapter.Db.NewInsert().Model(rooster).Exec(context.Background())
	return err
}

func (adapter UserAdapterPsql) RemoveRooster(id uuid.UUID) error {
	_, err := adapter.Db.NewDelete().Model(&entities.Rooster{}).Where("id=?", id).Exec(context.Background())
	return err
}

func (adapter UserAdapterPsql) UpdateEquippedRooster(user entities.User, callback func(entities.Rooster) entities.Rooster) error {
	galo := rinha.GetEquippedGalo(&user)
	cb := callback(*galo)
	galo = &cb
	_, err := adapter.Db.NewUpdate().Model(galo).Where("id = ?", galo.ID).Exec(context.Background())
	return err
}

func (adapter UserAdapterPsql) SortUsers(limit int, propertys ...string) (users []*entities.User) {
	query := adapter.Db.NewSelect().Model(&users)
	for _, property := range propertys {
		query = query.Order(fmt.Sprintf("%s DESC", property))
	}
	query.Limit(limit).Scan(context.Background())
	return
}

func (adapter UserAdapterPsql) SortUsersByRooster(limit int, propertys ...string) (users []*entities.User) {
	var roosters []*entities.Rooster
	query := adapter.Db.NewSelect().Model(&roosters)
	for _, property := range propertys {
		query = query.Order(fmt.Sprintf("%s DESC", property))
	}
	query.Limit(limit).Scan(context.Background())
	for _, rooster := range roosters {
		users = append(users, &entities.User{
			ID:    rooster.UserID,
			Galos: []*entities.Rooster{rooster},
		})
	}
	return
}
func (adapter UserAdapterPsql) InsertMission(id disgord.Snowflake, mission *entities.Mission) {
	mission.UserID = id
	adapter.Db.NewInsert().Model(mission).Exec(context.Background())
}

func (adapter UserAdapterPsql) UpdateMissions(id disgord.Snowflake, mission *entities.Mission, done bool) {
	if done {
		adapter.Db.NewDelete().Model(mission).WherePK().Exec(context.Background())
	} else {
		adapter.Db.NewUpdate().Model(mission).WherePK().Exec(context.Background())
	}
}
