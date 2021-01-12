package utils

import (
	"asura/src/handler"
	"context"
	"github.com/andersfylling/disgord"
	"strconv"
	"strings"
)

func StringToID(id string) disgord.Snowflake {
	converted, err := strconv.Atoi(id)
	if err == nil && converted >= 0 {
		snowflake := disgord.ParseSnowflakeString(id)
		return snowflake
	}
	return 0
}

func GetUser(msg *disgord.Message, args []string, session disgord.Session) *disgord.User {
	if len(msg.Mentions) > 0 {
		return msg.Mentions[0]
	}
	if len(args) > 0 {
		converted := StringToID(args[0])
		user, err := session.User(converted).Get()
		if err == nil && converted != 0  {
			return user
		} else {
			members, err := session.Guild(msg.GuildID).GetMembers(&disgord.GetMembersParams{
				Limit: 0,
			})
			if err == nil {
				username := strings.ToLower(strings.Join(args, " "))
				for _, member := range members {
					if strings.Contains(strings.ToLower(member.Nick), username) || strings.Contains(strings.ToLower(member.User.Username), username) {
						return member.User
					}
				}
			}
		}
	}
	return msg.Author
}

func Confirm(confirmMsg *disgord.Message, id disgord.Snowflake, callback func()) {
	Try(func() error {
		return confirmMsg.React(context.Background(), handler.Client, "✅")
	}, 5)
	done := false
	handler.RegisterHandler(confirmMsg, func(removed bool, emoji disgord.Emoji, u disgord.Snowflake) {
		if emoji.Name == "✅" && !removed && u == id && !done {
			done = true
			handler.DeleteHandler(confirmMsg)
			go handler.Client.Channel(confirmMsg.ChannelID).Message(confirmMsg.ID).Delete()
			callback()
		}
	}, 120)
}
