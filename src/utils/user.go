package utils

import (
	"github.com/andersfylling/disgord"
	"strings"
)

func StringToID(id string) disgord.Snowflake{
	return disgord.ParseSnowflakeString(id)
}


func GetUser(msg *disgord.Message, args []string,session disgord.Session) *disgord.User {
	if len(msg.Mentions) > 0 {
		return msg.Mentions[0]
	}
	if len(args) > 0 {
		converted := disgord.ParseSnowflakeString(args[0])
		user, err := session.User(converted).Get()
		if err == nil {
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
