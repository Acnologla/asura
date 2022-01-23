package adapter

import (
	"asura/src/entities"

	"github.com/andersfylling/disgord"
)

type UserAdapter interface {
	GetUser(id disgord.Snowflake) entities.User
	GetRoosters(id disgord.Snowflake) []entities.Rooster
	GetItems(id disgord.Snowflake) []entities.Item
	SetUser(user entities.User) error
	UpdateUser(id disgord.Snowflake, callback func(entities.User) entities.User) error
}
