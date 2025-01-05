package commands

import (
	"asura/src/database"
	"asura/src/handler"
	"asura/src/rinha"
	"context"
	"fmt"

	"github.com/andersfylling/disgord"
)

func init() {
	handler.RegisterCommand(handler.Command{
		Name:        "money",
		Description: "Ver dinheiro e xp",
		Run:         runMoney,
		Aliases:     []string{"dinheiro"},
		Cooldown:    5,
		Category:    handler.Profile,
	})
}

func runMoney(ctx context.Context, itc *disgord.InteractionCreate) *disgord.CreateInteractionResponse {
	user := database.User.GetUser(ctx, itc.Member.UserID)
	return &disgord.CreateInteractionResponse{
		Type: disgord.InteractionCallbackChannelMessageWithSource,
		Data: &disgord.CreateInteractionResponseData{
			Embeds: []*disgord.Embed{
				{
					Title:       "Money",
					Color:       65535,
					Description: fmt.Sprintf("Dinheiro: **%d**\nAsuraCoin: **%d**\nUserXP: **%d**\nTreinos: **%d/%d**\n\nUse `/lootbox view` para visualizar ou comprar lootboxes\nUse `/givemoney` para doar dinheiro\nUse `/daily` para pegar o bonus diario\nUse `/trial status` para fortificar seu galo\n\n**Entre no meu [Servidor](https://discord.gg/CfkBZyVsd7) para ganhar um bonus no daily**\n**[Comprar Moedas e XP](https://acnologla.github.io/asura-site/donate)**", user.Money, user.AsuraCoin, user.UserXp, user.TrainLimit, rinha.CalcLimit(&user)),
				},
			},
		},
	}
}
