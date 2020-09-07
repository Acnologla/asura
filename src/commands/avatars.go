package commands

import (
	"asura/src/database"
	"asura/src/handler"
	"asura/src/utils"
	"context"
	"fmt"
	"github.com/andersfylling/disgord"
	"strconv"
)

func init() {
	handler.Register(handler.Command{
		Aliases:   []string{"oldavatars", "avatars"},
		Run:       runAvatars,
		Available: true,
		Cooldown:  5,
		Usage:     "j!oldavatars <usuario>",
		Help:      "Veja os avatares antigos de alguem",
	})
}

func runAvatars(session disgord.Session, msg *disgord.Message, args []string) {
	user := utils.GetUser(msg, args, session)
	ctx := context.Background()
	var userinfo database.User
	var private bool
	id := strconv.FormatUint(uint64(user.ID), 10)
	database.Database.NewRef("users/"+id).Get(ctx, &userinfo)
	database.Database.NewRef("private/"+id).Get(ctx, &private)
	if private {
		msg.Reply(ctx, session, msg.Author.Mention()+", O historico de avatar desse usuario é privado")
		return
	}
	if len(userinfo.Avatars) == 0 {
		msg.Reply(ctx, session, msg.Author.Mention()+", Não tenho o historico de avatares desse usuario")
		return
	}
	count := 0
	message, err := msg.Reply(context.Background(), session, &disgord.CreateMessageParams{
		Embed: &disgord.Embed{
			Color: 65535,
			Title: fmt.Sprintf("Avatar numero %d", count+1),
			Image: &disgord.EmbedImage{
				URL: userinfo.Avatars[0],
			},
		},
	})
	if err != nil {
		return
	}
	utils.Try(func() error {
		return message.React(context.Background(), session, "⬅️")
	}, 3)
	utils.Try(func() error {
		return message.React(context.Background(), session, "➡️")
	}, 3)
	handler.RegisterHandler(message, func(removed bool, emoji disgord.Emoji, u disgord.Snowflake) {
		if !removed && msg.Author.ID == u {
			msgUpdater := handler.Client.UpdateMessage(ctx, msg.ChannelID, message.ID)
			if emoji.Name == "⬅️" {
				if count == 0 {
					count = len(userinfo.Avatars) - 1
				} else {
					count--
				}
			}
			if emoji.Name == "➡️" {
				if count == len(userinfo.Avatars)-1 {
					count = 0
				} else {
					count++
				}
			}
			if emoji.Name == "➡️" || emoji.Name == "⬅️" {
				msgUpdater.SetEmbed(&disgord.Embed{
					Color: 65535,
					Title: fmt.Sprintf("Avatar numero %d", count+1),
					Image: &disgord.EmbedImage{
						URL: userinfo.Avatars[count],
					},
				})
				msgUpdater.Execute()
				handler.Client.DeleteUserReaction(ctx, msg.ChannelID, message.ID, u, emoji.Name)
			}
		}
	}, 60*10)
}
