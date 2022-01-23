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
