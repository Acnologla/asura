package entities

import (
	"github.com/andersfylling/disgord"
	"github.com/uptrace/bun"
)

type Guild struct {
	bun.BaseModel `bun:"table:guilds,alias:guilds"`

	ID             disgord.Snowflake `bun:"id,pk"`
	DisableLootbox bool              `bun:"disablelootbox"`
	LootBoxChannel disgord.Snowflake `bun:"lootboxchannel"`
}
