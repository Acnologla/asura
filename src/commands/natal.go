package commands

import (
	"asura/src/entities"
	"asura/src/handler"
	"context"

	"github.com/andersfylling/disgord"
)

func init() {
	handler.RegisterCommand(handler.Command{
		Name:        "natal",
		Description: "Comando do evento de natal",
		Run:         runNatal,
		Cooldown:    5,
		Category:    handler.Rinha,
	})
}

func runNatal(ctx context.Context, itc *disgord.InteractionCreate) *disgord.CreateInteractionResponse {
	return entities.CreateMsg().
		Content("Evento acabou").
		Res()
}
