package commands

import (
	"asura/src/handler"
	"asura/src/translation"
	"asura/src/utils"
	"context"

	"github.com/andersfylling/disgord"
)

func init() {
	handler.RegisterCommand(handler.Command{
		Name:        "avatar",
		Description: translation.T("AvatarHelp", "pt"),
		Run:         runAvatar,
		Cooldown:    5,
		Options: utils.GenerateOptions(&disgord.ApplicationCommandOption{
			Name:        "user",
			Type:        disgord.OptionTypeUser,
			Description: "user avatar",
			Required:    false,
		}),
	})
}

func runAvatar(ctx context.Context, itc *disgord.InteractionCreate) *disgord.CreateInteractionResponse {
	user := itc.Member.User
	if len(itc.Data.Options) > 0 {
		user = utils.GetUser(itc, 0)
	}
	avatar, _ := user.AvatarURL(1024, true)
	return &disgord.CreateInteractionResponse{
		Type: disgord.InteractionCallbackChannelMessageWithSource,
		Data: &disgord.CreateInteractionResponseData{
			Embeds: []*disgord.Embed{
				{
					Title: user.Username,
					Image: &disgord.EmbedImage{
						URL: avatar,
					},
				},
			},
		},
	}
}
