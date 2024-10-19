package commands

import (
	"asura/src/database"
	"asura/src/entities"
	"asura/src/handler"
	"asura/src/utils"
	"context"
	"time"

	"asura/src/translation"

	"github.com/andersfylling/disgord"
)

const MAX_TRANSACTIONS = 10000

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
				MinValue:    120,
				MaxValue:    MAX_TRANSACTIONS,
				Name:        "money",
				Description: "money to give",
			}),
	})
}

func get24HoursTransactions(arr []*entities.Transaction) (transactions []*entities.Transaction) {
	for _, transaction := range arr {
		if time.Now().Unix()-transaction.CreatedAt < 60*60*24 {
			transactions = append(transactions, transaction)
		}
	}
	return
}

func sumTransactions(arr []*entities.Transaction) (sum int) {
	for _, transaction := range arr {
		sum += transaction.Amount
	}
	return
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
	database.User.UpdateUser(ctx, itc.Member.UserID, func(uAuthor entities.User) entities.User {
		if money > uAuthor.Money {
			msg = "NoMoney"
			return uAuthor
		}

		msg = "GiveMoney"
		database.User.UpdateUser(ctx, user.ID, func(u entities.User) entities.User {
			todayTransactions := get24HoursTransactions(u.Transactions)
			if sumTransactions(todayTransactions)+money >= MAX_TRANSACTIONS {
				msg = "MaxTransactions"
				return u
			}

			uAuthor.Money -= money
			u.Money += money
			database.User.InsertTransaction(ctx, u.ID, &entities.Transaction{
				Amount:    money,
				AuthorID:  itc.Member.UserID,
				CreatedAt: time.Now().Unix(),
			})

			return u
		}, "Transactions")
		return uAuthor
	})
	return &disgord.CreateInteractionResponse{
		Type: disgord.InteractionCallbackChannelMessageWithSource,
		Data: &disgord.CreateInteractionResponseData{
			Content: translation.T(msg, translation.GetLocale(itc), money),
		},
	}
}
