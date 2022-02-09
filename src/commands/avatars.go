package commands

import (
	"asura/src/firebase"
	"asura/src/handler"
	"asura/src/translation"
	"asura/src/utils"
	"context"
	"fmt"
	"strconv"

	"github.com/andersfylling/disgord"
)

func init() {
	handler.RegisterCommand(handler.Command{
		Name:        "avatars",
		Description: translation.T("AvatarsHelp", "pt"),
		Run:         runAvatars,
		Cooldown:    10,
		Options: utils.GenerateOptions(&disgord.ApplicationCommandOption{
			Name:        "user",
			Type:        disgord.OptionTypeUser,
			Description: "user avatar",
			Required:    true,
		}),
	})
}

func runAvatars(itc *disgord.InteractionCreate) *disgord.InteractionResponse {
	user := utils.GetUser(itc, 0)
	ctx := context.Background()
	var userinfo firebase.User
	var private bool
	id := strconv.FormatUint(uint64(user.ID), 10)
	firebase.Database.NewRef("users/"+id).Get(ctx, &userinfo)
	firebase.Database.NewRef("private/"+id).Get(ctx, &private)
	if private {
		return &disgord.InteractionResponse{
			Type: disgord.InteractionCallbackChannelMessageWithSource,
			Data: &disgord.InteractionApplicationCommandCallbackData{
				Content: translation.T("Private", translation.GetLocale(itc)),
			},
		}
	}
	if len(userinfo.Avatars) == 0 {
		return &disgord.InteractionResponse{
			Type: disgord.InteractionCallbackChannelMessageWithSource,
			Data: &disgord.InteractionApplicationCommandCallbackData{
				Content: translation.T("NotFoundAvatars", translation.GetLocale(itc)),
			},
		}
	}
	count := 0
	handler.Client.SendInteractionResponse(context.Background(), itc, &disgord.InteractionResponse{
		Type: disgord.InteractionCallbackChannelMessageWithSource,
		Data: &disgord.InteractionApplicationCommandCallbackData{
			Embeds: []*disgord.Embed{
				{
					Color: 65535,
					Title: fmt.Sprintf("Avatar %d", count+1),
					Image: &disgord.EmbedImage{
						URL: userinfo.Avatars[0],
					},
				},
			},
			Components: []*disgord.MessageComponent{
				{
					Type: disgord.MessageComponentActionRow,
					Components: []*disgord.MessageComponent{
						{
							Type:  disgord.MessageComponentButton,
							Style: disgord.Primary,
							Emoji: &disgord.Emoji{
								Name: "⬅️",
							},
							Label:    "\u200f",
							CustomID: "back",
						},
						{
							Type:  disgord.MessageComponentButton,
							Style: disgord.Primary,
							Emoji: &disgord.Emoji{
								Name: "➡️",
							},
							Label:    "\u200f",
							CustomID: "next",
						},
					},
				},
			},
		},
	})
	handler.RegisterHandler(itc.ID, func(interaction *disgord.InteractionCreate) {
		if itc.Member.User.ID == interaction.Member.User.ID {
			if interaction.Data.CustomID == "back" {
				if count == 0 {
					count = len(userinfo.Avatars) - 1
				} else {
					count--
				}
			} else {
				if count == len(userinfo.Avatars)-1 {
					count = 0
				} else {
					count++
				}
			}

			handler.Client.SendInteractionResponse(context.Background(), interaction, &disgord.InteractionResponse{
				Type: disgord.InteractionCallbackUpdateMessage,
				Data: &disgord.InteractionApplicationCommandCallbackData{
					Embeds: []*disgord.Embed{
						{

							Color: 65535,
							Title: fmt.Sprintf("Avatar %d", count+1),
							Image: &disgord.EmbedImage{
								URL: userinfo.Avatars[count],
							},
						},
					},
				},
			})
		}
	}, 60*5)
	return nil
}
