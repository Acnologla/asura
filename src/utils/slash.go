package utils

import (
	"strconv"

	"github.com/andersfylling/disgord"
)

func GenerateOptions(options ...*disgord.ApplicationCommandOption) []*disgord.ApplicationCommandOption {
	return options
}

func GetUser(itc *disgord.InteractionCreate, i int) *disgord.User {
	idStr := itc.Data.Options[i].Value.(string)
	id, _ := strconv.ParseUint(idStr, 10, 64)
	return itc.Data.Resolved.Users[disgord.Snowflake(id)]
}

func GetUserOrAuthor(itc *disgord.InteractionCreate, i int) *disgord.User {
	if i < len(itc.Data.Options) {
		return GetUser(itc, i)
	}
	return itc.Member.User
}
