package interaction

import "github.com/andersfylling/disgord"

type User struct {
	ID       disgord.Snowflake `json:"id"`
	Username string            `json:"username"`
}

type Member struct {
	User *User `json:"user"`
}
