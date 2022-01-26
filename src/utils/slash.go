package utils

import (
	"asura/src/handler"
	"strconv"

	"github.com/andersfylling/disgord"
)

func GenerateOptions(options ...*disgord.ApplicationCommandOption) []*disgord.ApplicationCommandOption {
	return options
}

func GetUser(options []*disgord.ApplicationCommandDataOption, i int) *disgord.User {
	idStr := options[i].Value.(string)
	id, _ := strconv.ParseUint(idStr, 10, 64)
	user, _ := handler.Client.User(disgord.Snowflake(id)).Get()
	return user
}
