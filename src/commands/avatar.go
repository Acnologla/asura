package commands

import (
	"asura/src/handler"
	"asura/src/translation"
	"asura/src/utils"

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
			Required:    true,
		}),
	})
}

func runAvatar(itc *disgord.InteractionCreate) *disgord.CreateInteractionResponse {
	user := utils.GetUser(itc, 0)
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
