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

func (adapter UserAdapterPsql) GetUser(id disgord.Snowflake) (user entities.User) {
	adapter.Db.NewSelect().Model(&user).Where("? = ?", bun.Ident("id"), id).Scan(context.Background())
	return user
}

func (adapter UserAdapterPsql) UpdateUser(id disgord.Snowflake, callback func(entities.User) entities.User) error {

	return adapter.Db.RunInTx(context.Background(), &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		tx.ExecContext(ctx, fmt.Sprintf("SELECT pg_advisory_xact_lock(%d)", id))
		user := adapter.GetUser(id)
		user = callback(user)
		_, err := adapter.Db.NewUpdate().Model(&user).Where("? = ?", bun.Ident("id"), user.ID).Exec(context.Background())
		return err
	})
}

func (adapter UserAdapterPsql) SetUser(user entities.User) error {
	_, err := adapter.Db.NewInsert().Model(&user).Exec(context.Background())
	return err
}

func (adapter UserAdapterPsql) GetRoosters(id disgord.Snowflake) []entities.Rooster {
	var rooster []entities.Rooster
	adapter.Db.NewInsert().Model(&rooster).Where("? = ?", bun.Ident("userid"), id).Exec(context.Background())
	return rooster
}

func (adapter UserAdapterPsql) GetItems(id disgord.Snowflake) []entities.Item {
	var items []entities.Item
	adapter.Db.NewInsert().Model(&items).Where("? = ?", bun.Ident("userid"), id).Exec(context.Background())
	return items
}
