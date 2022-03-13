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
	UpdateMissions(id disgord.Snowflake, mission *entities.Mission, done bool)
	InsertMission(id disgord.Snowflake, mission *entities.Mission)
	InsertItem(id disgord.Snowflake, items []*entities.Item, itemID int, itemType entities.ItemType) error
	UpdateItem(user *entities.User, id uuid.UUID, callback func(entities.Item) entities.Item) error
	RemoveItem(items []*entities.Item, itemUUID uuid.UUID) error
	InsertRooster(rooster *entities.Rooster) error
	RemoveRooster(id uuid.UUID) error
	UpdateRooster(user *entities.User, id uuid.UUID, callback func(entities.Rooster) entities.Rooster) error
	UpdateEquippedRooster(user entities.User, callback func(entities.Rooster) entities.Rooster) error
	SortUsers(limit int, propertys ...string) []*entities.User
	SortUsersByRooster(limit int, propertys ...string) []*entities.User
}
