package utils

import (
	"asura/src/handler"
	"image"
	"strconv"
	"strings"

	"github.com/andersfylling/disgord"
)

func IsNumber(text string) bool {
	_, err := strconv.Atoi(text)
	return err == nil
}

func StringToID(id string) disgord.Snowflake {
	converted, err := strconv.Atoi(id)
	if err == nil && converted >= 0 {
		snowflake := disgord.ParseSnowflakeString(id)
		return snowflake
	}
	return 0
}

func DownloadAvatar(id disgord.Snowflake, size int, gif bool) (image.Image, error) {
	user, _ := handler.Client.User(id).WithFlags(disgord.IgnoreCache).Get()
	avatar, _ := user.AvatarURL(size, gif)
	avatar = strings.Replace(avatar, ".webp", ".png", 1)
	return DownloadImage(avatar)
}

func GetAllStringArgs(args []*disgord.ApplicationCommandDataOption) []string {
	arr := []string{}
	for i := range args {
		str, _ := handler.GetStringArg(args, i)
		arr = append(arr, str)
	}
	return arr
}

func GetUser(msg *disgord.Message, args []*disgord.ApplicationCommandDataOption, session disgord.Session) *disgord.User {
	if len(msg.Mentions) > 0 {
		return msg.Mentions[0]
	}
	if len(args) > 0 {
		if args[0].Type == disgord.USER {
			user, _ := args[0].Value.(*disgord.User)
			return user
		}
		stringArg, _ := handler.GetStringArg(args, 0)
		converted := StringToID(stringArg)
		user, err := session.User(converted).Get()
		if err == nil && converted != 0 {
			return user
		} else {
			members, err := session.Guild(msg.GuildID).GetMembers(&disgord.GetMembersParams{
				Limit: 0,
			})
			if err == nil {
				username := strings.ToLower(strings.Join(GetAllStringArgs(args), " "))
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

func Confirm(title string, channel, id disgord.Snowflake, callback func()) {
	msg, err := handler.Client.Channel(channel).CreateMessage(&disgord.CreateMessageParams{
		Embed: &disgord.Embed{
			Title: title,
			Color: 65535,
		},
		Components: []*disgord.MessageComponent{
			{
				Type: disgord.MessageComponentActionRow,
				Components: []*disgord.MessageComponent{
					{
						Type:     disgord.MessageComponentButton,
						Label:    "Aceitar",
						Style:    disgord.Success,
						CustomID: "yes",
					},
					{
						Type:     disgord.MessageComponentButton,
						Label:    "Recusar",
						Style:    disgord.Danger,
						CustomID: "no",
					},
				},
			},
		},
	})
	if err != nil {
		return
	}
	done := false
	handler.RegisterBHandler(msg.ID, func(interaction *disgord.InteractionCreate) {
		if id == interaction.Member.User.ID && !done {
			done = true
			handler.DeleteBHandler(msg.ID)
			go handler.Client.Channel(msg.ChannelID).Message(msg.ID).Delete()
			if interaction.Data.CustomID == "yes" {
				callback()
			}
		}
	}, 120)
}
