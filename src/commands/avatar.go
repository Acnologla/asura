package commands

import (
	"asura/src/handler"
	"asura/src/interaction"
	"asura/src/translation"
	"asura/src/utils"
)

func init() {
	handler.RegisterCommand(handler.Command{
		Name:        "avatar",
		Description: translation.T("AvatarHelp", "pt"),
		Run:         runAvatar,
		Cooldown:    5,
		Options: utils.GenerateOptions(&interaction.ApplicationCommandOption{
			Name:        "user",
			Type:        interaction.USER,
			Description: "user avatar",
			Required:    true,
		}),
	})
}

func runAvatar(itc interaction.Interaction) *interaction.InteractionResponse {
	user := utils.GetUser(itc.Data.Options, 0)
	return &interaction.InteractionResponse{
		Type: interaction.CHANNEL_MESSAGE_WITH_SOURCE,
		Data: &interaction.InteractionCallbackData{
			Embeds: []*interaction.Embed{
				{
					Title: user.Username,
					Image: &interaction.Image{
						URL: user.AvatarURL(128),
					},
				},
			},
		},
	}
}
