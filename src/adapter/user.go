package adpter

import (
	"asura/src/database"
	"asura/src/entities"
	"context"

	"github.com/andersfylling/disgord"
	"github.com/uptrace/bun"
)

type UserAdapter interface {
	GetUser(id disgord.Snowflake) entities.User
	GetEquippedGalo(id disgord.Snowflake) entities.Rooster
	GetGalo(id disgord.Snowflake) []entities.Rooster

	SetUser(user entities.User) error
	UpdateUser(user entities.User) error
}

type userAdapterImpl struct {
	db *bun.DB
}

func (adpter userAdapterImpl) GetUser(id disgord.Snowflake) (user entities.User) {
	adpter.db.NewSelect().Model(&user).Where("? = ?", bun.Ident("id"), id).Scan(context.Background())
	return user
}

var UserAdapterImpl = &userAdapterImpl{
	db: database.Database,
}
