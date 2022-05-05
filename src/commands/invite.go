package commands

import (
	"asura/src/entities"
	"asura/src/handler"
	"asura/src/translation"
	"context"

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

func runInvite(ctx context.Context, itc *disgord.InteractionCreate) *disgord.CreateInteractionResponse {
	return entities.CreateMsg().
		Content("https://discordapp.com/oauth2/authorize?client_id=470684281102925844&scope=applications.commands%20bot&permissions=8").
		Res()
}
