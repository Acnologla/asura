package commands

import (
	"asura/src/database"
	"asura/src/entities"
	"asura/src/handler"
	"asura/src/utils"
	"context"

	"asura/src/translation"

	"github.com/andersfylling/disgord"
)

func init() {
	handler.RegisterCommand(handler.Command{
		Name:        "givemoney",
		Description: translation.T("GiveMoneyHelp", "pt"),
		Run:         runGiveMoney,
		Cooldown:    15,
		Category:    handler.Profile,
		Options: utils.GenerateOptions(
			&disgord.ApplicationCommandOption{
				Type:        disgord.OptionTypeUser,
				Required:    true,
				Name:        "user",
				Description: "user to give money",
			},
			&disgord.ApplicationCommandOption{
				Type:        disgord.OptionTypeNumber,
				Required:    true,
				MinValue:    0,
				MaxValue:    12000,
				Name:        "money",
				Description: "money to give",
			}),
	})
}

func runGiveMoney(ctx context.Context, itc *disgord.InteractionCreate) *disgord.CreateInteractionResponse {
	user := utils.GetUser(itc, 0)
	if user.Bot || user.ID == itc.Member.UserID {
		return &disgord.CreateInteractionResponse{
			Type: disgord.InteractionCallbackChannelMessageWithSource,
			Data: &disgord.CreateInteractionResponseData{
				Content: "invalid user",
			},
		}
	}
	money := int(itc.Data.Options[1].Value.(float64))
	if money < 1 {
		return &disgord.CreateInteractionResponse{
			Type: disgord.InteractionCallbackChannelMessageWithSource,
			Data: &disgord.CreateInteractionResponseData{
				Content: translation.T("GiveMoneyInvalid", translation.GetLocale(itc)),
			},
		}
	}
	var msg string
	database.User.UpdateUser(ctx, itc.Member.UserID, func(u entities.User) entities.User {
		if money > u.Money {
			msg = "NoMoney"
			return u
		}
		msg = "GiveMoney"
		u.Money -= money
		database.User.UpdateUser(ctx, user.ID, func(u entities.User) entities.User {
			u.Money += money
			return u
		})
		return u
	})
	return &disgord.CreateInteractionResponse{
		Type: disgord.InteractionCallbackChannelMessageWithSource,
		Data: &disgord.CreateInteractionResponseData{
			Content: translation.T(msg, translation.GetLocale(itc), money),
		},
	}
}
