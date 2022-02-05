package adapter

import (
	"asura/src/entities"
	"context"
	"database/sql"
	"fmt"

	"github.com/andersfylling/disgord"
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
	}
	return user
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

func (adapter UserAdapterPsql) InsertLootbox(id disgord.Snowflake, items []*entities.Item, lootType int) error {
	var itemUpdate *entities.Item
	for _, item := range items {
		if item.Type == entities.LootboxType {
			if item.ItemID == lootType {
				itemUpdate = item
			}
		}
	}
	if itemUpdate != nil {
		itemUpdate.Quantity++
		_, err := adapter.Db.NewUpdate().Model(itemUpdate).Where("id = ?", itemUpdate.ID).Exec(context.Background())

		return err
	}
	newItem := entities.Item{
		Type:     entities.LootboxType,
		Quantity: 1,
		ItemID:   lootType,
		UserID:   id,
	}
	_, err := adapter.Db.NewInsert().Model(&newItem).Exec(context.Background())
	return err
}
