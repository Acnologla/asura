package adapter

import (
	"asura/src/entities"

	"github.com/andersfylling/disgord"
	"github.com/google/uuid"
)

type UserAdapter interface {
	GetUser(id disgord.Snowflake, relations ...string) entities.User
	GetRoosters(id disgord.Snowflake) []entities.Rooster
	GetItems(id disgord.Snowflake) []entities.Item
	SetUser(user entities.User) error
	UpdateUser(id disgord.Snowflake, callback func(entities.User) entities.User, relations ...string) error
	InsertItem(id disgord.Snowflake, items []*entities.Item, itemID int, itemType entities.ItemType) error
	RemoveItem(items []*entities.Item, itemUUID uuid.UUID) error
	InsertRooster(rooster *entities.Rooster) error
	RemoveRooster(id uuid.UUID) error
	UpdateEquippedRooster(user entities.User, callback func(entities.Rooster) entities.Rooster) error
}
