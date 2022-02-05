package adapter

import (
	"asura/src/entities"

	"github.com/andersfylling/disgord"
)

type UserAdapter interface {
	GetUser(id disgord.Snowflake, relations ...string) entities.User
	GetRoosters(id disgord.Snowflake) []entities.Rooster
	GetItems(id disgord.Snowflake) []entities.Item
	SetUser(user entities.User) error
	UpdateUser(id disgord.Snowflake, callback func(entities.User) entities.User, relations ...string) error
	InsertLootbox(id disgord.Snowflake, items []*entities.Item, lootType int) error
}
