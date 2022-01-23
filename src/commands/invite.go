package commands

import (
	"asura/src/handler"
	"asura/src/interaction"
	"asura/src/translation"
)

func init() {
	handler.RegisterCommand(handler.Command{
		Name:        "invite",
		Description: translation.T("InviteHelp", "pt"),
		Run:         runInvite,
		Cooldown:    5,
	})
}

func runInvite(itc interaction.Interaction) *interaction.InteractionResponse {
	return &interaction.InteractionResponse{
		Type: interaction.CHANNEL_MESSAGE_WITH_SOURCE,
		Data: &interaction.InteractionCallbackData{
			Content: " https://discordapp.com/oauth2/authorize?client_id=470684281102925844&scope=applications.commands%%20bot&permissions=8",
		},
	}
}
