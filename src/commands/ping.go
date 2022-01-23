package commands

import (
	"asura/src/handler"
	"asura/src/interaction"

	"asura/src/translation"
)

func init() {
	handler.RegisterCommand(handler.Command{
		Name:        "ping",
		Description: "Pinga o usuario",
		Run:         runPing,
		Cooldown:    3,
	})
}

func runPing(itc interaction.Interaction) *interaction.InteractionResponse {
	return &interaction.InteractionResponse{
		Type: interaction.CHANNEL_MESSAGE_WITH_SOURCE,
		Data: &interaction.InteractionCallbackData{
			Content: translation.T("Ping", itc.GuildLocale, "12ms"),
		},
	}
}
