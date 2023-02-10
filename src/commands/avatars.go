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
		Cache:       60,
		Options: utils.GenerateOptions(&disgord.ApplicationCommandOption{
			Name:        "user",
			Type:        disgord.OptionTypeUser,
			Description: "user avatar",
			Required:    true,
		}),
	})
}

func runAvatars(ctx context.Context, itc *disgord.InteractionCreate) *disgord.CreateInteractionResponse {
	user := utils.GetUser(itc, 0)
	var userinfo firebase.User
	var private bool
	id := strconv.FormatUint(uint64(user.ID), 10)
	firebase.Database.NewRef("users/"+id).Get(ctx, &userinfo)
	firebase.Database.NewRef("private/"+id).Get(ctx, &private)
	if private {
		return &disgord.CreateInteractionResponse{
			Type: disgord.InteractionCallbackChannelMessageWithSource,
			Data: &disgord.CreateInteractionResponseData{
				Content: translation.T("Private", translation.GetLocale(itc)),
			},
		}
	}
	if len(userinfo.Avatars) == 0 {
		return &disgord.CreateInteractionResponse{
			Type: disgord.InteractionCallbackChannelMessageWithSource,
			Data: &disgord.CreateInteractionResponseData{
				Content: translation.T("NotFoundAvatars", translation.GetLocale(itc)),
			},
		}
	}
	count := 0
	components := []*disgord.MessageComponent{
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
	}
	handler.Client.SendInteractionResponse(ctx, itc, &disgord.CreateInteractionResponse{
		Type: disgord.InteractionCallbackChannelMessageWithSource,
		Data: &disgord.CreateInteractionResponseData{
			Embeds: []*disgord.Embed{
				{
					Color: 65535,
					Title: fmt.Sprintf("Avatar %d", count+1),
					Image: &disgord.EmbedImage{
						URL: userinfo.Avatars[0],
					},
				},
			},
			Components: components,
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

			handler.Client.SendInteractionResponse(ctx, interaction, &disgord.CreateInteractionResponse{
				Type: disgord.InteractionCallbackUpdateMessage,
				Data: &disgord.CreateInteractionResponseData{
					Embeds: []*disgord.Embed{
						{

							Color: 65535,
							Title: fmt.Sprintf("Avatar %d", count+1),
							Image: &disgord.EmbedImage{
								URL: userinfo.Avatars[count],
							},
						},
					},
					Components: components,
				},
			})
		}
	}, 60*5)
	return nil
}
