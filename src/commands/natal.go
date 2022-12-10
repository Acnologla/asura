package commands

import (
	"asura/src/database"
	"asura/src/entities"
	"asura/src/handler"
	"asura/src/rinha"
	"asura/src/rinha/engine"
	"context"

	"github.com/andersfylling/disgord"
)

func init() {
	handler.RegisterCommand(handler.Command{
		Name:        "natal",
		Description: "Comando do evento de natal",
		Run:         runNatal,
		Cooldown:    5,
	})
}

func runNatal(ctx context.Context, itc *disgord.InteractionCreate) *disgord.CreateInteractionResponse {
	discordUser := itc.Member.User
	authorRinha := isInRinha(ctx, discordUser)
	user := database.User.GetUser(ctx, itc.Member.UserID, "Galos")
	galo := rinha.GetEquippedGalo(&user)
	if authorRinha != "" {
		handler.Client.Channel(itc.ChannelID).CreateMessage(&disgord.CreateMessage{
			Content: rinhaMessage(discordUser.Username, authorRinha).Data.Content,
		})
		return rinhaMessage(discordUser.Username, authorRinha)
	}
	galoType := 32
	advLevel := 200
	galoAdv := entities.Rooster{
		Xp:     rinha.CalcXP(advLevel) + 1,
		Type:   galoType,
		Equip:  true,
		Resets: 5,
	}
	lockEvent(ctx, discordUser.ID, "Boss evento"+rinha.Classes[galoAdv.Type].Name)
	defer unlockEvent(ctx, discordUser.ID)

	itc.Reply(ctx, handler.Client, &disgord.CreateInteractionResponse{
		Type: disgord.InteractionCallbackChannelMessageWithSource,
		Data: &disgord.CreateInteractionResponseData{
			Content: "A batalha esta iniciando",
		},
	})
	winner, _ := engine.ExecuteRinha(itc, handler.Client, engine.RinhaOptions{
		GaloAuthor: &user,
		GaloAdv: &entities.User{
			Galos: []*entities.Rooster{&galoAdv},
		},
		IDs: [2]disgord.Snowflake{discordUser.ID},

		AuthorName:  rinha.GetName(discordUser.Username, *galo),
		AdvName:     "Boss evento " + rinha.Classes[galoAdv.Type].Name,
		AuthorLevel: rinha.CalcLevel(galo.Xp),
		AdvLevel:    rinha.CalcLevel(galoAdv.Xp),
		NoItems:     true,
	}, false)
	if winner == -1 {
		return nil
	}
	ch := handler.Client.Channel(disgord.Snowflake(itc.ChannelID))

	if winner == 0 {
		database.User.UpdateUser(ctx, discordUser.ID, func(u entities.User) entities.User {
			database.User.InsertItem(ctx, discordUser.ID, u.Items, 43, entities.NormalType)
			return u
		}, "Items")
		ch.CreateMessage(&disgord.CreateMessage{
			Content: "Voce ganhou um presente de natal",
		})
	} else {
		ch.CreateMessage(&disgord.CreateMessage{
			Content: "Parabens voce perdeu",
		})
	}
	return nil
}
