package utils

import (
	"asura/src/handler"
	"context"
	"image"
	"strconv"
	"strings"

	"github.com/andersfylling/disgord"
)

type IDUsername struct {
	ID       disgord.Snowflake
	Username string
}

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
	user, _ := handler.Client.User(id).Get(disgord.IgnoreCache)
	avatar, _ := user.AvatarURL(size, gif)
	avatar = strings.Replace(avatar, ".webp", ".png", 1)
	return DownloadImage(avatar)
}

func GetUser(msg *disgord.Message, args []string, session disgord.Session) *disgord.User {
	if len(msg.Mentions) > 0 {
		return msg.Mentions[0]
	}
	if len(args) > 0 {
		converted := StringToID(args[0])
		user, err := session.User(converted).Get()
		if err == nil && converted != 0 {
			return user
		} else if !IsNumber(args[0]) {
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

func IndexId(id disgord.Snowflake, ids []IDUsername) int {
	for i, id2 := range ids {
		if id == id2.ID {
			return i
		}
	}
	return -1
}

func SpliceID(i int, ids []IDUsername) []IDUsername {
	return append(ids[:i], ids[i+1:]...)
}

func IDUsernameToString(ids []IDUsername) string {
	var str string
	for _, id := range ids {
		str += id.Username + "\n"
	}
	return str
}

func IncludesID(id disgord.Snowflake, ids []IDUsername) bool {
	for _, id2 := range ids {
		if id == id2.ID {
			return true
		}
	}
	return false
}

func ConfirmArray(title string, channel disgord.Snowflake, ids []IDUsername, callback func()) {
	msg, err := handler.Client.Channel(channel).CreateMessage(&disgord.CreateMessageParams{
		Embed: &disgord.Embed{
			Title:       title,
			Color:       65535,
			Description: "Usuarios que não aceitaram: \n" + IDUsernameToString(ids),
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
				},
			},
		},
	})
	if err != nil {
		return
	}
	done := false
	handler.RegisterBHandler(msg, func(interaction *disgord.InteractionCreate) {
		id := interaction.Member.User.ID
		if IncludesID(id, ids) && !done {
			if interaction.Data.CustomID == "yes" {
				index := IndexId(id, ids)
				ids = SpliceID(index, ids)
				if len(ids) == 0 {
					done = true
					handler.DeleteBHandler(msg)
					go handler.Client.Channel(msg.ChannelID).Message(msg.ID).Delete()
					callback()
				} else {
					handler.Client.SendInteractionResponse(context.Background(), interaction, &disgord.InteractionResponse{
						Type: disgord.UpdateMessage,
						Data: &disgord.InteractionApplicationCommandCallbackData{
							Embeds: []*disgord.Embed{
								{
									Title:       title,
									Color:       65535,
									Description: "Usuarios que não aceitaram: \n" + IDUsernameToString(ids),
								},
							},
						},
					})
				}
			}
		}
	}, 120)
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
	handler.RegisterBHandler(msg, func(interaction *disgord.InteractionCreate) {
		if id == interaction.Member.User.ID && !done {
			done = true
			handler.DeleteBHandler(msg)
			go handler.Client.Channel(msg.ChannelID).Message(msg.ID).Delete()
			if interaction.Data.CustomID == "yes" {
				callback()
			}
		}
	}, 120)
}

func UsersToIDUsername(users []*disgord.User) []IDUsername {
	var ids []IDUsername
	for _, user := range users {
		ids = append(ids, IDUsername{
			ID:       user.ID,
			Username: user.Username,
		})
	}
	return ids
}
