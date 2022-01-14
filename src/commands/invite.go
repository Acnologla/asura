package commands

import (
	"asura/src/handler"
	"asura/src/interaction"
	"fmt"
)

func init() {
	handler.RegisterCommand(handler.Command{
		Name:        "invite",
		Description: "Convidar o bot para seu servidor",
		Run:         runInvite,
	})
}

func runInvite(itc interaction.Interaction) *interaction.InteractionResponse {
	return &interaction.InteractionResponse{
		Type: interaction.CHANNEL_MESSAGE_WITH_SOURCE,
		Data: &interaction.InteractionCallbackData{
			Content: fmt.Sprintf("https://discordapp.com/oauth2/authorize?client_id=%s&scope=bot&permissions=8", "470684281102925844"),
		},
	}
}
