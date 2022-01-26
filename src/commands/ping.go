package commands

import (
	"asura/src/handler"

	"asura/src/translation"

	"github.com/andersfylling/disgord"
)

func init() {
	handler.RegisterCommand(handler.Command{
		Name:        "ping",
		Description: translation.T("PingHelp", "pt"),
		Run:         runPing,
		Cooldown:    3,
	})
}

func runPing(itc *disgord.InteractionCreate) *disgord.InteractionResponse {
	return &disgord.InteractionResponse{
		Type: disgord.InteractionCallbackChannelMessageWithSource,
		Data: &disgord.InteractionApplicationCommandCallbackData{
			Content: translation.T("Ping", translation.GetLocale(itc), "12ms"),
		},
	}
}
