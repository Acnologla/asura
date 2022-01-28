package commands

import (
	"asura/src/handler"

	"github.com/andersfylling/disgord"
)

func init() {
	handler.RegisterCommand(handler.Command{
		Name:        "suporte",
		Description: "suport server",
		Run:         runSupport,
		Cooldown:    3,
	})
}

func runSupport(itc *disgord.InteractionCreate) *disgord.InteractionResponse {
	return &disgord.InteractionResponse{
		Type: disgord.InteractionCallbackChannelMessageWithSource,
		Data: &disgord.InteractionApplicationCommandCallbackData{
			Content: "https://discord.gg/tdVWQGV",
		},
	}
}
