package utils

import (
	"asura/src/handler"
	"context"
	"github.com/andersfylling/disgord"
	"strconv"
	"strings"
)

func GetUser(msg *disgord.Message, args []string) *disgord.User{
	if len(msg.Mentions) >0{
		return msg.Mentions[0]
	}
	if len(args) > 0 {
		converted, _ := strconv.Atoi(args[0])
		user, err := handler.Client.GetUser(context.Background(), disgord.NewSnowflake(uint64(converted)))
		if err == nil {
			return user
		}else{
			members, err := handler.Client.GetMembers(context.Background(),msg.GuildID,&disgord.GetMembersParams{
				Limit: 0,
			})
			if err == nil{
				username := strings.ToLower(strings.Join(args," "))
				for _, member := range members{
					if strings.Contains(strings.ToLower(member.Nick),username) || strings.Contains(strings.ToLower(member.User.Username),username){
						return member.User
					} 
				}
			}
		}
	}
	return msg.Author
}