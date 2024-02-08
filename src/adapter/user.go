package adapter

import (
	"asura/src/entities"
	"context"

	"github.com/andersfylling/disgord"
	"github.com/google/uuid"
)

type UserAdapter interface {
	GetUser(ctx context.Context, id disgord.Snowflake, relations ...string) entities.User
	GetRoosters(ctx context.Context, id disgord.Snowflake) []entities.Rooster
	GetItems(ctx context.Context, id disgord.Snowflake) []entities.Item
	SetUser(ctx context.Context, user entities.User) error
	UpdateUser(ctx context.Context, id disgord.Snowflake, callback func(entities.User) entities.User, relations ...string) error
	UpdateMissions(ctx context.Context, id disgord.Snowflake, mission *entities.Mission, done bool)
	InsertMission(ctx context.Context, id disgord.Snowflake, mission *entities.Mission)
	InsertItem(ctx context.Context, id disgord.Snowflake, items []*entities.Item, itemID int, itemType entities.ItemType) error
	InsertManyItem(ctx context.Context, id disgord.Snowflake, items []*entities.Item, itemID int, itemType entities.ItemType, quantity int) error
	UpdateItem(ctx context.Context, user *entities.User, id uuid.UUID, callback func(entities.Item) entities.Item) error
	RemoveItem(ctx context.Context, items []*entities.Item, itemUUID uuid.UUID) error
	InsertRooster(ctx context.Context, rooster *entities.Rooster) error
	RemoveRooster(ctx context.Context, id uuid.UUID) error
	UpdateRooster(ctx context.Context, user *entities.User, id uuid.UUID, callback func(entities.Rooster) entities.Rooster) error
	UpdateEquippedRooster(ctx context.Context, user entities.User, callback func(entities.Rooster) entities.Rooster) error
	SortUsers(ctx context.Context, limit int, propertys ...string) []*entities.User
	SortUsersByRooster(ctx context.Context, limit int, propertys ...string) []*entities.User
	UpdateBp(ctx context.Context, user *entities.User, rooster *entities.Rooster) int
}
