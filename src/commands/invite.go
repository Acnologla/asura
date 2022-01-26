package commands

import (
	"asura/src/handler"
	"asura/src/translation"

	"github.com/andersfylling/disgord"
)

func init() {
	handler.RegisterCommand(handler.Command{
		Name:        "invite",
		Description: translation.T("InviteHelp", "pt"),
		Run:         runInvite,
		Cooldown:    5,
	})
}

func runInvite(itc *disgord.InteractionCreate) *disgord.InteractionResponse {
	return &disgord.InteractionResponse{
		Type: disgord.InteractionCallbackChannelMessageWithSource,
		Data: &disgord.InteractionApplicationCommandCallbackData{
			Content: " https://discordapp.com/oauth2/authorize?client_id=470684281102925844&scope=applications.commands%%20bot&permissions=8",
		},
	}
}
