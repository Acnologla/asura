package database

import (
	"context"
	"fmt"

	"github.com/andersfylling/disgord"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type Arena struct {
	Active    bool              `bun:"active"`
	Win       int               `bun:"win"`
	Lose      int               `bun:"lose"`
	LastFight disgord.Snowflake `bun:"lastFight"`
}

type Daily struct {
	Last    uint64 `bun:"last"`
	Strikes int    `bun:"strikes"`
}

type Cooldowns struct {
	Trade uint64 `bun:"trade"`
	Daily Daily  `bun:"daily"`
}

type Vip struct {
	Vip        uint64 `bun:"vip"`
	Background string `bun:"background"`
}

type Dungeon struct {
	Floor  int `bun:"floor"`
	Resets int `bun:"resets"`
}
type Status struct {
	Win  int `bun:"win"`
	Lose int `bun:"lose"`
}

type MissionProgress struct {
	Type     int `bun:"type"`
	Level    int `bun:"level"`
	Progress int `bun:"progress"`
	Adv      int `bun:"adv"`
}
type Mission struct {
	Trade    uint64             `bun:"trade"`
	Last     uint64             `bun:"last"`
	Missions []*MissionProgress `bun:"missions,array"`
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

	ID         disgord.Snowflake `bun:"id,pk"`
	Xp         int               `bun:"xp"`
	Galos      []*Rooster        `bun:"rel:has-many,join:id=userID"`
	Items      []*Item           `bun:"rel:has-many,join:id=userID"`
	Upgrades   []int             `bun:"upgrades,array"`
	Status     *Status           `bun:"status"`
	Money      int               `bun:"money"`
	Clan       string            `bun:"clan,nullzero"`
	Dungeon    *Dungeon          `bun:"dungeon"`
	Mission    *Mission          `bun:"mission"`
	Vip        *Vip              `bun:"vip"`
	TrainLimit int               `bun:"trainlimit"`
	AsuraCoin  int               `bun:"asuracoin"`
	Arena      *Arena            `bun:"arena"`
	Rank       int               `bun:"rank"`
	Cooldowns  *Cooldowns        `bun:"cooldowns"`
	Pity       int               `bun:"pity"`
}

func GetGaloDb(ID disgord.Snowflake) User {
	var user User
	db.NewSelect().Model(&user).Where("? = ?", bun.Ident("id"), ID).Scan(context.Background())
	return user
}

func InsertGaloDb(galo User) {
	_, err := db.NewInsert().Model(&galo).Exec(context.Background())
	fmt.Println(err)
}
