package commands

import (
	"asura/src/entities"
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
	return entities.CreateMsg().Content("https://discord.gg/CfkBZyVsd7").Res()
}
