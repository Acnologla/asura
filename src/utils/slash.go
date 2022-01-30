package utils

import (
	"asura/src/handler"
	"strconv"

	"github.com/andersfylling/disgord"
)

func GenerateOptions(options ...*disgord.ApplicationCommandOption) []*disgord.ApplicationCommandOption {
	return options
}

func GetUser(itc *disgord.InteractionCreate, i int) *disgord.User {
	opt := itc.Data.Options[i]
	idStr := opt.Value.(string)
	id, _ := strconv.ParseUint(idStr, 10, 64)
	if opt.Type == disgord.OptionTypeString {
		u, err := handler.Client.User(disgord.Snowflake(id)).Get()
		if err != nil {
			return itc.Member.User
		}
		return u
	}
	return itc.Data.Resolved.Users[disgord.Snowflake(id)]
}

func GetUserOrAuthor(itc *disgord.InteractionCreate, i int) *disgord.User {
	if i < len(itc.Data.Options) {
		return GetUser(itc, i)
	}
	return itc.Member.User
}
