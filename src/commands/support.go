package commands

import (
	"asura/src/handler"
	"context"

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

func runSupport(ctx context.Context, itc *disgord.InteractionCreate) *disgord.CreateInteractionResponse {
	return &disgord.CreateInteractionResponse{
		Type: disgord.InteractionCallbackChannelMessageWithSource,
		Data: &disgord.CreateInteractionResponseData{
			Content: "https://discord.gg/tdVWQGV",
		},
	}
}
