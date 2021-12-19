package commands

import (
	"asura/src/handler"
	"asura/src/interaction"
)

func init() {
	handler.RegisterCommand(handler.Command{
		Name:        "ping",
		Description: "Pinga o usuario",
		Run:         runPing,
	})
}

func runPing(itc interaction.Interaction) *interaction.InteractionResponse {
	return &interaction.InteractionResponse{
		Type: interaction.CHANNEL_MESSAGE_WITH_SOURCE,
		Data: &interaction.InteractionCallbackData{
			Content: "Pinguei vc lindo",
		},
	}
}
