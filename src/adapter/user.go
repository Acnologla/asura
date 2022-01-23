package adapter

import (
	"asura/src/entities"

	"github.com/andersfylling/disgord"
)

type UserAdapter interface {
	GetUser(id disgord.Snowflake) entities.User
	UpdateUser(id disgord.Snowflake, callback func(entities.User) entities.User) error
}
