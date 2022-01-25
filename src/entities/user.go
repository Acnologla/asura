package entities

import (
	"github.com/andersfylling/disgord"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type Mission struct {
	bun.BaseModel `bun:"table:mission,alias:mission"`
	UserID        disgord.Snowflake `bun:"userID, pk"`
	Type          int               `bun:"type"`
	Level         int               `bun:"level"`
	Progress      int               `bun:"progress"`
	Adv           int               `bun:"adv"`
}

type Item struct {
	bun.BaseModel `bun:"table:item,alias:item"`
	UserID        disgord.Snowflake `bun:"userID"`
	Quatity       int               `bun:"quatity"`
	ItemID        int               `bun:"itemID"`
	Equip         bool              `bun:"equip"`
	Type          int               `bun:"type"`
}
type Rooster struct {
	bun.BaseModel `bun:"table:rooster,alias:galo"`

	ID       uuid.UUID         `bun:"id,pk"`
	UserID   disgord.Snowflake `bun:"userID"`
	Name     string            `bun:"name"`
	Resets   int               `bun:"resets"`
	Equip    bool              `bun:"equip"`
	Xp       int               `bun:"xp"`
	Type     int               `bun:"type"`
	Equipped []int             `bun:"equipped,array"`
}

type User struct {
	bun.BaseModel `bun:"table:users,alias:u"`

	ID             disgord.Snowflake `bun:"id,pk"`
	Xp             int               `bun:"xp"`
	Galos          []*Rooster        `bun:"rel:has-many,join:id=userID"`
	Items          []*Item           `bun:"rel:has-many,join:id=userID"`
	Upgrades       []int             `bun:"upgrades,array"`
	Win            int               `bun:"win"`
	Lose           int               `bun:"lose"`
	Money          int               `bun:"money"`
	Clan           string            `bun:"clan,nullzero"`
	Dungeon        int               `bun:"dungeon"`
	DungeonReset   int               `bun:"dungeonReset"`
	TradeMission   uint64            `bun:"tradeMission"`
	LastMission    uint64            `bun:"lastMission"`
	Missions       []*Mission        `bun:"rel:has-many,join:id=userID"`
	Vip            uint64            `bun:"vip"`
	VipBackground  string            `bun:"vipBackground"`
	TrainLimit     int               `bun:"trainlimit"`
	AsuraCoin      int               `bun:"asuracoin"`
	ArenaActive    bool              `bun:"arenaActive"`
	ArenaWin       int               `bun:"arenaWin"`
	ArenaLose      int               `bun:"arenaLose"`
	ArenaLastFight disgord.Snowflake `bun:"arenaLastFight"`
	Rank           int               `bun:"rank"`
	TradeItem      uint64            `bun:"tradeItem"`
	Daily          uint64            `bun:"daily"`
	DailyStrikes   int               `bun:"dailyStrikes"`
	Pity           int               `bun:"pity"`
}
