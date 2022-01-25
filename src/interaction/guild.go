package interaction

import (
	"fmt"
	"strings"

	"github.com/andersfylling/disgord"
)

type User struct {
	ID            disgord.Snowflake `json:"id"`
	Username      string            `json:"username"`
	Discriminator int               `json:"discriminator"`
	Avatar        string            `json:"avatar"`
}

func (u *User) AvatarURL(size int) string {
	if u.Avatar == "" {
		return fmt.Sprintf("https://cdn.discordapp.com/embed/avatars/%d.png?size=%d", u.Discriminator%5, size)
	} else if strings.HasPrefix(u.Avatar, "a_") {
		return fmt.Sprintf("https://cdn.discordapp.com/avatars/%d/%s.gif?size=%d", u.ID, u.Avatar, size)
	} else {
		return fmt.Sprintf("https://cdn.discordapp.com/avatars/%d/%s.webp?size=%d", u.ID, u.Avatar, size)
	}
}

type Member struct {
	User *User `json:"user"`
}
